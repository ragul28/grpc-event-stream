package grpcutil

import (
	"log"

	"github.com/ragul28/grpc-event-stream/pkg/getenv"
	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

func GetgrpcConn() *grpc.ClientConn {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(getenv.GetEnv("ORDER_GRPC_ADDR", address), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn
}
