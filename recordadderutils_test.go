package main

import (
	"fmt"
	"testing"

	rbpb "github.com/brotherlogic/recordbudget/proto"
	"golang.org/x/net/context"

	pbgd "github.com/brotherlogic/godiscogs/proto"
	pb "github.com/brotherlogic/recordadder/proto"
	pbrc "github.com/brotherlogic/recordcollection/proto"
)

type testCollection struct {
	addedRecord *pbrc.Record
	fail        bool
}

type testBudget struct {
	fail bool
}

func (p *testBudget) getBudget(ctx context.Context) (*rbpb.GetBudgetResponse, error) {
	if p.fail {
		return nil, fmt.Errorf("Built to fail")
	}
	return &rbpb.GetBudgetResponse{Spends: 100, Budget: 200}, nil
}

func (p *testCollection) addRecord(ctx context.Context, r *pb.AddRecordRequest) (int32, error) {
	if p.fail {
		return -1, fmt.Errorf("Built to fail")
	}
	p.addedRecord = &pbrc.Record{Release: &pbgd.Release{Id: r.Id}}
	return 123, nil
}

func TestBasicRunThrough(t *testing.T) {
	s := InitTestServer()
	tc := &testCollection{}
	s.rc = tc

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 12, Cost: 12, Arrived: true, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err != nil {
		t.Errorf("Error processing queue: %v", err)
	}

	if tc.addedRecord == nil || tc.addedRecord.Release.Id != 12 {
		t.Errorf("Record was not added: %v", tc.addedRecord)
	}
}

func TestBasicRunThroughWithFanoutFail(t *testing.T) {
	s := InitTestServer()
	tc := &testCollection{}
	s.rc = tc
	s.testingFail = true

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 12, Cost: 12, Arrived: true, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err != nil {
		t.Errorf("Error processing queue: %v", err)
	}

	if tc.addedRecord == nil || tc.addedRecord.Release.Id != 12 {
		t.Errorf("Record was not added: %v", tc.addedRecord)
	}
}

func TestBasicRunThroughWithBudgetFail(t *testing.T) {
	s := InitTestServer()
	tc := &testCollection{}
	tb := &testBudget{fail: true}
	s.budget = tb
	s.rc = tc

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 12, Cost: 12, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err != nil {
		t.Errorf("Error processing queue: %v", err)
	}
}

func TestEmptyRunThrough(t *testing.T) {
	s := InitTestServer()

	err := s.processQueue(context.Background())
	if err != nil {
		t.Errorf("Processing an empty queue led to errors")
	}
}

func TestBadReadRunThrough(t *testing.T) {
	s := InitTestServer()
	s.GoServer.KSclient.Fail = true

	err := s.processQueue(context.Background())

	if err == nil {
		t.Errorf("Failed to read the queue with failing read")
	}
}

func TestRunThroughWithAddFail(t *testing.T) {
	s := InitTestServer()
	tc := &testCollection{fail: true}
	s.rc = tc

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 12, Cost: 12, Arrived: true, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err == nil {
		t.Errorf("No error processing the queue with failing add")
	}
}

func TestRunThroughWithDigitalAddFail(t *testing.T) {
	s := InitTestServer()
	tc := &testCollection{fail: true}
	s.rc = tc

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 242018, Cost: 12, Arrived: true, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err == nil {
		t.Errorf("No error processing the queue with failing add")
	}
}

func TestWithBadAdd(t *testing.T) {
	s := InitTestServer()

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: -1, Folder: 12, Cost: 12, Arrived: true, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err == nil {
		t.Errorf("No error processing the queue with failing add")
	}
}

func TestBasicRunThroughWithNoAction(t *testing.T) {
	s := InitTestServer()
	tc := &testCollection{}
	s.rc = tc

	s.AddRecord(context.Background(), &pb.AddRecordRequest{Id: 12, Folder: 12, Cost: 12, Arrived: false, PurchaseLocation: "discogs"})

	err := s.processQueue(context.Background())
	if err != nil {
		t.Errorf("Error processing queue: %v", err)
	}

}
