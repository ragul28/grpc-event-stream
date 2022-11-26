package service

import (
	"context"
	"encoding/json"

	"github.com/nats-io/nats.go"
	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/model"
	"github.com/ragul28/grpc-event-stream/internal/order/repository"
	"google.golang.org/grpc/codes"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/status"
)

const (
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

func (s *Server) GetEvents(filter *pb.GetEventFilter, stream pb.Event_GetEventsServer) error {
	res, err := s.Repo.GetAllOrders(0, 10)
	if err != nil {
		return err
	}

	for _, ge := range res {
		if err := stream.Send(ge); err != nil {
			return err
		}
	}

	return nil
}

func (s *Server) Check(ctx context.Context, req *healthpb.HealthCheckRequest) (*healthpb.HealthCheckResponse, error) {
	return &healthpb.HealthCheckResponse{Status: healthpb.HealthCheckResponse_SERVING}, nil
}

func (s *Server) Watch(req *healthpb.HealthCheckRequest, ws healthpb.Health_WatchServer) error {
	return status.Errorf(codes.Unimplemented, "health check via Watch not implemented")
}
