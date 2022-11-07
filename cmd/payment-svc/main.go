package main

import (
	"encoding/json"
	"log"
	"net"

	"github.com/nats-io/nats.go"
	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/model"
	psql "github.com/ragul28/grpc-event-stream/pkg/db"
	"github.com/ragul28/grpc-event-stream/pkg/repository"
	"github.com/ragul28/grpc-event-stream/pkg/stream"
	"google.golang.org/grpc"
)

const (
	port               = ":50051"
	streamName         = "ORDERS"
	streamSubjectsname = "ORDERS.new"
)

type server struct {
	pb.UnimplementedEventServer
	repo repository.Repository
	nats nats.JetStreamContext
}

func consumeEvent(js nats.JetStreamContext) {
	_, err := js.Subscribe(streamSubjectsname, func(m *nats.Msg) {
		err := m.Ack()

		if err != nil {
			log.Println("Unable to Ack", err)
			return
		}

		var orderEvt model.OrderEvent

		err = json.Unmarshal(m.Data, &orderEvt)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Consumer  =>  Subject: %s  -  ID: %s  -  Name: %s\n", m.Subject, orderEvt.Id, orderEvt.Name)
	})

	if err != nil {
		log.Println("Subscribe failed")
		return
	}
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

	consumeEvent(js)

	log.Printf("gRPC Server listening on %v", port)
	s.Serve(lis)
}
