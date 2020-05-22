package models

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/configs"
	"github.com/felipefrizzo/brazilian-zipcode-api/internals/helpers"
	"gopkg.in/mgo.v2/bson"
)

const body string = `
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:cli="http://cliente.bean.master.sigep.bsb.correios.com.br/">
   <soapenv:Header/>
   <soapenv:Body>
      <cli:consultaCEP>
         <cep>%s</cep>
      </cli:consultaCEP>
   </soapenv:Body>
</soapenv:Envelope>
`

// SOAPResponse struct has information to parse xml return
type SOAPResponse struct {
	Address Address `xml:"Body>consultaCEPResponse>return"`
}

// Address struct has information from brazilian addresses
type Address struct {
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

// AddressIsUpdated this function is responsible to get the current address and validate if is updated
// based on UpdatedAt and Correios data
func (a *Address) AddressIsUpdated(zipcode string) error {
	var sevenDaysAgo time.Time = time.Now().UTC().AddDate(0, 0, -7)
	if a.Zipcode != "" && a.UpdatedAt.After(sevenDaysAgo) {
		return nil
	}

	c, err := fetchCorreiosAddress(zipcode)
	if err != nil {
		log.Printf("ADDRESS_VALIDATE_ERROR - Error for fetching data from correios API - %v", err)
		return fmt.Errorf("Error for fetching data from correios API - %v", err)
	}

	if a.Zipcode == "" {
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
	if a.UpdatedAt.Before(sevenDaysAgo) {
		a.FederativeUnit = c.FederativeUnit
		a.City = c.City
		a.Neighborhood = c.Neighborhood
		a.AddressName = c.AddressName
		a.Complement = c.Complement
		a.Zipcode = c.Zipcode
		a.UpdatedAt = time.Now().UTC()
	}

	return nil
}

func fetchCorreiosAddress(zipcode string) (*Address, error) {
	var soap SOAPResponse
	var URL string = configs.Config.CorreiosURL

	response, err := helpers.RequestXML(URL, body, []interface{}{zipcode})
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Status code: %v", response.StatusCode)
	}

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	xml.Unmarshal(helpers.ISO8859ToUTF8(data), &soap)

	return &soap.Address, nil
}
