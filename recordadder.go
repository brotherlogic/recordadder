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
	_, err = client.AddRecord(ctx, &pbrc.AddRecordRequest{ToAdd: &pbrc.Record{
		Release:  &pbgd.Release{Id: r.Id},
		Metadata: &pbrc.ReleaseMetadata{Cost: r.Cost, GoalFolder: r.Folder},
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
	return []*pbg.State{}
}

func (s *Server) runTimedTask(ctx context.Context) (time.Time, error) {
	s.Log("Running the timed task")

	// Wait 24 hours between additions
	return time.Now().Add(time.Hour * 24), nil
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
	server.RegisterServer("recordadder", false)

	if *init {
		return
	}

	server.RegisterLockingTask(server.runTimedTask, "run_timed_task")

	fmt.Printf("%v", server.Serve())
}
