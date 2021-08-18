package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordadder/proto"
	google_protobuf "github.com/golang/protobuf/ptypes/any"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"

	pbgd "github.com/brotherlogic/godiscogs"
	qpb "github.com/brotherlogic/queue/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

const (
	// QUEUE - Where we store incoming requests
	QUEUE = "/github.com/brotherlogic/recordadder/queue"
)

//AddRecord adds a record into the system
func (s *Server) AddRecord(ctx context.Context, req *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return nil, err
	}
	queue := data.(*pb.Queue)

	for _, entry := range queue.Requests {
		if entry.Id == req.Id {
			return nil, fmt.Errorf("This record is already in the queue")
		}
	}

	req.DateAdded = time.Now().Unix()

	queue.Requests = append(queue.Requests, req)

	// Run the fanout
	for _, server := range s.fanout {
		// Use a new context for fanout
		ctxfinner, cancelfinner := utils.ManualContext("rasave", time.Minute)
		err := s.runFanout(ctxfinner, server, req.GetId())
		code := status.Convert(err)
		if code.Code() != codes.OK && code.Code() != codes.Unavailable {
			s.RaiseIssue(fmt.Sprintf("Fanout for %v failed", server), fmt.Sprintf("Error was %v (%v)", err, code.Code()))
		}
		cancelfinner()
	}

	err = s.KSclient.Save(ctx, QUEUE, queue)
	return &pb.AddRecordResponse{ExpectedAdditionDate: time.Now().Add(time.Hour * time.Duration((24 * len(queue.Requests)))).Unix()}, err
}

//ListQueue lists the entries in the queue
func (s *Server) ListQueue(ctx context.Context, req *pb.ListQueueRequest) (*pb.ListQueueResponse, error) {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return nil, err
	}
	queue := data.(*pb.Queue)

	return &pb.ListQueueResponse{Requests: queue.GetRequests()}, nil
}

//UpdateRecord updates a record
func (s *Server) UpdateRecord(ctx context.Context, req *pb.UpdateRecordRequest) (*pb.UpdateRecordResponse, error) {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return nil, err
	}
	queue := data.(*pb.Queue)

	updated := false
	for _, entry := range queue.Requests {
		if entry.Id == req.Id {
			if req.GetAvailable() {
				entry.Arrived = true
				updated = true
			}
		}
	}

	if !updated {
		return nil, status.Errorf(codes.NotFound, "Unable to locate %v for update", req.GetId())
	}

	return &pb.UpdateRecordResponse{}, s.KSclient.Save(ctx, QUEUE, queue)
}

//DeleteRecord remove a record from the queue
func (s *Server) DeleteRecord(ctx context.Context, req *pb.DeleteRecordRequest) (*pb.DeleteRecordResponse, error) {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return nil, err
	}
	queue := data.(*pb.Queue)

	nqueue := []*pb.AddRecordRequest{}
	for _, entry := range queue.Requests {
		if entry.Id != req.Id {
			nqueue = append(nqueue, entry)
		}
	}
	queue.Requests = nqueue

	return &pb.DeleteRecordResponse{}, s.KSclient.Save(ctx, QUEUE, queue)

}

//ProcAdded processes the added queue
func (s *Server) ProcAdded(ctx context.Context, req *pb.ProcAddedRequest) (*pb.ProcAddedResponse, error) {
	conf, err := s.loadConfig(ctx)
	if err != nil {
		return nil, err
	}

	val, ok := conf.GetAddedMap()[req.GetType()]
	s.Log(fmt.Sprintf("ADDED the MAP: %v (%v)", time.Since(time.Unix(val, 0)), time.Unix(val, 0)))
	if !ok || time.Since(time.Unix(val, 0)) > time.Hour*24 {
		s.Log("Adding!")

		conn, err := s.FDialServer(ctx, "recordcollection")
		if err != nil {
			return nil, err
		}
		defer conn.Close()

		client := pbrc.NewRecordCollectionServiceClient(conn)
		res, err := client.QueryRecords(ctx, &pbrc.QueryRecordsRequest{Query: &pbrc.QueryRecordsRequest_Category{Category: pbrc.ReleaseMetadata_ARRIVED}})
		if err != nil {
			return nil, err
		}

		var recs []*pbrc.Record
		for _, id := range res.GetInstanceIds() {
			rec, err := client.GetRecord(ctx, &pbrc.GetRecordRequest{InstanceId: id})
			if err != nil {
				return nil, err
			}

			if rec.GetRecord().GetMetadata().GetFiledUnder().String() == req.GetType() {
				recs = append(recs, rec.GetRecord())
			}
		}

		sort.SliceStable(recs, func(i, j int) bool {
			return recs[i].GetMetadata().GetDateAdded() < recs[j].GetMetadata().GetDateAdded()
		})

		s.Log(fmt.Sprintf("Found: %v with %v", len(recs), req.GetType()))

		if len(recs) > 0 {
			_, err = client.UpdateRecord(ctx, &pbrc.UpdateRecordRequest{
				Reason: "Updating for addition",
				Update: &pbrc.Record{Release: &pbgd.Release{InstanceId: recs[0].GetRelease().GetInstanceId()},
					Metadata: &pbrc.ReleaseMetadata{Category: pbrc.ReleaseMetadata_UNLISTENED}},
			})
			if err != nil {
				return nil, err
			}

			conf.AddedMap[req.GetType()] = time.Now().Unix()
			err = s.saveConfig(ctx, conf)
			if err != nil {
				return nil, err
			}
		}

		if len(recs) <= 1 {
			val = time.Now().Unix()
		}

		runTime := time.Unix(val, 0).Add(time.Hour * 24).Unix()

		conn2, err2 := s.FDialServer(ctx, "queue")
		if err2 != nil {
			return nil, err2
		}
		defer conn2.Close()
		qclient := qpb.NewQueueServiceClient(conn)
		upup := &pb.ProcAddedRequest{
			Type: req.GetType(),
		}
		data, _ := proto.Marshal(upup)
		_, err = qclient.AddQueueItem(ctx, &qpb.AddQueueItemRequest{
			QueueName: "record_adder",
			RunTime:   runTime,
			Payload:   &google_protobuf.Any{Value: data},
			Key:       fmt.Sprintf("%v", req.GetType()),
		})
	}

	return &pb.ProcAddedResponse{}, err
}
