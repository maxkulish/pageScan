package database

import (
	"fmt"
	"os"

	"github.com/jackc/pgx"
)

func Connection() *pgx.Conn {

	dbURI := os.Getenv("DATABASE_URI")

	config, err := pgx.ParseConnectionString(dbURI)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to parse environment:", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}
