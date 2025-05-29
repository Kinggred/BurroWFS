package config

import (
	"log"

	"github.com/caarlos0/env/v9"
)

type Config struct {
	Debug            bool   `env:"DEBUG" envDefault:"false"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"5432"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
}

var CONFIG Config

func LoadVariables() {
	err := env.Parse(&CONFIG)
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
}
