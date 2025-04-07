package main

import (
	"embed"
	"errors"
	"fmt"
	"log"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/joho/godotenv"
)

//go:embed migrations/*.sql
var migrationsFS embed.FS

func main() {
	// Load environment variables from .env if it exists
	err := godotenv.Load()
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Create database URL
	dbURL := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD")),
		Host:   fmt.Sprintf("%s:%s", os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT")),
		Path:   os.Getenv("POSTGRES_DB"),
	}

	q := dbURL.Query()
	q.Add("sslmode", os.Getenv("POSTGRES_SSLMODE"))
	dbURL.RawQuery = q.Encode()

	// Initialize migration source
	d, err := iofs.New(migrationsFS, "migrations")
	if err != nil {
		log.Fatalf("Failed to create migration source: %v", err)
	}

	// Create migrator instance
	m, err := migrate.NewWithSourceInstance("iofs", d, dbURL.String())
	if err != nil {
		log.Fatalf("Migration initialization failed: %v", err)
	}
	defer m.Close()

	// Run migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		log.Fatalf("Failed to get migration version: %v", err)
	}

	log.Printf("Database migrated successfully to version %d (dirty=%v)", version, dirty)
}
