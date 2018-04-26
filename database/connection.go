package database

import (
	"os"

	"github.com/jackc/pgx"
	log "gopkg.in/inconshreveable/log15.v2"
)

var pool *pgx.ConnPool

func Connection() *pgx.Conn {

	config, err := pgx.ParseConnectionString(os.Getenv("DATABASE_URI"))
	if err != nil {
		log.Crit("Unable to parse environment: %v", err)
		os.Exit(1)
	}

	conn, err := pgx.Connect(config)
	if err != nil {
		log.Crit("Unable to connection to database: %v\n", err)
		os.Exit(1)
	}

	return conn
}

func PoolConnection() *pgx.ConnPool {

	//logger := log15adapter.NewLogger(log.New("module", "pgx"))

	var err error

	config, err := pgx.ParseConnectionString(os.Getenv("DATABASE_URI"))
	//config.Logger = logger

	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: 5,
		AfterConnect:   afterConnect,
	}
	pool, err = pgx.NewConnPool(connPoolConfig)
	if err != nil {
		log.Crit("Unable to create connection pool", "error", err)
		os.Exit(1)
	}

	return pool
}

// afterConnect creates the prepared statements that this application uses
func afterConnect(conn *pgx.Conn) (err error) {
	_, err = conn.Prepare("updateContent", `
		UPDATE sitemap_pages
		SET title = $1, h1 = $2, description = $3
		WHERE id = $4`)
	if err != nil {
		return
	}

	return
}
