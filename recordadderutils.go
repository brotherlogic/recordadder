package main

import (
	"fmt"
	"time"

	"github.com/brotherlogic/goserver/utils"
	pb "github.com/brotherlogic/recordadder/proto"
	"golang.org/x/net/context"
)

func (s *Server) processQueue(ctx context.Context) error {
	s.Log(fmt.Sprintf("Processing Queue"))
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return fmt.Errorf("Error reading queue: %v", err)
	}
	queue := data.(*pb.Queue)

	budget, err := s.budget.getBudget(ctx)
	if err != nil {
		return err
	}

	available := budget.GetBudget() - budget.GetSpends()
	s.Log(fmt.Sprintf("Found %v entries in the queue with %v in the budget (%v)", len(queue.Requests), available, time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) >= time.Hour*24))
	time.Sleep(time.Second * 2)

	if len(queue.Requests) > 0 && time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) >= time.Hour*24 {
		for i, req := range queue.GetRequests() {
			if req.GetId() <= 0 || req.GetResetFolder() > 0 {
				queue.Requests = append(queue.Requests[:i], queue.Requests[i+1:]...)
				s.KSclient.Save(ctx, QUEUE, queue)
				return fmt.Errorf("Bad entry in the queue")
			}
			if req.GetCost() < available && req.GetArrived() {
				err = s.rc.addRecord(ctx, queue.Requests[i])
				s.Log(fmt.Sprintf("Adding %v -> %v", queue.Requests[i], err))
				if err != nil {
					return fmt.Errorf("Error adding record: %v", err)
				}

				// Run the fanout
				for _, server := range s.fanout {
					// Use a new context for fanout
					ctxfinner, cancelfinner := utils.ManualContext("rasave", "rasave", time.Minute, true)
					err := s.runFanout(ctxfinner, server, req.GetId())
					if err != nil {
						s.RaiseIssue(fmt.Sprintf("Fanout for %v failed", server), fmt.Sprintf("Error was %v", err))
					}
					cancelfinner()
				}

				queue.LastAdditionDate = time.Now().Unix()
				queue.Requests = append(queue.Requests[:i], queue.Requests[i+1:]...)

				// We need to refresh the context for the save since the fanout may have run out the clock
				ctxinner, cancelinner := utils.ManualContext("rasave", "rasave", time.Minute, true)
				err = s.KSclient.Save(ctxinner, QUEUE, queue)
				cancelinner()

				return err
			}
			time.Sleep(time.Second * 5)
		}
	}

	s.Log(fmt.Sprintf("Still %v to go!", time.Now().Sub(time.Unix(queue.LastAdditionDate, 0))))

	return nil
}
