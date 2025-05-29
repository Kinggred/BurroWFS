package config

import (
	"log"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
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
	err := godotenv.Load(getEnvFile("./.env"))
	if err != nil {
		log.Printf("Warning: could not load .env file: %v", err)
	}

	err = env.Parse(&CONFIG)
	if err != nil {
		log.Fatalf("Error loading environment variables: %v", err)
	}
}

func getEnvFile(filePath string) string {
	localPath := "./.env.local"
	if !fileCheck(filePath) && filePath != localPath {
		return getEnvFile(localPath)
	}
	return filePath
}

func fileCheck(filePath string) bool {
	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	if err != nil {
		log.Fatal(err)
		return false
	}
	return !info.IsDir()
}
