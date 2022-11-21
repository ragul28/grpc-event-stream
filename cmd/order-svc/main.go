package main

import (
	"context"
	"log"
	"net"

	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/order/repository"
	order "github.com/ragul28/grpc-event-stream/internal/order/service"
	psql "github.com/ragul28/grpc-event-stream/pkg/db"
	"github.com/ragul28/grpc-event-stream/pkg/opentelemetry"
	"github.com/ragul28/grpc-event-stream/pkg/stream"
	"github.com/ragul28/grpc-event-stream/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
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

	tp := opentelemetry.InitTracerProvider()
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	port := utils.GetEnv("PORT", port)
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	// Creates a new gRPC server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()),
		grpc.StreamInterceptor(otelgrpc.StreamServerInterceptor()),
	)

	tp.Tracer("order-service")

	server := &order.Server{
		Repo: repository,
		Nats: js,
	}
	pb.RegisterEventServer(s, server)
	healthpb.RegisterHealthServer(s, server)

	log.Printf("gRPC Server listening on %v", ":"+port)
	s.Serve(lis)
}
