package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ragul28/grpc-event-stream/internal/gateway/router"
	"github.com/ragul28/grpc-event-stream/pkg/opentelemetry"
	"github.com/ragul28/grpc-event-stream/pkg/utils"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux"
)

func main() {

	tp, err := opentelemetry.InitTracer()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	rt := mux.NewRouter()
	// Enable otelmux
	rt.Use(otelmux.Middleware("gateway"))

	router.InitRouter(rt)

	port := utils.GetEnv("PORT", "8080")
	log.Printf("Server started on Port: %v", port)
	log.Fatal(http.ListenAndServe(":"+port, rt))
}
