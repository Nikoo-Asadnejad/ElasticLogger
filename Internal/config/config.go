package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Elasticsearch struct {
		URL      string
		Username string
		Password string
		Index    string
	}
	RabbitMQ struct {
		ConnectionString string `mapstructure:"connection_string"`
	}
}

func Load() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatalf("Unable to decode config into struct: %v", err)
	}
	return &cfg
}
