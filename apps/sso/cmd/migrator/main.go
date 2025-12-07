package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	connString, migrationsPath, migrationsTable := fetchMigratorPaths()

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("%s?%s", connString, migrationsTable),
	)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Printf("no migrations to apply")

			return
		}

		panic(err)
	}
}

// fetchMigratorPaths fetches the paths for the migration, migrations table and env for connString.
// Priority: flag > env > default
// connString and migrationPath cannot be empty
// Default value: migrationPath: , migrationsTable: "migrations"
func fetchMigratorPaths() (string, string, string) {
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations directory")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("name"),
		os.Getenv("sslmode"),
	)

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}

	return connString, migrationsPath, migrationsTable
}
