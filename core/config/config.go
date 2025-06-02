package config

import (
	"fmt"
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
	PostgresDSN      string `env:"POSTGRES_DSN" envDefault:""`
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

	if CONFIG.PostgresDSN == "" {
		CONFIG.PostgresDSN = buildPostgresDSN(
			CONFIG.PostgresUser,
			CONFIG.PostgresPassword,
			CONFIG.PostgresHost,
			CONFIG.PostgresPort,
			CONFIG.PostgresDB,
		)
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

func buildPostgresDSN(user, password, host string, port int, dbname string) string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		user, password, host, port, dbname,
	)
}
