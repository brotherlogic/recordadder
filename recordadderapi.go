package main

import (
	"fmt"
	"time"

	pb "github.com/brotherlogic/recordadder/proto"
	"golang.org/x/net/context"
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

	queue.Requests = append(queue.Requests, req)

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

	for _, entry := range queue.Requests {
		if entry.Id == req.Id {
			if req.GetAvailable() {
				entry.Arrived = true
			}
		}
	}

	return &pb.UpdateRecordResponse{}, s.KSclient.Save(ctx, QUEUE, queue)
}

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
