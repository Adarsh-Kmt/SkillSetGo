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

	err := godotenv.Load("/Users/adarsh-kmt/Desktop/Programming/Go/src/Personal/SkillSetGo/.env") // Specify the filename if it's named differently

	if err != nil {
		logger.Fatalf("Error loading .env file: %v", err)
	}

	config.password = os.Getenv("DB_PASSWORD")
	config.username = os.Getenv("DB_USERNAME")
	config.port = os.Getenv("DB_PORT")
	config.host = os.Getenv("DB_HOST")
	config.database = os.Getenv("DB_DATABASE")

	logger.Println("password : " + config.password)
	logger.Println("username : " + config.username)
	logger.Println("port : " + config.port)
	logger.Println("host : " + config.host)
	logger.Println("database : " + config.database)

	return config, nil

}

func PostgresDBClientInit() error {

	config, err := postgresConfiguration()

	if err != nil {
		return err
	}

	connString := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", config.username, config.password, config.host, config.port, config.database)
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
