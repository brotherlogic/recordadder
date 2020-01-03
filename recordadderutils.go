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

	s.Log(fmt.Sprintf("Found %v entries in the queue", len(queue.Requests)))

	if len(queue.Requests) > 0 && time.Now().Sub(time.Unix(queue.LastAdditionDate, 0)) > time.Hour*24 {
		queueSize := len(queue.Requests)
		err = s.rc.addRecord(ctx, queue.Requests[0])
		if err != nil {
			return fmt.Errorf("Error adding record: %v", err)
		}
		queue.LastAdditionDate = time.Now().Unix()
		queue.Requests = queue.Requests[1:]
		err = s.KSclient.Save(ctx, QUEUE, queue)
		s.Log(fmt.Sprintf("Ran the queue %v -> %v with %v", queueSize, len(queue.Requests), err))
		return err
	}

	s.Log(fmt.Sprintf("Still %v to go!", time.Now().Sub(time.Unix(queue.LastAdditionDate, 0))))

	return nil
}
