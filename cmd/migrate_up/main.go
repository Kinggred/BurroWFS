package main

import (
	"burrowfs/core/config"
	"errors"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log.Println("Running migrate-tool: migration entrypoint reached")

	// Print current working directory
	if cwd, err := os.Getwd(); err == nil {
		log.Println("Current working directory:", cwd)
	} else {
		log.Println("Could not get working directory:", err)
	}

	// Print migration directory contents
	files, err := os.ReadDir("core/db/migrate")
	if err != nil {
		log.Fatalf("Could not read migration directory: %v", err)
	}
	log.Println("Migration files:")
	for _, f := range files {
		log.Println(" -", f.Name())
	}

	config.LoadVariables()

	log.Println("PostgresDSN:", config.CONFIG.PostgresDSN)
	m, err := migrate.New("file://core/db/migrate", config.CONFIG.PostgresDSN)
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Println("Migration successful")
}
