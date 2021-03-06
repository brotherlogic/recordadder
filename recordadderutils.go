package main

import (
	"fmt"
	"time"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordadder/proto"
	"golang.org/x/net/context"
)

func (s *Server) processQueue(ctx context.Context) error {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return fmt.Errorf("Error reading queue: %v", err)
	}
	queue := data.(*pb.Queue)

	budget, err := s.budget.getBudget(ctx)
	if err != nil {
		return err
	}

	available := budget.GetBudget() + budget.GetSolds() - budget.GetSpends()

	lowest := int32(999999)
	for _, entry := range queue.GetRequests() {
		if entry.GetCost() < lowest {
			lowest = (entry.GetCost())
		}
	}

	s.Log(fmt.Sprintf("Found %v entries in the queue with %v in the budget (%v vs %v) -> %v", len(queue.Requests), available, lowest, time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) >= time.Hour*24, ctx))
	if len(queue.Requests) > 0 && time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) >= time.Hour*24 {
		for i, req := range queue.GetRequests() {
			if req.GetId() <= 0 || req.GetResetFolder() > 0 {
				queue.Requests = append(queue.Requests[:i], queue.Requests[i+1:]...)
				s.KSclient.Save(ctx, QUEUE, queue)
				return fmt.Errorf("Bad entry in the queue")
			}
			if !isDigital(req) && (req.GetCost() < available || int(req.GetAccountingYear()) != time.Now().Year()) && req.GetArrived() {
				err = s.rc.addRecord(ctx, queue.Requests[i])
				s.Log(fmt.Sprintf("Adding (%v) %v -> %v", i, queue.Requests[i], err))
				if err != nil {
					return fmt.Errorf("Error adding record: %v", err)
				}

				queue.LastAdditionDate = time.Now().Unix()
				queue.Requests = append(queue.Requests[:i], queue.Requests[i+1:]...)

				// We need to refresh the context for the save since the fanout may have run out the clock
				ctxinner, cancelinner := utils.ManualContext("rasave", time.Minute)
				err = s.KSclient.Save(ctxinner, QUEUE, queue)
				cancelinner()

				return err
			}
		}
	}

	err = s.runDigital(ctx, queue, available)
	if err != nil {
		return err
	}

	s.Log(fmt.Sprintf("Still %v to go (with digital %v) !", time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)), err))

	return nil
}

func isDigital(req *pb.AddRecordRequest) bool {
	//Digital is CD, Bandcamp or Computer
	return req.GetFolder() == 242018 ||
		req.GetFolder() == 1782105 ||
		req.GetFolder() == 2274270
}

func (s *Server) runDigital(ctx context.Context, queue *pb.Queue, available int32) error {
	if len(queue.Requests) > 0 && time.Now().Sub(time.Unix(queue.GetLastDigitalAddition(), 0)) >= time.Hour*24 {
		for i, req := range queue.GetRequests() {
			if isDigital(req) && (req.GetCost() < available || int(req.GetAccountingYear()) != time.Now().Year()) && req.GetArrived() {
				err := s.rc.addRecord(ctx, queue.Requests[i])
				s.Log(fmt.Sprintf("DIGITAL Adding %v -> %v", queue.Requests[i], err))
				if err != nil {
					return fmt.Errorf("Error adding digital record: %v", err)
				}

				queue.LastDigitalAddition = time.Now().Unix()
				queue.Requests = append(queue.Requests[:i], queue.Requests[i+1:]...)

				// We need to refresh the context for the save since the fanout may have run out the clock
				ctxinner, cancelinner := utils.ManualContext("rasave", time.Minute)
				err = s.KSclient.Save(ctxinner, QUEUE, queue)
				cancelinner()

				return err
			}
		}
	}
	return nil
}
