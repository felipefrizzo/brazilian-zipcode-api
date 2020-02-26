package middleware

import (
	"time"

	"github.com/felipefrizzo/brazilian-zipcode-api/internals/configs"
	"gopkg.in/mgo.v2"
)

// MongoConnection open connection with mongodb
func MongoConnection() (*mgo.Session, error) {
	config := configs.Config.Mongo

	mongo := &mgo.DialInfo{
		Addrs:    []string{config.Host},
		Timeout:  60 * time.Second,
		Database: config.Database,
		Username: config.Username,
		Password: config.Password,
	}
	return mgo.DialWithInfo(mongo)
}
