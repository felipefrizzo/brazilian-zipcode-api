package server

import (
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address/redis"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/config"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/correios"
	"github.com/felipefrizzo/brazilian-zipcode-api/internal/zipcode"
	"github.com/julienschmidt/httprouter"
	"github.com/richardwilkes/toolbox/errs"
)

type Server struct {
	Router *httprouter.Router
	Server *http.Server
}

// Handler represents a gateway handler
type Handler interface {
	AddHandlers(*httprouter.Router)
}

type handlerServices struct {
	address  address.AddressRepository
	correios correios.Correios
}

// New creates a new server instance
func New(logger *slog.Logger, cfg *config.Config) (*Server, error) {
	logger.Info("starting up server")

	services, err := createServices(cfg)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	router := httprouter.New()
	router.GET("/health", healthcheckHandler)
	router.PanicHandler = panicHandler
	server := &http.Server{
		Addr:    ":" + cfg.ServerPort,
		Handler: router,
	}

	srvr := &Server{
		Router: router,
		Server: server,
	}

	srvr.registerHandlers(services)
	return srvr, nil
}

func (srvr *Server) AddHandlers(handlers ...Handler) {
	for _, h := range handlers {
		h.AddHandlers(srvr.Router)
	}
}

func createServices(cfg *config.Config) (*handlerServices, error) {
	var addrService address.AddressRepository

	correiosService := correios.New(cfg.CorreiosURL)
	switch cfg.DatabaseDriver {
	case "redis":
		db, err := strconv.ParseInt(cfg.DatabaseName, 10, 32) // database name is translated to int as database number in redis
		if err != nil {
			return nil, errs.Wrap(err)
		}

		addrService, err = redis.NewClient(
			fmt.Sprintf("%s:%s", cfg.DatabaseURL, cfg.DatabasePort),
			cfg.DatabaseUsername,
			cfg.DatabasePassword,
			int(db),
			correiosService,
			cfg.CacheTTL,
		)
		if err != nil {
			return nil, errs.Wrap(err)
		}
	default:
		return nil, errs.Newf("unsupported database driver: %s", cfg.DatabaseDriver)
	}

	return &handlerServices{
		address:  addrService,
		correios: correiosService,
	}, nil
}

func (srvr *Server) registerHandlers(svcs *handlerServices) {
	zipcodeHandler := zipcode.New(svcs.address)

	srvr.AddHandlers(zipcodeHandler)
}

func healthcheckHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func panicHandler(w http.ResponseWriter, req *http.Request, rcv any) {
	w.WriteHeader(http.StatusInternalServerError)
}
