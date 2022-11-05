package main

import (
	"io"
	"log"

	"golang.org/x/net/context"

	pb "github.com/ragul28/grpc-event-stream/event"
)

const (
	address = "localhost:5000"
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
