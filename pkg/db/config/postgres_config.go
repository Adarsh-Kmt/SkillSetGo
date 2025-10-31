package db

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/adarsh-kmt/skillsetgo/pkg/db/sqlc"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	Client     *sqlc.Queries
	PostgresDB *sql.DB
	connPool   *pgxpool.Pool
	logger     = log.New(os.Stdout, "SKILLSETGO DATABASE CONFIG >> ", 0)
)

type postgresConfig struct {
	host     string
	port     string
	username string
	password string
	database string
}

func postgresConfiguration() (*postgresConfig, error) {
	config := &postgresConfig{}

	// Try to load .env if present (for local dev)
	_ = godotenv.Load()

	config.password = os.Getenv("DB_PASSWORD")
	config.username = os.Getenv("DB_USERNAME")
	config.port = os.Getenv("DB_PORT")
	config.host = os.Getenv("DB_HOST")
	config.database = os.Getenv("DB_DATABASE")

	logger.Println("Loaded local database configuration:")
	logger.Println("  host:", config.host)
	logger.Println("  port:", config.port)
	logger.Println("  database:", config.database)
	logger.Println("  username:", config.username)

	return config, nil
}

func PostgresDBClientInit() error {
	ctx := context.Background()

	// Prefer DATABASE_URL (Render) if available
	connString := os.Getenv("DATABASE_URL")

	if connString == "" {
		// Fall back to manual config (local dev)
		config, err := postgresConfiguration()
		if err != nil {
			return err
		}

		connString = "postgresql://" + config.username + ":" + config.password +
			"@" + config.host + ":" + config.port + "/" + config.database + "?sslmode=require"
		logger.Println("Using local DB connection string.")
	} else {
		logger.Println("Using DATABASE_URL from environment.")
	}

	// Ensure SSL mode is set (Render requires this)
	if !strings.Contains(connString, "sslmode=") {
		if strings.Contains(connString, "?") {
			connString += "&sslmode=require"
		} else {
			connString += "?sslmode=require"
		}
	}

	// Create connection pool
	var err error
	connPool, err = pgxpool.New(ctx, connString)
	if err != nil {
		return err
	}

	Client = sqlc.New(connPool)

	PostgresDB, err = sql.Open("postgres", connString)
	if err != nil {
		return err
	}

	logger.Println("✅ Connected to PostgreSQL successfully.")
	return nil
}
