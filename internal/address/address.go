package address

import (
	"context"
	"time"
)

// Address struct has information from brazilian addresses
type Address struct {
	FederativeUnit string    `json:"federative_unit" bson:"federative_unit" xml:"uf"`
	City           string    `json:"city" bson:"city" xml:"cidade"`
	Neighborhood   string    `json:"neighborhood" bson:"neighborhood" xml:"bairro"`
	AddressName    string    `json:"address_name" bson:"address_name" xml:"end"`
	Complement     string    `json:"complement" bson:"complement" xml:"complemento2"`
	Zipcode        string    `json:"zipcode" bson:"zipcode" xml:"cep"`
	CreatedAt      time.Time `json:"created_at" bson:"created_at"`
}

// AddressRepository interface defines methods for interacting with address data.
type AddressRepository interface {
	Get(ctx context.Context, zipcode string) (*Address, error)
	Save(ctx context.Context, address *Address, zipcode string) error
	Delete(ctx context.Context, zipcode string) error
}
