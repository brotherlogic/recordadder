package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbgd "github.com/brotherlogic/godiscogs"
	pbg "github.com/brotherlogic/goserver/proto"
	pb "github.com/brotherlogic/recordadder/proto"
	rbpb "github.com/brotherlogic/recordbudget/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

type budget interface {
	getBudget(ctx context.Context) (*rbpb.GetBudgetResponse, error)
}

type prodBudget struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodBudget) getBudget(ctx context.Context) (*rbpb.GetBudgetResponse, error) {
	conn, err := p.dial(ctx, "recordbudget")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := rbpb.NewRecordBudgetServiceClient(conn)
	return client.GetBudget(ctx, &rbpb.GetBudgetRequest{Year: int32(time.Now().Year())})
}

type collection interface {
	addRecord(ctx context.Context, r *pb.AddRecordRequest) (int32, error)
}

type prodCollection struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodCollection) addRecord(ctx context.Context, r *pb.AddRecordRequest) (int32, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	// Vanilla Addition
	resp, err := client.AddRecord(ctx, &pbrc.AddRecordRequest{ToAdd: &pbrc.Record{
		Release:  &pbgd.Release{Id: r.Id},
		Metadata: &pbrc.ReleaseMetadata{Cost: r.Cost, GoalFolder: r.Folder, AccountingYear: r.AccountingYear},
	}})
	if err != nil {
		return -1, err
	}
	return resp.GetAdded().GetRelease().GetInstanceId(), nil

}

//Server main server type
type Server struct {
	*goserver.GoServer
	rc          collection
	budget      budget
	running     bool
	fanout      []string
	testing     bool
	testingFail bool
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		running:  true,
		fanout: []string{
			"digitalwantlist",
		},
		testing:     false,
		testingFail: false,
	}
	s.rc = &prodCollection{dial: s.FDialServer}
	s.budget = &prodBudget{dial: s.FDialServer}
	return s
}

func (s *Server) runFanout(ctx context.Context, server string, id int32) error {
	if s.testing {
		if s.testingFail {
			return fmt.Errorf("Build to fail")
		}
		return nil
	}

	conn, err := s.FDialServer(ctx, server)
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pb.NewClientAddUpdateServiceClient(conn)
	_, err = client.ClientAddUpdate(ctx, &pb.ClientAddUpdateRequest{Id: id})
	return err
}

// DoRegister does RPC registration
func (s *Server) DoRegister(server *grpc.Server) {
	pb.RegisterAddRecordServiceServer(server, s)
}

// ReportHealth alerts if we're not healthy
func (s *Server) ReportHealth() bool {
	return true
}

// Shutdown the server
func (s *Server) Shutdown(ctx context.Context) error {
	return nil
}

// Mote promotes/demotes this server
func (s *Server) Mote(ctx context.Context, master bool) error {
	return nil
}

// GetState gets the state of the server
func (s *Server) GetState() []*pbg.State {
	return []*pbg.State{}
}

var (
	backlog = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "recordadder_backlog",
		Help: "The number of records we know of that have arrived but we haven't moved yet",
	})
	purchased = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "recordadder_purchased",
		Help: "The number of records we know of that haven't arrived",
	})
)

func (s *Server) load(ctx context.Context) (*pb.Queue, error) {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return nil, err
	}
	queue := data.(*pb.Queue)

	s.validateQueue(ctx, queue)

	return queue, nil
}

func (s *Server) validateQueue(ctx context.Context, queue *pb.Queue) {
	for _, entry := range queue.GetRequests() {
		if time.Now().Sub(time.Unix(entry.GetDateAdded(), 0)) > time.Hour*24*30 && !entry.GetArrived() {
			s.RaiseIssue("Old record in add queue", fmt.Sprintf("%v is stale in the add queue", entry.GetId()))
		}

		if entry.GetDateAdded() == 0 {
			entry.DateAdded = time.Now().Unix()
		}
	}

	backlog.Set(float64(len(queue.GetAdded())))
	purchased.Set(float64(len(queue.GetRequests())))
}

func min(t1, t2 time.Duration) time.Duration {
	if t1 < t2 {
		return t1
	}
	return t2
}

func (s *Server) runTimedTask() error {
	time.Sleep(time.Second * 10)

	var queue *pb.Queue
	var err error
	for true {
		ctx, cancel := utils.ManualContext("adder-load", time.Minute)
		s.Log(fmt.Sprintf("First read"))
		queue, err = s.load(ctx)
		cancel()
		if err == nil {
			break
		}
		time.Sleep(time.Minute)
	}
	for s.running {
		minTime := min(time.Unix(queue.LastAdditionDate, 0).Add(time.Hour*24).Sub(time.Now()), time.Unix(queue.GetLastDigitalAddition(), 0).Add(time.Hour*24).Sub(time.Now()))
		s.Log(fmt.Sprintf("Sleeping for %v", minTime))
		time.Sleep(minTime)
		ctx, cancel := utils.ManualContext("adder-load", time.Minute)
		queue, err = s.load(ctx)
		cancel()
		if err == nil && (time.Now().After(time.Unix(queue.LastAdditionDate, 0).Add(time.Hour*24)) || time.Now().After(time.Unix(queue.GetLastDigitalAddition(), 0).Add(time.Hour*24))) {

			done, err := s.Elect()
			if err == nil {
				ctx, cancel = utils.ManualContext("adder-load", time.Minute)
				err := s.processQueue(ctx)
				s.Log(fmt.Sprintf("Ran queue: %v", err))
				cancel()
			}
			done()
		}

		time.Sleep(time.Minute)
	}

	return nil
}

func main() {
	var run = flag.Bool("run", true, "Run the background adder")
	flag.Parse()

	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("recordadder", false, true)
	if err != nil {
		return
	}

	if *run {
		go func() {
			err := server.runTimedTask()
			if err != nil {
				log.Fatalf("Unable to run timed task: %v", err)
			}
		}()
	}

	fmt.Printf("%v", server.Serve())
}
