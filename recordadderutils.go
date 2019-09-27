package main

import (
	"fmt"

	pb "github.com/brotherlogic/recordadder/proto"
	"golang.org/x/net/context"
)

func (s *Server) processQueue(ctx context.Context) error {
	data, _, err := s.KSclient.Read(ctx, QUEUE, &pb.Queue{})
	if err != nil {
		return err
	}
	queue := data.(*pb.Queue)

	if len(queue.Requests) > 0 {
		queueSize := len(queue.Requests)
		err = s.rc.addRecord(ctx, queue.Requests[0])
		if err != nil {
			return err
		}
		queue.Requests = queue.Requests[1:]
		err = s.KSclient.Save(ctx, QUEUE, queue)
		s.Log(fmt.Sprintf("Ran the queue %v -> %v with %v", queueSize, len(queue.Requests), err))
		return err
	}

	return nil
}
