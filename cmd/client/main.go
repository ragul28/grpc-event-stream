package main

import (
	"io"
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"

	pb "github.com/ragul28/grpc-event-stream/event"
)

const (
	address = "localhost:50050"
)

// createEvent calls the RPC method CreateEvent of EventServer
func createEvent(client pb.EventClient, Event *pb.EventRequest) {
	resp, err := client.CreateEvent(context.Background(), Event)
	if err != nil {
		log.Fatalf("Could not create Event: %v", err)
	}
	if resp.Success {
		log.Printf("A new Event has been added with id: %s", resp.Id)
	}
}

// getEvents calls the RPC method GetEvents of EventServer
func getEvents(client pb.EventClient, filter *pb.GetEventFilter) {
	// calling the streaming API
	stream, err := client.GetEvent(context.Background(), filter)
	if err != nil {
		log.Fatalf("Error on get Events: %v", err)
	}
	for {
		Event, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.GetEvents(_) = _, %v", client, err)
		}
		log.Printf("Event: %v", Event)
	}
}

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	// Creates a new CustomerClient
	client := pb.NewEventClient(conn)

	event := &pb.EventRequest{
		Id:   "cc0c1003-4461-4b96-b5d6-5e111ba441f7",
		Name: "order1",
	}

	// Create a new event
	createEvent(client, event)

	event = &pb.EventRequest{
		Id:   "3ec8b696-1e45-4622-9c32-30e8b3d91ff0",
		Name: "order2",
	}

	// Create a new event
	createEvent(client, event)

	// Filter with an empty Keyword
	filter := &pb.GetEventFilter{Id: "cc0c1003-4461-4b96-b5d6-5e111ba441f7"}
	getEvents(client, filter)
}
