package main

import (
	"testing"

	pb "github.com/brotherlogic/recordadder/proto"
	"golang.org/x/net/context"
)

func TestFullRun(t *testing.T) {
	s := InitTestServer()

	_, err := s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 12, Cost: 12})
	if err != nil {
		t.Errorf("Bad add: %v", err)
	}

	list, err := s.ListQueue(context.Background(), &pb.ListQueueRequest{})
	if err != nil {
		t.Errorf("Bad list: %v", err)
	}

	if len(list.GetRequests()) != 1 || list.GetRequests()[0].Id != 12 {
		t.Errorf("Bad list response: %v", err)
	}

	_, err = s.UpdateRecord(context.Background(), &pb.UpdateRecordRequest{Id: 12, Available: true})
	if err != nil {
		t.Errorf("Bad update: %v", err)
	}

	list, err = s.ListQueue(context.Background(), &pb.ListQueueRequest{})
	if err != nil {
		t.Errorf("Bad list: %v", err)
	}

	if len(list.GetRequests()) != 1 || list.GetRequests()[0].Id != 12 || !list.GetRequests()[0].GetArrived() {
		t.Errorf("Bad list response: %v", err)
	}

	_, err = s.DeleteRecord(context.Background(), &pb.DeleteRecordRequest{Id: 12})
	if err != nil {
		t.Errorf("Bad delete: %v", err)
	}

	list, err = s.ListQueue(context.Background(), &pb.ListQueueRequest{})
	if err != nil {
		t.Errorf("Bad list: %v", err)
	}

	if len(list.GetRequests()) != 0 {
		t.Errorf("Bad list response: %v", err)
	}
}
