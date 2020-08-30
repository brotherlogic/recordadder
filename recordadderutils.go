package main

import (
	"fmt"
	"time"

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

	available := budget.GetBudget() - budget.GetSpends()
	s.Log(fmt.Sprintf("Found %v entries in the queue with %v in the budget (%v)", len(queue.Requests), available, time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) >= time.Hour*24))
	time.Sleep(time.Second * 2)

	if len(queue.Requests) > 0 && time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) >= time.Hour*24 {
		for i, req := range queue.GetRequests() {
			if req.GetId() <= 0 {
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
				queue.LastAdditionDate = time.Now().Unix()
				queue.Requests = append(queue.Requests[:i], queue.Requests[i+1:]...)
				err = s.KSclient.Save(ctx, QUEUE, queue)
				return err
			}

			s.RaiseIssue("Addition Error", fmt.Sprintf("Could not add %v -> %v", available))
		}
	}

	s.Log(fmt.Sprintf("Still %v to go!", time.Now().Sub(time.Unix(queue.LastAdditionDate, 0))))

	return nil
}
