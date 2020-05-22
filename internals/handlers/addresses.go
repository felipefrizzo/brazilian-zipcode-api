package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/middleware"
	"github.com/felipefrizzo/brazilian-zipcode-api/internals/models"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

// ZipcodeHandler function for returns the address corresponding to the zipcode
func ZipcodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	var params map[string]string = mux.Vars(r)
	var address models.Address
	var zipcode string = params["zipcode"]

	session, err := middleware.MongoConnection()
	if err != nil {
		log.Printf("ZIPCODE_HANDLER_ERROR - Error to close connection with MongoDB - %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer session.Close()

	session.DB("zipcode").C("addresses").Find(bson.M{"zipcode": zipcode}).One(&address)
	if err := address.AddressIsUpdated(zipcode); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	session.DB("zipcode").C("addresses").Upsert(bson.M{"zipcode": zipcode}, address)

	json.NewEncoder(w).Encode(address)
}
