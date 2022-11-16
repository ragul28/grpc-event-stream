package main

import (
	"log"
	"net"

	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/order"
	"github.com/ragul28/grpc-event-stream/internal/repository"
	psql "github.com/ragul28/grpc-event-stream/pkg/db"
	"github.com/ragul28/grpc-event-stream/pkg/getenv"
	"github.com/ragul28/grpc-event-stream/pkg/stream"
	"google.golang.org/grpc"
)

const (
	port       = ":50050"
	streamName = "ORDERS"
)

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

	port := getenv.GetEnv("PORT", port)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer()

	server := &order.Server{
		Repo: repository,
		Nats: js,
	}
	pb.RegisterEventServer(s, server)

	log.Printf("gRPC Server listening on %v", ":"+port)
	s.Serve(lis)
}
