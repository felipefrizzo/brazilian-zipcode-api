package configs

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

type mongoConfig struct {
	MongoURI string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
}

type config struct {
	Mongo       mongoConfig
	CorreiosURL string
}

// Config values collect from environment variables
var Config config

func init() {
	mongo := mongoConfig{}
	if err := env.Parse(&mongo); err != nil {
		panic(fmt.Errorf("Error to collect mongo environment values. Error: %s", err))
	}

	Config.Mongo = mongo
	Config.CorreiosURL = "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl"
}
