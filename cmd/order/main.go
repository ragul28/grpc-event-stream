package main

import (
	"context"
	"log"
	"net"
	"strings"

	pb "github.com/ragul28/grpc-event-stream/event"
	"google.golang.org/grpc"
)

const (
	port = ":50050"
)

type server struct {
	savedEvent []*pb.EventRequest
	pb.UnimplementedEventServer
}

// CreateEvent creates a new Event
func (s *server) CreateEvent(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	s.savedEvent = append(s.savedEvent, in)
	return &pb.EventResponse{Id: in.Id, Success: true}, nil
}

// GetEvents returns all Events by given filter
func (s *server) GetEvent(filter *pb.GetEventFilter, stream pb.Event_GetEventServer) error {
	for _, Event := range s.savedEvent {
		if filter.Id != "" {
			if !strings.Contains(Event.Id, filter.Id) {
				continue
			}
		}
		if err := stream.Send(Event); err != nil {
			return err
		}
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()
	pb.RegisterEventServer(s, &server{})
	s.Serve(lis)
}
