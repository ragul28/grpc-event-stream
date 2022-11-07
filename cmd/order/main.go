package main

import (
	"context"
	"encoding/json"
	"log"
	"net"

	"github.com/nats-io/nats.go"
	pb "github.com/ragul28/grpc-event-stream/event"
	psql "github.com/ragul28/grpc-event-stream/pkg/db"
	"github.com/ragul28/grpc-event-stream/pkg/repository"
	"github.com/ragul28/grpc-event-stream/pkg/stream"
	"google.golang.org/grpc"
)

const (
	port       = ":50050"
	streamName = "ORDERS"
)

type server struct {
	pb.UnimplementedEventServer
	repo repository.Repository
	nats nats.JetStreamContext
}

// CreateEvent creates a new Event
func (s *server) CreateEvent(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	// Store the order on DB
	res, err := s.repo.CreateOrder(in)
	if err != nil {
		return nil, err
	}

	values := map[string]string{"id": in.Id, "name": in.Name}
	jsonValue, _ := json.Marshal(values)

	// Publish order on nats jetstream
	s.nats.Publish(streamName, jsonValue)
	return res, nil
}

func (s *server) GetEvent(ctx context.Context, filter *pb.GetEventFilter) (*pb.GetEventResponse, error) {
	res, err := s.repo.GetOrder(filter)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func main() {
	db, err := psql.CreateConnection()
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("DB Successfully connected!")
	repository := &repository.OrderRepository{DB: db}

	js, err := stream.JetStreamInit(streamName)
	if err != nil {
		log.Println(err)
		return
	}

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()

	server := &server{
		repo: repository,
		nats: js,
	}
	pb.RegisterEventServer(s, server)

	log.Printf("gRPC Server listening on %v", port)
	s.Serve(lis)
}
