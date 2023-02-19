package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

const configFileName = "config.toml"

type Config struct {
	BaseApiUrl      string   `mapstructure:"base_api_url" validate:"required,url"`
	BaseAuthUrl     string   `mapstructure:"base_auth_url" validate:"required,url"`
	ClientId        string   `mapstructure:"client_id" validate:"required"`
	ClientSecret    string   `mapstructure:"client_secret" validate:"required"`
	FallbackAuthUrl string   `mapstructure:"token_generator_fallback" validate:"required,url"`
	AllowedOrigins  []string `mapstructure:"allowed_origins" validate:"required"`
}

func GetConfig() (c Config) {
	viper.SetConfigFile(configFileName)
	viper.AddConfigPath(".")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("config file %s not found", configFileName)
	}
	if err := viper.Unmarshal(&c); err != nil {
		log.Fatalf("config file %s not found", configFileName)
	}
	validate := validator.New()
	if err := validate.Struct(&c); err != nil {
		log.Fatalf(err.Error())
	}
	return
}
