package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/middleware"
	"github.com/felipefrizzo/brazilian-zipcode-api/internals/models"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// ZipcodeHandler function for returns the address corresponding to the zipcode
func ZipcodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var params map[string]string = mux.Vars(r)
	var address models.Address

	zipcode, err := strconv.Atoi(params["zipcode"])
	if err != nil {
		log.Printf("ZIPCODE_HANDLER_ERROR - Error to convert the zipcode parameter - %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session, err := middleware.MongoConnection()
	if err != nil {
		log.Printf("ZIPCODE_HANDLER_ERROR - Error to close connection with MongoDB - %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	session.DB("zipcode").C("addresses").Find(bson.M{"zipcode": zipcode}).One(&address)
	err = address.AddressIsUpdated(zipcode)
	if err != nil {
		log.Printf("ZIPCODE_HANDLER_ERROR - Invalid address - %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	session.DB("zipcode").C("addresses").Upsert(bson.M{"zipcode": zipcode}, address)

	json.NewEncoder(w).Encode(address)
}
