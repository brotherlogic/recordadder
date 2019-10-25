package main

import (
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

	queue.Requests = append(queue.Requests, req)

	err = s.KSclient.Save(ctx, QUEUE, queue)
	return &pb.AddRecordResponse{ExpectedAdditionDate: time.Now().Add(time.Hour * time.Duration((24 * len(queue.Requests)))).Unix()}, err
}

// Test test function
func (s *Server) Test(ctx context.Context, req *pb.AddRecordRequest) (*pb.AddRecordResponse, error) {
	time.Sleep(time.Minute)
	return nil, nil
}
