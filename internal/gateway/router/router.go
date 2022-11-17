package router

import (
	"github.com/gorilla/mux"
	"github.com/ragul28/grpc-event-stream/internal/gateway/handler"
	"github.com/ragul28/grpc-event-stream/internal/gateway/service"
)

func InitRouter(router *mux.Router) {

	repos := &service.SrvRepository{}
	mws := &service.Middleware{}

	h := handler.NewHandler(repos)

	r := router.PathPrefix("/api/gw").Subrouter()

	r.HandleFunc("/", h.HealthEndpoint).Methods("GET")
	r.HandleFunc("/order/create", mws.Middleware(h.CreateOrder)).Methods("POST")
	r.HandleFunc("/order/get/{id}", mws.Middleware(h.GetOrder)).Methods("GET")
}
