package server

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/felipefrizzo/brazilian-zipcode-api/handlers"
)

// Server struct has router and db instances
type Server struct {
	Router *mux.Router
}

// Initialize initializer server with predefined configuration
func (s *Server) Initialize() {
	s.Router = mux.NewRouter()
	s.setRouters()
}

// Set all required routers
func (s *Server) setRouters() {
	s.Router.HandleFunc("/", handlers.HomeHandler).Methods("GET")
}

// Run the app on it's router
func (s *Server) Run(host string) {
	log.Fatal(http.ListenAndServe(host, s.Router))
}
