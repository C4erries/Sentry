package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	_ "github.com/lib/pq"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Run database migrations",
	Run: func(cmd *cobra.Command, args []string) {
		dsn := os.Getenv("POSTGRES_DSN")
		log.Printf("Running migrations on %s\n", dsn)
		runMigrations(dsn)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}

func runMigrations(dsn string) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("migrate: failed to open DB: %v", err)
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("migrate: failed to create driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", driver)
	if err != nil {
		log.Fatalf("migrate: failed to init: %v", err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatalf("migrate: up failed: %v", err)
	}

	log.Println("migrate: database is up to date")
}
