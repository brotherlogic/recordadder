package main

import (
	"fmt"
	"time"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordadder/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		ctxfinner, cancelfinner := utils.ManualContext("rasave", "rasave", time.Minute, true)
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
