package handlers

import (
	"encoding/json"
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
	var params = mux.Vars(r)
	var address models.Address

	zipcode, err := strconv.Atoi(params["zipcode"])
	if err != nil {
		panic(err)
	}

	session, err := middleware.MongoConnection()
	if err != nil {
		panic(err)
	}

	session.DB("zipcode").C("addresses").Find(bson.M{"zipcode": zipcode}).One(&address)
	json.NewEncoder(w).Encode(address)
}
