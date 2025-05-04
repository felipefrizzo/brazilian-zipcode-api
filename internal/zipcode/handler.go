package zipcode

import (
	"encoding/json"
	"net/http"

	"github.com/felipefrizzo/brazilian-zipcode-api/internal/address"
	"github.com/julienschmidt/httprouter"
)

// Handler some description
type Handler struct {
	Address address.AddressRepository
}

// New creates a new instance of Handler
func New(addr address.AddressRepository) *Handler {
	return &Handler{
		Address: addr,
	}
}

// AddHandlers some description
func (h *Handler) AddHandlers(r *httprouter.Router) {
	r.GET("/zipcode/:zipcode", h.FetchAddressByZipcode())
}

// FetchAddressByZipcode fetches an address by zipcode
func (h *Handler) FetchAddressByZipcode() httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		zipcode := params.ByName("zipcode")
		if zipcode == "" {
			http.Error(w, "missing zipcode", http.StatusBadRequest)
			return
		}

		addr, err := h.Address.Get(r.Context(), zipcode)
		if err != nil {
			http.Error(w, "address not found", http.StatusNotFound)
			return
		}

		response, err := json.Marshal(addr)
		if err != nil {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
