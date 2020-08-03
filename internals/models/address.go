package models

import (
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/configs"
	"github.com/felipefrizzo/brazilian-zipcode-api/internals/helpers"
	"gopkg.in/mgo.v2/bson"
)

const (
	body string = `
		<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cli="http://cliente.bean.master.sigep.bsb.correios.com.br/">
			<soapenv:Header/>
			<soapenv:Body>
				<cli:consultaCEP>
					<cep>%s</cep>
				</cli:consultaCEP>
			</soapenv:Body>
		</soapenv:Envelope>
	`
	zipcodeNotFound string = "CEP NAO ENCONTRADO"
	zipcodeInvalid  string = "CEP INVÁLIDO"
)

// ErrAddressNotFound returns if address was not found
// ErrAddressInvalid returns if zip code sent is invalid
var (
	ErrAddressNotFound = errors.New("Address not found")
	ErrAddressInvalid  = errors.New("The informed zip code is invalid")
)

type (
	AddressInterface interface {
		Create(zipcode string) error
		IsUpdated() bool
		Update(zipcode string) error
	}

	soapResponse struct {
		Address Address `xml:"Body>consultaCEPResponse>return"`
	}

	soapResponseError struct {
		FaultError string `xml:"Body>Fault>faultstring"`
	}

	// Address struct has information from brazilian addresses
	Address struct {
		ID             bson.ObjectId `json:"-" bson:"_id,omitempty"`
		FederativeUnit string        `json:"federative_unit" bson:"federative_unit" xml:"uf"`
		City           string        `json:"city" bson:"city" xml:"cidade"`
		Neighborhood   string        `json:"neighborhood" bson:"neighborhood" xml:"bairro"`
		AddressName    string        `json:"address_name" bson:"address_name" xml:"end"`
		Complement     string        `json:"complement" bson:"complement" xml:"complemento2"`
		Zipcode        string        `json:"zipcode" bson:"zipcode" xml:"cep"`
		CreatedAt      time.Time     `json:"created_at" bson:"created_at"`
		UpdatedAt      time.Time     `json:"updated_at" bson:"updated_at"`
	}
)

// Create this function create address if doesn't exists
func (a *Address) Create(zipcode string) error {
	c, err := fetchCorreiosAddress(zipcode)
	if err != nil {
		log.Printf("ADDRESS_CREATE_ERROR - %v", err)
		return err
	}

	a.FederativeUnit = c.FederativeUnit
	a.City = c.City
	a.Neighborhood = c.Neighborhood
	a.AddressName = c.AddressName
	a.Complement = c.Complement
	a.Zipcode = c.Zipcode
	a.CreatedAt = time.Now().UTC()
	a.UpdatedAt = time.Now().UTC()

	return nil
}

// IsUpdated returns if the address is older than 7 days
func (a *Address) IsUpdated() bool {
	var sevenDaysAgo time.Time = time.Now().UTC().AddDate(0, 0, -7)
	return a.UpdatedAt.After(sevenDaysAgo)
}

// Update this function update address if outdated
func (a *Address) Update() error {
	c, err := fetchCorreiosAddress(a.Zipcode)
	if err != nil {
		log.Printf("ADDRESS_UPDATE_ERROR - %v", err)
		return err
	}

	a.FederativeUnit = c.FederativeUnit
	a.City = c.City
	a.Neighborhood = c.Neighborhood
	a.AddressName = c.AddressName
	a.Complement = c.Complement
	a.Zipcode = c.Zipcode
	a.UpdatedAt = time.Now().UTC()

	return nil
}

func fetchCorreiosAddress(zipcode string) (*Address, error) {
	var soap soapResponse
	var URL string = configs.Config.CorreiosURL

	response, err := helpers.RequestXML(URL, body, []interface{}{zipcode})
	if err != nil {
		return nil, fmt.Errorf("Error for fetching data from correios API. %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var soapError soapResponseError

		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Errorf("Error to read request body. %v", err)
		}

		xml.Unmarshal(helpers.ISO8859ToUTF8(data), &soapError)
		if soapError.FaultError == zipcodeNotFound {
			return nil, ErrAddressNotFound
		}

		if soapError.FaultError == zipcodeInvalid {
			return nil, ErrAddressInvalid
		}

		return nil, fmt.Errorf("Status code: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("Error to read request body. %v", err)
	}

	xml.Unmarshal(helpers.ISO8859ToUTF8(data), &soap)

	return &soap.Address, nil
}
