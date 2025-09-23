package config

import (
	"crypto/sha256"
	"fmt"
	"log"
	"os"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
)

type Config struct {
	Port             string `env:"PORT" envDefault:"80"`
	APIPrefix        string `env:"API_PREFIX" envDefault:"/api"`
	Debug            bool   `env:"DEBUG" envDefault:"false"`
	Secret           string `env:"SECRET,required"`
	PostgresHost     string `env:"POSTGRES_HOST,required"`
	PostgresPort     int    `env:"POSTGRES_PORT" envDefault:"6432"`
	PostgresUser     string `env:"POSTGRES_USER,required"`
	PostgresPassword string `env:"POSTGRES_PASSWORD,required"`
	PostgresDB       string `env:"POSTGRES_DB,required"`
	PostgresDSN      string `env:"POSTGRES_DSN" envDefault:""`

	// AWS S3 Configuration
	AWSRegion          string `env:"AWS_REGION" envDefault:"us-east-1"`
	AWSBucket          string `env:"AWS_BUCKET,required"`
	AWSAccessKeyID     string `env:"AWS_ACCESS_KEY_ID,required"`
	AWSSecretAccessKey string `env:"AWS_SECRET_ACCESS_KEY,required"`

	// Computed values
	HashedSecret []byte
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

	CONFIG.HashedSecret = hashSecret(CONFIG.Secret)
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

func hashSecret(secret string) []byte {
	encodedSecret := sha256.Sum256([]byte(secret))
	return encodedSecret[:16]
}
