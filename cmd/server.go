package server

import (
	"log"
	"net/http"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/middleware"
	"github.com/felipefrizzo/brazilian-zipcode-api/internals/zipcode"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2"
)

// Server struct has router and db instances
type Server struct {
	Router *mux.Router
	Mongo  *mgo.Session
}

// Initialize initializer server with predefined configuration
func (s *Server) Initialize() {
	s.Router = mux.NewRouter()
	s.Router.Use(mux.CORSMethodMiddleware(s.Router))

	s.InitializeMongo()
	s.InitializeZipcode()
}

// InitializeMongo initialize mongo db server
func (s *Server) InitializeMongo() {
	var err error

	s.Mongo, err = middleware.MongoConnection()
	if err != nil {
		log.Printf("ZIPCODE_GET_ERROR - Error to close connection with MongoDB - %v", err)
		panic("Error to close connection with MongoDB")
	}
}

// InitializeZipcode initialize zipcode service
func (s *Server) InitializeZipcode() {
	service := zipcode.New(s.Mongo)
	handler := zipcode.Handlers{
		Service: service,
	}
	handler.AddHandlers(s.Router)
}

// Run the app on it's router
func (s *Server) Run(host string) {
	log.Fatal(http.ListenAndServe(host, s.Router))
}
