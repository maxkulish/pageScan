package database

import (
	"log"
	"os"

	"github.com/jackc/pgx"
)

func Connection() *pgx.Conn {

	dbURI := os.Getenv("DATABASE_URI")

	config, err := pgx.ParseConnectionString(dbURI)
	if err != nil {
		log.Printf("Unable to parse environment: %v", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		log.Printf("Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
