package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// SOAPResponse struct has information to parse xml return
type SOAPResponse struct {
	Address Address `xml:"Body>consultaCEPResponse>return"`
}

// Address struct has information from brazilian addresses
type Address struct {
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	FederativeUnit string        `json:"federative_unit" bson:"federative_unit" xml:"uf"`
	City           string        `json:"city" bson:"city" xml:"cidade"`
	Neighborhood   string        `json:"neighborhood" bson:"neighborhood" xml:"bairro"`
	AddressName    string        `json:"address_name" bson:"address_name" xml:"end"`
	Complement     string        `json:"complement" bson:"complement" xml:"complemento2"`
	Zipcode        int32         `json:"zipcode" bson:"zipcode" xml:"cep"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
}
