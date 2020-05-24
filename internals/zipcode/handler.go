package zipcode

import (
	"encoding/json"
	"net/http"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/models"
	"github.com/gorilla/mux"
)

// Handlers some description
type Handlers struct {
	Service Service
}

// AddHandlers some description
func (h *Handlers) AddHandlers(r *mux.Router) {
	r.HandleFunc("/zipcode/{zipcode:[0-9]+}", h.FetchAddressByZipcode).Methods("GET")
}

// FetchAddressByZipcode function for returns the address corresponding to the zipcode
func (h *Handlers) FetchAddressByZipcode(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var params map[string]string = mux.Vars(r)
	var zipcode string = params["zipcode"]

	address, err := h.Service.FetchAddressByZipcode(zipcode)
	if err != nil {
		if err == models.ErrAddressNotFound {
			http.Error(w, models.ErrAddressNotFound.Error(), http.StatusNotFound)
			return
		}
		if err == models.ErrAddressInvalid {
			http.Error(w, models.ErrAddressInvalid.Error(), http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(address)
}
