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
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	dspb "github.com/brotherlogic/dstore/proto"
	pbgd "github.com/brotherlogic/godiscogs"
	pbg "github.com/brotherlogic/goserver/proto"
	pb "github.com/brotherlogic/recordadder/proto"
	rbpb "github.com/brotherlogic/recordbudget/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
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
	var pl pbrc.ReleaseMetadata_PurchaseLocation
	switch r.GetPurchaseLocation() {
	case "amoeba":
		pl = pbrc.ReleaseMetadata_AMOEBA
	case "hercules":
		pl = pbrc.ReleaseMetadata_HERCULES
	case "stranded":
		pl = pbrc.ReleaseMetadata_STRANDED
	case "discogs":
		pl = pbrc.ReleaseMetadata_DISCOGS
	case "gift":
		pl = pbrc.ReleaseMetadata_GIFT
	case "bandcamp":
		pl = pbrc.ReleaseMetadata_PBANDCAMP
	default:
		return -1, fmt.Errorf("Unknown location %v", r.GetPurchaseLocation())
	}
	resp, err := client.AddRecord(ctx, &pbrc.AddRecordRequest{ToAdd: &pbrc.Record{
		Release:  &pbgd.Release{Id: r.Id},
		Metadata: &pbrc.ReleaseMetadata{Cost: r.Cost, GoalFolder: r.Folder, AccountingYear: r.AccountingYear, PurchaseLocation: pl},
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

	for s.running {
		time.Sleep(time.Minute)
		ctx, cancel := utils.ManualContext("adder-load", time.Minute*10)

		done, err := s.RunLockingElection(ctx, "recordadder")
		if err == nil {
			err := s.processQueue(ctx)
			s.Log(fmt.Sprintf("Ran queue: %v", err))
		} else {
			s.Log(fmt.Sprintf("Unable to get lock: %v", err))
		}
		s.ReleaseLockingElection(ctx, "recordadder", done)

		cancel()

		time.Sleep(time.Hour)
	}

	return nil
}

const (
	CONFIG_KEY = "github.com/brotherlogic/recordadder/moverconfig"
)

func (s *Server) loadConfig(ctx context.Context) (*pb.MoverConfig, error) {
	conn, err := s.FDialServer(ctx, "dstore")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := dspb.NewDStoreServiceClient(conn)
	res, err := client.Read(ctx, &dspb.ReadRequest{Key: CONFIG_KEY})
	if err != nil {
		if status.Convert(err).Code() == codes.NotFound {
			return &pb.MoverConfig{AddedMap: make(map[string]int64)}, nil
		}

		return nil, err

	}

	if res.GetConsensus() < 0.5 {
		return nil, fmt.Errorf("could not get read consensus (%v)", res.GetConsensus())
	}

	config := &pb.MoverConfig{}
	err = proto.Unmarshal(res.GetValue().GetValue(), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (s *Server) saveConfig(ctx context.Context, config *pb.MoverConfig) error {
	conn, err := s.FDialServer(ctx, "dstore")
	if err != nil {
		return err
	}
	defer conn.Close()

	data, err := proto.Marshal(config)
	if err != nil {
		return err
	}

	client := dspb.NewDStoreServiceClient(conn)
	res, err := client.Write(ctx, &dspb.WriteRequest{Key: CONFIG_KEY, Value: &google_protobuf.Any{Value: data}})
	if err != nil {
		return err
	}

	if res.GetConsensus() < 0.5 {
		return fmt.Errorf("could not get read consensus (%v)", res.GetConsensus())
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
