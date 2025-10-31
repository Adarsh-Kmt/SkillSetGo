package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

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

    // Try to load .env if present (local dev), but don't fail if missing
    _ = godotenv.Load() // Looks for ".env" in current dir by default

    // Read values from environment variables
    config.password = os.Getenv("DB_PASSWORD")
    config.username = os.Getenv("DB_USERNAME")
    config.port = os.Getenv("DB_PORT")
    config.host = os.Getenv("DB_HOST")
    config.database = os.Getenv("DB_DATABASE")

    logger.Println("Loaded database configuration:")
    logger.Println("host:", config.host)
    logger.Println("port:", config.port)
    logger.Println("database:", config.database)
    logger.Println("username:", config.username)
    // ⚠️ Do NOT log password in real deployments

    return config, nil
}


func PostgresDBClientInit() error {

	config, err := postgresConfiguration()

	if err != nil {
		return err
	}

	connString := os.Getenv("DATABASE_URL")
	ctx := context.Background()
	connPool, err = pgxpool.New(ctx, connString)

	if err != nil {
		return err
	}

	Client = sqlc.New(connPool)
	if PostgresDB, err = sql.Open("postgres", connString); err != nil {
		return err
	}
	return nil
}
