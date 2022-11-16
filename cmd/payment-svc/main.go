package main

import (
	"log"
	"net"

	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/payment"
	"github.com/ragul28/grpc-event-stream/pkg/getenv"
	"github.com/ragul28/grpc-event-stream/pkg/stream"
	"google.golang.org/grpc"
)

const (
	port       = ":50051"
	streamName = "ORDERS"
)

func main() {

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

	server := &payment.Server{
		Nats: js,
	}
	pb.RegisterEventServer(s, server)

	server.ConsumeEvent(js)

	log.Printf("gRPC Server listening on %v", ":"+port)
	s.Serve(lis)
}
