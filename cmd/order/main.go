package main

import (
	"context"
	"log"
	"net"

	pb "github.com/ragul28/grpc-event-stream/event"
	psql "github.com/ragul28/grpc-event-stream/pkg/db"
	"github.com/ragul28/grpc-event-stream/pkg/repository"
	"google.golang.org/grpc"
)

const (
	port = ":50050"
)

type server struct {
	pb.UnimplementedEventServer
	repo repository.Repository
}

// CreateEvent creates a new Event
func (s *server) CreateEvent(ctx context.Context, in *pb.EventRequest) (*pb.EventResponse, error) {
	res, err := s.repo.CreateOrder(in)
	if err != nil {
		return nil, err
	}

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

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()

	server := &server{
		repo: repository,
	}
	pb.RegisterEventServer(s, server)

	log.Printf("gRPC Server listening on %v", port)
	s.Serve(lis)
}
