package config

import (
	"os"

	"github.com/mitchellh/mapstructure"
	"github.com/richardwilkes/toolbox/errs"
	"github.com/spf13/viper"
)

const (
	defaultPort     = "8080"
	correiosURL     = "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl"
	defaultCacheTTL = "3600"
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

	CacheTTL string `mapstructure:"cache_ttl_seconds"`
}

// New creates a new config instance.
func New() (*Config, error) {
	viper.SetEnvPrefix("APP")
	viper.AutomaticEnv()

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		viper.SetConfigType("env")
		if err := viper.ReadInConfig(); err != nil {
			return nil, errs.Wrap(err)
		}
	}

	viper.SetDefault("port", defaultPort)
	viper.SetDefault("correios_url", correiosURL)
	viper.SetDefault("cache_ttl_seconds", defaultCacheTTL)

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
