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
	pbgd "github.com/brotherlogic/godiscogs/proto"
	pbg "github.com/brotherlogic/goserver/proto"
	pb "github.com/brotherlogic/recordadder/proto"
	rbpb "github.com/brotherlogic/recordbudget/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
)

var (
	adds = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordadder_adds_24hours",
		Help: "The number of records we know of that have arrived but we haven't moved yet",
	}, []string{"dest"})
	dones = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordadder_dones",
		Help: "The number of records we know of that have arrived but we haven't moved yet",
	}, []string{"dest"})
)

func (s *Server) updateMetrics(ctx context.Context, queue *pb.Queue) {
	recentAdds := float64(0)
	recentAddsSeven := float64(0)
	for _, entry := range queue.GetAdded() {
		if time.Since(time.Unix(entry.GetDateAdded(), 0)) < time.Hour*18 && entry.GetFolderId() == 242017 {
			recentAdds++
		}
		if time.Since(time.Unix(entry.GetDateAdded(), 0)) < time.Hour*18 && entry.GetFolderId() == 267116 {
			recentAddsSeven++
		}
	}
	adds.With(prometheus.Labels{"dest": "12"}).Set((recentAdds))
	adds.With(prometheus.Labels{"dest": "7"}).Set((recentAddsSeven))
}

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
	getRecord(ctx context.Context, id int32) (*pbrc.Record, error)
}

type prodCollection struct {
	dial func(ctx context.Context, server string) (*grpc.ClientConn, error)
}

func (p *prodCollection) getRecord(ctx context.Context, id int32) (*pbrc.Record, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	resp, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: id})
	if err != nil {
		return nil, err
	}
	return resp.GetRecord(), nil
}

func (p *prodCollection) addRecord(ctx context.Context, r *pb.AddRecordRequest) (int32, error) {
	conn, err := p.dial(ctx, "recordcollection")
	if err != nil {
		return -1, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	// Vanilla Addition
	var pl pbrc.PurchaseLocation
	switch r.GetPurchaseLocation() {
	case "amoeba":
		pl = pbrc.PurchaseLocation_AMOEBA
	case "hercules":
		pl = pbrc.PurchaseLocation_HERCULES
	case "stranded":
		pl = pbrc.PurchaseLocation_STRANDED
	case "discogs":
		pl = pbrc.PurchaseLocation_DISCOGS
	case "gift":
		pl = pbrc.PurchaseLocation_GIFT
	case "bandcamp":
		pl = pbrc.PurchaseLocation_PBANDCAMP
	case "download":
		pl = pbrc.PurchaseLocation_DOWNLOAD
	case "cherry":
		pl = pbrc.PurchaseLocation_CHERRY_RED
	case "bleep":
		pl = pbrc.PurchaseLocation_BLEEP
	case "direct":
		pl = pbrc.PurchaseLocation_DIRECT
	case "groovemerchant":
		pl = pbrc.PurchaseLocation_GROOVE_MERCHANT
	case "sacredbones":
		pl = pbrc.PurchaseLocation_SACRED_BONES
	default:
		return -1, fmt.Errorf("Unknown location %v", r.GetPurchaseLocation())
	}
	resp, err := client.AddRecord(ctx, &pbrc.AddRecordRequest{ToAdd: &pbrc.Record{
		Release:  &pbgd.Release{Id: r.Id},
		Metadata: &pbrc.ReleaseMetadata{WasParents: r.WasParents, Cost: r.Cost, GoalFolder: r.Folder, AccountingYear: r.AccountingYear, PurchaseLocation: pl},
	}})
	if err != nil {
		return -1, err
	}

	// Also dial grambridge and ping it
	conn2, err := p.dial(ctx, "grambridge")
	if err != nil {
		return -1, err
	}
	gbc := pbrc.NewClientUpdateServiceClient(conn2)
	gbc.ClientUpdate(ctx, &pbrc.ClientUpdateRequest{InstanceId: resp.GetAdded().GetRelease().GetInstanceId()})

	return resp.GetAdded().GetRelease().GetInstanceId(), nil

}

func (s *Server) getUnlistenedSevens(ctx context.Context) (int, error) {
	conn, err := s.FDialServer(ctx, "recordcollection")
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{
		Query: &pbrc.QueryRecordsRequest_Category{Category: pbrc.ReleaseMetadata_UNLISTENED},
	})
	if err != nil {
		return 0, err
	}

	count := 0
	for _, r := range recs.GetInstanceIds() {
		rec, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: r})
		if err != nil {
			return 0, err
		}
		if rec.GetRecord().GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_7_INCH {
			count++
		}
	}

	return count, nil
}

func (s *Server) getUnlistenedCDs(ctx context.Context) (int, error) {
	conn, err := s.FDialServer(ctx, "recordcollection")
	defer conn.Close()
	if err != nil {
		return 0, err
	}
	client := pbrc.NewRecordCollectionServiceClient(conn)
	recs, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{
		Query: &pbrc.QueryRecordsRequest_Category{Category: pbrc.ReleaseMetadata_UNLISTENED},
	})
	if err != nil {
		return 0, err
	}

	count := 0
	for _, r := range recs.GetInstanceIds() {
		rec, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: r})
		if err != nil {
			return 0, err
		}
		if rec.GetRecord().GetMetadata().GetFiledUnder() == pbrc.ReleaseMetadata_FILE_CD {
			count++
		}
	}

	return count, nil
}

// Server main server type
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
		GoServer:    &goserver.GoServer{},
		running:     true,
		fanout:      []string{},
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
	pbrc.RegisterClientUpdateServiceServer(server, s)
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
	added = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "recordadder_last_add",
		Help: "The last addition for each type",
	}, []string{"type"})
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

		done, err := s.RunLockingElection(ctx, "recordadder", "Locking to add a record")
		if err == nil {
			err := s.processQueue(ctx)
			s.CtxLog(ctx, fmt.Sprintf("Ran queue: %v", err))
		} else {
			s.CtxLog(ctx, fmt.Sprintf("Unable to get lock: %v", err))
		}
		s.ReleaseLockingElection(ctx, "recordadder", done)

		cancel()

		time.Sleep(time.Minute * 15)
	}

	return nil
}

const (
	CONFIG_KEY = "github.com/brotherlogic/recordadder/moverconfig"
)

func (s *Server) loadConfig(ctx context.Context) (*pb.MoverConfig, error) {
	if s.testing {
		return &pb.MoverConfig{}, nil
	}

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

	if config.TodayFolders == nil {
		config.TodayFolders = make(map[int32]int32)
	}

	for t, d := range config.GetAddedMap() {
		added.With(prometheus.Labels{"type": t}).Set(float64(d))
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
	server.PrepServer("recordadder")
	server.Register = server

	err := server.RegisterServerV2(false)
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
