package models

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAddressIsUpdated(t *testing.T) {
	var a Address = Address{
		FederativeUnit: "PR",
		City:           "Cascavel",
		Neighborhood:   "Santa Felicidade",
		AddressName:    "Rua Major João Ribeiro Pinheiro",
		Complement:     "",
		Zipcode:        "85803260",
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	assert.True(t, a.IsUpdated())
}

func TestAddressIsNotUpdated(t *testing.T) {
	var a Address = Address{
		FederativeUnit: "PR",
		City:           "Cascavel",
		Neighborhood:   "Santa Felicidade",
		AddressName:    "Rua Major João Ribeiro Pinheiro",
		Complement:     "",
		Zipcode:        "85803260",
		CreatedAt:      time.Now().UTC().AddDate(0, 0, -8),
		UpdatedAt:      time.Now().UTC().AddDate(0, 0, -8),
	}

	assert.False(t, a.IsUpdated())
}
