package main

import (
	"burrowfs/config"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	config.LoadVariables()

	log.Println(config.CONFIG.PostgresDSN)
	m, err := migrate.New("file://db/migrate", config.CONFIG.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Println("Migration successful")
}
