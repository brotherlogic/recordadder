package main

import (
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
		err = s.rc.addRecord(ctx, queue.Requests[0])
		if err != nil {
			return err
		}
		queue.Requests = queue.Requests[1:]
		err = s.KSclient.Save(ctx, QUEUE, queue)
		return err
	}

	return nil
}
