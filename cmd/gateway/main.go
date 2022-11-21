package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ragul28/grpc-event-stream/internal/gateway/router"
	"github.com/ragul28/grpc-event-stream/pkg/utils"
)

func main() {

	rt := mux.NewRouter()

	router.InitRouter(rt)

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Server started on Port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, rt))
}
