package config

import (
	"errors"
	"os"

	"github.com/spf13/viper"
)

type EnvVars struct {
	PORT              string `mapstructure:"PORT"`
	POSTGRES_PORT     string `mapstructure:"POSTGRES_PORT"`
	POSTGRES_HOST     string `mapstructure:"POSTGRES_HOST"`
	POSTGRES_NAME     string `mapstructure:"POSTGRES_NAME"`
	POSTGRES_USER     string `mapstructure:"POSTGRES_USER"`
	POSTGRES_PASSWORD string `mapstructure:"POSTGRES_PASSWORD"`
	POSTGRES_SSL      string `mapstructure:"POSTGRES_SSL"`
	AUTH0_DOMAIN      string `mapstructure:"AUTH0_DOMAIN"`
	AUTH0_AUDIENCE    string `mapstructure:"AUTH0_AUDIENCE"`
}

func LoadConfig() (config EnvVars, err error) {
	env := os.Getenv("GO_ENV")
	if env == "production" {
		return EnvVars{
			PORT:              os.Getenv("PORT"),
			POSTGRES_PORT:     os.Getenv("POSTGRES_PORT"),
			POSTGRES_HOST:     os.Getenv("POSTGRES_HOST"),
			POSTGRES_NAME:     os.Getenv("POSTGRES_NAME"),
			POSTGRES_USER:     os.Getenv("POSTGRES_USER"),
			POSTGRES_PASSWORD: os.Getenv("POSTGRES_PASSWORD"),
		}, nil
	}

	viper.AddConfigPath(".")
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)

	// validate config here
	if config.POSTGRES_HOST == "" {
		err = errors.New("POSTGRES_HOST is required")
		return
	}

	if config.POSTGRES_NAME == "" {
		err = errors.New("POSTGRES_NAME is required")
		return
	}

	return
}
