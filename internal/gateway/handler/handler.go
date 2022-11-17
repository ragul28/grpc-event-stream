package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	pb "github.com/ragul28/grpc-event-stream/event"
	"github.com/ragul28/grpc-event-stream/internal/gateway/service"
	"github.com/ragul28/grpc-event-stream/internal/model"
	"github.com/ragul28/grpc-event-stream/pkg/getenv"

	"google.golang.org/grpc"
)

const (
	address = "localhost:8080"
)

type Handler struct {
	repository service.Repository
}

func NewHandler(sr service.Repository) *Handler {
	return &Handler{
		repository: sr,
	}
}

// HealthEndpoint for srv status
func (h *Handler) HealthEndpoint(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service Healthy"))
}

// Create Order
func (h *Handler) CreateOrder(w http.ResponseWriter, req *http.Request) {

	var order model.OrderEvent
	_ = json.NewDecoder(req.Body).Decode(&order)

	if order.Id == "" || order.Name == "" {
		err := "error missing id or name"
		log.Println(err)
		http.Error(w, err, http.StatusBadRequest)
		return
	}

	conn := getgrpcConn()
	defer conn.Close()

	client := pb.NewEventClient(conn)

	event := &pb.EventRequest{
		Id:   order.Id,
		Name: order.Name,
	}
	err := h.repository.CreateOrder(client, event)
	if err != nil {
		log.Panicln(err)
	}

	w.WriteHeader(http.StatusOK)
}

// Get Order
func (h *Handler) GetOrder(w http.ResponseWriter, req *http.Request) {

	reqId := mux.Vars(req)["id"]
	if reqId == "" {
		log.Println("Id cannot be empty")
	}
	log.Println("Get userid:", reqId)

	conn := getgrpcConn()
	defer conn.Close()

	client := pb.NewEventClient(conn)

	// Filter with an empty Keyword
	filter := &pb.GetEventFilter{Id: reqId}
	resOrder, err := h.repository.GetOrder(client, filter)
	if err != nil {
		log.Println(err)
	}

	json.NewEncoder(w).Encode(resOrder)
	w.WriteHeader(http.StatusOK)
}

func getgrpcConn() *grpc.ClientConn {
	// Set up a connection to the gRPC server.
	conn, err := grpc.Dial(getenv.GetEnv("ORDER_GRPC_ADDR", address), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return conn
}
