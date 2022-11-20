package service

import (
	"fmt"
	"log"

	"golang.org/x/net/context"

	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/model"
	"github.com/ragul28/grpc-event-stream/pkg/utils"
)

type Repository interface {
	CreateOrder(client pb.EventClient, Event *pb.EventRequest) error
	GetOrder(client pb.EventClient, filter *pb.GetEventFilter) (*model.OrderEvent, error)
}

type SrvRepository struct{}

// createEvent calls the RPC method CreateEvent of EventServer
func (s *SrvRepository) CreateOrder(client pb.EventClient, Event *pb.EventRequest) error {
	resp, err := client.CreateEvent(context.Background(), Event)
	if err != nil {
		return fmt.Errorf("Could not create Event: %v", err)
	}
	if resp.Success {
		log.Printf("Event added id: %s", resp.Id)
	}

	return nil
}

// getEvents calls the RPC method GetEvents of EventServer
func (s *SrvRepository) GetOrder(client pb.EventClient, filter *pb.GetEventFilter) (*model.OrderEvent, error) {
	// calling the streaming API
	resp, err := client.GetEvent(context.Background(), filter)
	if err != nil {
		return nil, fmt.Errorf("Error on get Event: %v", err)
	}
	log.Printf("Get Event id %s: %s", utils.Sanitize(filter.Id), resp.Name)

	return &model.OrderEvent{Id: resp.Id, Name: resp.Name}, nil
}
