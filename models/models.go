package models

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/loitd/vabackend/config"
	_ "gopkg.in/goracle.v2"
)

// the `DBConn` struct will implement the ImportItf interface
// It also takes the sql DB connection object ~ database connection
type DBConn struct {
	DB *sql.DB
}

func InitDB(cfg *config.Config) {
	// Please remember that sql.Open only validate argument without creating connection. Ping/PingContext do that.
	db, err := sql.Open(cfg.DATABASE_DRIVER, cfg.DATABASE_URL)
	if err != nil {
		// log.Fatal(err)
		log.Println(err) // continue for dev only. Should be Fatal in production mode
	}
	defer db.Close()
	// Create database connection with context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err = db.PingContext(ctx)
	if err != nil {
		// log.Fatal(err)
		log.Println(err) //should be Fatal in production mode
	}
	// Pass the database connection to the interfaces (many interfaces need many assignments)
	ImportItf = &DBConn{DB: db}
}
