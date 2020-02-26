package server

import (
	"log"
	"net/http"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/handlers"
	"github.com/gorilla/mux"
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
	s.Router.HandleFunc("/zipcode/{zipcode:[0-9]+}", handlers.ZipcodeHandler).Methods("GET")
}

// Run the app on it's router
func (s *Server) Run(host string) {
	log.Fatal(http.ListenAndServe(host, s.Router))
}
