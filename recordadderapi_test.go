package main

import (
	"context"
	"testing"
	"time"

	keystoreclient "github.com/brotherlogic/keystore/client"
	pb "github.com/brotherlogic/recordadder/proto"
)

//InitTestServer gets a test version of the server
func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.SkipIssue = true
	s.testing = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	s.GoServer.KSclient.Save(context.Background(), QUEUE, &pb.Queue{})
	return s
}

func TestAddRequest(t *testing.T) {
	s := InitTestServer()

	val, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{})
	if err != nil {
		t.Fatalf("Add Record failed: %v", err)
	}

	if time.Now().After(time.Unix(val.ExpectedAdditionDate, 0)) {
		t.Errorf("Time is not set correct: %v", val)
	}
}

func TestAddRequestFail(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	val, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{})
	if err == nil {
		t.Errorf("Add Record with failing read did not fail: %v", val)
	}
}

func TestKeystoreFails(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	val, err := s.ListQueue(context.Background(), &pb.ListQueueRequest{})
	if err == nil {
		t.Errorf("Add Record with failing read did not fail: %v", val)
	}

	valup, err := s.UpdateRecord(context.Background(), &pb.UpdateRecordRequest{})
	if err == nil {
		t.Errorf("Update Record with failing read did not fail: %v", valup)
	}

	valdel, err := s.DeleteRecord(context.Background(), &pb.DeleteRecordRequest{})
	if err == nil {
		t.Errorf("Delete Record with failing read did not fail: %v", valdel)
	}
}

func TestDoubleAddRequest(t *testing.T) {
	s := InitTestServer()

	val, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 123})
	if err != nil {
		t.Fatalf("Add Record failed: %v", err)
	}

	val, err = s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 123})

	if err == nil {
		t.Fatalf("Double addition should have failed: %v", val)
	}
}

func TestListQueue(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 123})
	if err != nil {
		t.Fatalf("Add Record failed: %v", err)
	}

	q, err := s.ListQueue(context.Background(), &pb.ListQueueRequest{})

	if err != nil {
		t.Fatalf("Error listing queue: %v", err)
	}

	if len(q.GetRequests()) != 1 {
		t.Errorf("Wrong number of requests in queue: %v", q.GetRequests())
	}
}

func TestEmptyDelete(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 123})
	if err != nil {
		t.Fatalf("Add Record failed: %v", err)
	}

	_, err = s.DeleteRecord(context.Background(), &pb.DeleteRecordRequest{Id: 12})
	if err != nil {
		t.Errorf("Bad delete: %v", err)
	}

	q, err := s.ListQueue(context.Background(), &pb.ListQueueRequest{})
	if err != nil {
		t.Errorf("Bad list: %v", err)
	}

	if len(q.GetRequests()) != 1 {
		t.Errorf("Delete failed: %v", q)
	}
}
