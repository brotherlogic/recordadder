package main

import (
	"context"
	"testing"
	"time"

	"github.com/brotherlogic/keystore/client"

	pb "github.com/brotherlogic/recordadder/proto"
)

//InitTestServer gets a test version of the server
func InitTestServer() *Server {
	s := Init()
	s.SkipLog = true
	s.GoServer.KSclient = *keystoreclient.GetTestClient(".test")
	return s
}

func TestAddRequest(t *testing.T) {
	s := InitTestServer()

	val, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{})
	if err != nil {
		t.Errorf("Add Record failed: %v", err)
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
