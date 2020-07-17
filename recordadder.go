package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/goserver/utils"
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
	addRecord(ctx context.Context, r *pb.AddRecordRequest) error
}

type prodCollection struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodCollection) addRecord(ctx context.Context, r *pb.AddRecordRequest) error {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	// Vanilla Addition
	if r.GetResetFolder() == 0 {
		_, err = client.AddRecord(ctx, &pbrc.AddRecordRequest{ToAdd: &pbrc.Record{
			Release:  &pbgd.Release{Id: r.Id},
			Metadata: &pbrc.ReleaseMetadata{Cost: r.Cost, GoalFolder: r.Folder, AccountingYear: r.AccountingYear},
		}})
		return err
	}

	// Folder reset
	_, err = client.UpdateRecord(ctx, &pbrc.UpdateRecordRequest{
		Reason: "recordadder-folderupdate",
		Update: &pbrc.Record{
			Release:  &pbgd.Release{InstanceId: r.GetId()},
			Metadata: &pbrc.ReleaseMetadata{GoalFolder: r.GetResetFolder(), SetRating: -1},
		}})
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	rc      collection
	budget  budget
	running bool
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
		running:  true,
	}
	s.rc = &prodCollection{dial: s.FDialServer}
	s.budget = &prodBudget{dial: s.FDialServer}
	return s
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

func (s *Server) load(ctx context.Context) (*pb.Queue, error) {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return nil, err
	}
	queue := data.(*pb.Queue)
	return queue, nil
}

func (s *Server) runTimedTask() {
	time.Sleep(time.Second * 10)

	ctx, cancel := utils.ManualContext("adder-load", "adder-load", time.Minute, true)
	queue, err := s.load(ctx)
	cancel()
	if err != nil {
		log.Fatalf("Unable to initialize with queue")
	}
	for s.running {
		s.Log(fmt.Sprintf("Sleeping for %v -> %v, %v, %v", time.Unix(queue.LastAdditionDate, 0).Add(time.Hour*24).Sub(time.Now()), time.Now(), time.Unix(queue.LastAdditionDate, 0), time.Unix(queue.LastAdditionDate, 0).Add(time.Hour*24)))
		time.Sleep(time.Unix(queue.LastAdditionDate, 0).Add(time.Hour * 24).Sub(time.Now()))
		ctx, cancel = utils.ManualContext("adder-load", "adder-load", time.Minute, true)
		queue, err = s.load(ctx)
		cancel()
		if err == nil && time.Now().After(time.Unix(queue.LastAdditionDate, 0).Add(time.Hour*24)) {
			done, err := s.Elect()
			if err == nil {
				ctx, cancel = utils.ManualContext("adder-load", "adder-load", time.Minute, true)
				err := s.processQueue(ctx)
				s.Log(fmt.Sprintf("Ran queue: %v", err))
				cancel()
			}
			done()
		}

		time.Sleep(time.Minute)
	}
}

func main() {
	var quiet = flag.Bool("quiet", false, "Show all output")
	var init = flag.Bool("init", false, "Prep server")
	flag.Parse()

	//Turn off logging
	if *quiet {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
	server := Init()
	server.PrepServer()
	server.Register = server

	err := server.RegisterServerV2("recordadder", false, true)
	if err != nil {
		return
	}

	if *init {
		ctx, cancel := utils.BuildContext("recordadder", "recordadder")
		defer cancel()

		err := server.KSclient.Save(ctx, QUEUE, &pb.Queue{ProcessedRecords: 1})
		fmt.Printf("Initialised: %v\n", err)
		return
	}

	go server.runTimedTask()

	fmt.Printf("%v", server.Serve())
}
