package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Address struct has information from brazilian addresses
type Address struct {
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	FederativeUnit string        `json:"federative_unit" bson:"federative_unit"`
	City           string        `json:"city" bson:"city"`
	Neighborhood   string        `json:"neighborhood" bson:"neighborhood"`
	AddressName    string        `json:"address_name" bson:"address_name"`
	Complement     string        `json:"complement" bson:"complement"`
	Zipcode        int32         `json:"zipcode" bson:"zipcode"`
	CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
}
