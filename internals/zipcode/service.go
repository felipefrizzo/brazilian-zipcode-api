package zipcode

import (
	"fmt"
	"log"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/models"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// Service represents zipcode application interface
type Service interface {
	FetchAddressByZipcode(zipcode string) (*models.Address, error)
}

type zipcode struct {
	MongoSession *mgo.Session
}

// New initialize new zipcode service
func New(session *mgo.Session) Service {
	return &zipcode{
		MongoSession: session,
	}
}

// FetchAddressByZipcode returns the address corresponding to the zipcode
func (z *zipcode) FetchAddressByZipcode(zipcode string) (*models.Address, error) {
	var address models.Address

	if err := z.MongoSession.DB("zipcode").C("addresses").Find(bson.M{"zipcode": zipcode}).One(&address); err != nil {
		if err == mgo.ErrNotFound {
			if err := address.Create(zipcode); err != nil {
				return nil, err
			}
			if err := z.MongoSession.DB("zipcode").C("addresses").Insert(address); err != nil {
				log.Printf("FETCH_ADDRESS_BY_ZIPCODE_ERROR - Error to create address - %v", err)
				return nil, fmt.Errorf("Error to create address. %v", err)
			}

			return &address, nil
		}

		return nil, models.ErrAddressNotFound
	}

	if !address.IsUpdated() {
		if err := address.Update(); err != nil {
			return nil, err
		}

		if _, err := z.MongoSession.DB("zipcode").C("addresses").Upsert(bson.M{"zipcode": zipcode}, address); err != nil {
			log.Printf("FETCH_ADDRESS_BY_ZIPCODE_ERROR - Error to update address - %v", err)
			return nil, fmt.Errorf("Error to update address. %v", err)
		}
	}

	return &address, nil
}
