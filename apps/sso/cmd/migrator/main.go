package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net"
	"net/url"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	log := slog.New(
		slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
	)

	connString, migrationsPath, migrationsTable := fetchMigratorPaths()

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("%s&x-migrations-table=%s", connString, migrationsTable),
	)
	if err != nil {
		log.Error("failed to create migrator: %v", slog.Any("error", err))
	}
	defer m.Close()

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("no migrations to apply")
			return
		}
		log.Error("failed to apply migrations: %v", slog.Any("error", err))
	}

	log.Info("migrations applied successfully")
}

// fetchMigratorPaths fetches the paths for the migration, migrations table and env for connString.
// Priority: flag > env > default
// connString and migrationPath cannot be empty
// Default value: migrationPath: , migrationsTable: "migrations".
func fetchMigratorPaths() (string, string, string) {
	var migrationsPath, migrationsTable string

	flag.StringVar(&migrationsPath, "migrations-path", "", "path to migrations directory")
	flag.StringVar(&migrationsTable, "migrations-table", "migrations", "name of migrations table")
	flag.Parse()

	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}
	if migrationsPath == "" {
		log.Fatal("migrations-path is required: set flag --migrations-path or env MIGRATIONS_PATH")
	}

	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	if dbUser == "" || dbPassword == "" || dbHost == "" || dbPort == "" || dbName == "" {
		log.Fatal("database credentials are required: DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME")
	}

	if dbSSLMode == "" {
		dbSSLMode = "disable"
	}

	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(dbUser, dbPassword),
		Host:   net.JoinHostPort(dbHost, dbPort),
		Path:   dbName,
	}

	q := u.Query()
	q.Set("sslmode", dbSSLMode)
	u.RawQuery = q.Encode()

	connString := u.String()

	return connString, migrationsPath, migrationsTable
}
