package order

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/model"
	"github.com/ragul28/grpc-event-stream/internal/repository"
)

const (
	port               = ":50050"
	streamName         = "ORDERS"
	streamSubjectsname = "ORDERS.new"
)

type Server struct {
	pb.UnimplementedEventServer
	Repo repository.Repository
	Nats nats.JetStreamContext
}

// CreateEvent creates a new Event
func (s *Server) CreateEvent(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	// Store the order on DB
	res, err := s.Repo.CreateOrder(in)
	if err != nil {
		return nil, err
	}

	orderVal := model.OrderEvent{Id: in.Id, Name: in.Name}
	jsonValue, _ := json.Marshal(orderVal)

	// Publish order on nats jetstream
	s.Nats.Publish(streamSubjectsname, jsonValue)
	return res, nil
}

func (s *Server) GetEvent(ctx context.Context, filter *pb.GetEventFilter) (*pb.GetEventResponse, error) {
	res, err := s.Repo.GetOrder(filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}