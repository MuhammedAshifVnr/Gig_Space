package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	DBUrl     string
	GRPCPort  string
	JWTSecret string
}

func LoadConfig() error {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
	viper.AutomaticEnv()
	return nil
}
