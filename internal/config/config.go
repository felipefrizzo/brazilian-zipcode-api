package config

import (
	"os"
	"strings"

	"github.com/mitchellh/mapstructure"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/spf13/viper"
)

const (
	correiosURl = "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl"
)

// Config represents the application configuration.
type Config struct {
	ServerPort string `mapstructure:"port"`

	CorreiosURL string `mapstructure:"correios_url"`

	DatabaseDriver   string `mapstructure:"database_driver"`
	DatabaseURL      string `mapstructure:"database_url"`
	DatabasePort     string `mapstructure:"database_port"`
	DatabaseName     string `mapstructure:"database_name"`
	DatabaseUsername string `mapstructure:"database_username"`
	DatabasePassword string `mapstructure:"database_password"`
}

// New creates a new config instance.
func New() (*Config, error) {
	viper.SetEnvPrefix("APP")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, errs.Wrap(err)
		}
	}

	viper.SetDefault("port", "8080")
	viper.SetDefault("correios.url", correiosURl)

	var result map[string]any
	var config Config
	if err := viper.Unmarshal(&result); err != nil {
		return nil, errs.Wrap(err)
	}

	if err := mapstructure.Decode(result, &config); err != nil {
		return nil, errs.Wrap(err)
	}

	return &config, nil
}
