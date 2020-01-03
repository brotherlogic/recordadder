package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/brotherlogic/goserver"
	"github.com/brotherlogic/keystore/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pbgd "github.com/brotherlogic/godiscogs"
	pbg "github.com/brotherlogic/goserver/proto"
	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordadder/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

type collection interface {
	addRecord(ctx context.Context, r *pb.AddRecordRequest) error
}

type prodCollection struct {
	dial func(server string) (*grpc.ClientConn, error)
}

func (p *prodCollection) addRecord(ctx context.Context, r *pb.AddRecordRequest) error {
	conn, err := p.dial("recordcollection")
	if err != nil {
		return err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)

	// Vanilla Addition
	if r.GetResetFolder() == 0 {
		_, err = client.AddRecord(ctx, &pbrc.AddRecordRequest{ToAdd: &pbrc.Record{
			Release:  &pbgd.Release{Id: r.Id},
			Metadata: &pbrc.ReleaseMetadata{Cost: r.Cost, GoalFolder: r.Folder},
		}})
		return err
	}

	// Folder reset
	_, err = client.UpdateRecord(ctx, &pbrc.UpdateRecordRequest{Update: &pbrc.Record{
		Release:  &pbgd.Release{InstanceId: r.GetId()},
		Metadata: &pbrc.ReleaseMetadata{GoalFolder: r.GetResetFolder(), SetRating: -1},
	}})
	return err
}

//Server main server type
type Server struct {
	*goserver.GoServer
	rc collection
}

// Init builds the server
func Init() *Server {
	s := &Server{
		GoServer: &goserver.GoServer{},
	}
	s.rc = &prodCollection{dial: s.DialMaster}
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
	return []*pbg.State{
		&pbg.State{Key: "no", Value: int64(233)},
	}
}

func (s *Server) runTimedTask(ctx context.Context) (time.Time, error) {
	err := s.processQueue(ctx)

	// Wait 24 hours between additions
	return time.Now().Add(time.Hour * 24), err
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
	server.GoServer.KSclient = *keystoreclient.GetClient(server.DialMaster)
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

	server.RegisterLockingTask(server.runTimedTask, "run_timed_task")

	fmt.Printf("%v", server.Serve())
}
