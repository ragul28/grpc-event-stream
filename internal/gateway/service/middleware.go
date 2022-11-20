package service

import (
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/ragul28/grpc-event-stream/pkg/utils"
)

type Middleware struct{}

// General noauth middleware
func (m *Middleware) Middleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
		log.Println(r.Method, r.RequestURI, time.Since(startTime))
	})
}

// allowCORS from any origin. TODO disable in prod
func (m *Middleware) AllowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

// preflightHandler adds headers and CROS from ACAM
func preflightHandler(w http.ResponseWriter, r *http.Request) {
	headers := []string{"Content-Type", "Accept"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
	log.Println("preflight request for", utils.Sanitize(r.URL.Path))
}
