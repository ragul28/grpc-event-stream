package grpcutil

import (
	"log"

	"google.golang.org/grpc"
)

func GetgrpcConn(grpcAddr string) *grpc.ClientConn {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn
}
