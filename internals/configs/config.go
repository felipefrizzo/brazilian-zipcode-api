package configs

import (
	"fmt"

	"github.com/caarlos0/env"
)

type mongoConfig struct {
	Host     string `env:"MONGO_HOST" envDefault:"mongo"`
	Username string `env:"MONGO_USERNAME"`
	Password string `env:"MONGO_PASSWORD"`
	Database string `env:"MONGO_DATABASE"`
}

type config struct {
	Mongo mongoConfig
}

// Config values collect from environment variables
var Config config

func init() {
	mongo := mongoConfig{}
	err := env.Parse(&mongo)
	if err != nil {
		panic(fmt.Errorf("Error to collect mongo environment values. Error: %s", err))
	}

	Config.Mongo = mongo
}
