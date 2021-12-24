package zipcode

import (
	"fmt"
	"log"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Service represents zipcode application interface
type Service interface {
	FetchAddressByZipcode(zipcode string) (*models.Address, error)
}

type zipcode struct {
	MongoSession *mongo.Client
}

// New initialize new zipcode service
func New(session *mongo.Client) Service {
	return &zipcode{
		MongoSession: session,
	}
}

// FetchAddressByZipcode returns the address corresponding to the zipcode
func (z *zipcode) FetchAddressByZipcode(zipcode string) (*models.Address, error) {
	var address models.Address

	if err := z.MongoSession.Database("zipcode").Collection("addresses").FindOne(nil, bson.M{"zipcode": zipcode}).Decode(&address); err != nil {
		if err == mongo.ErrNoDocuments {
			if err := address.Create(zipcode); err != nil {
				return nil, err
			}
			insert, err := z.MongoSession.Database("zipcode").Collection("addresses").InsertOne(nil, address)
			if err != nil {
				log.Printf("Error inserting address %v", err)
				return nil, err
			}

			address.ID = insert.InsertedID.(primitive.ObjectID)
			return &address, nil
		}
		return nil, models.ErrAddressNotFound
	}

	if address.IsUpdated() {
		return &address, nil
	}

	if err := address.Update(); err != nil {
		z.MongoSession.Database("zipcode").Collection("addresses").DeleteOne(nil, bson.M{"_id": address.ID})
		return nil, err
	}

	if _, err := z.MongoSession.Database("zipcode").Collection("addresses").ReplaceOne(nil, address, bson.M{"upsert": true}); err != nil {
		errMessage := fmt.Sprintf("Error updating address %v", err)
		log.Println(errMessage)
		return nil, fmt.Errorf(errMessage)
	}

	return &address, nil
}
