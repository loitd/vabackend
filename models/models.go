package models

import (
	"database/sql"
	"log"

	"github.com/loitd/vabackend/config"
	_ "gopkg.in/goracle.v2"
)

// the `DBConn` struct will implement the ImportItf interface
// It also takes the sql DB connection object ~ database connection
type DBConn struct {
	DB *sql.DB
}

// The data strcuture get from database
type ImportBatch struct {
	id                  int
	bank_code           string
	quantity            string
	batch_code          string
	file_name_root      string
	parent_account_epay string
}

var ImportBatchVar ImportBatch

type ImportStatus struct {
	ResponseCode int    "json:ResponseCode"
	Message      string "json:Message"
	TotalSuccess int    "json:TotalSuccess"
	TotalError   int    "json:TotalError"
	TotalRecords int    "json:TotalRecords total records read from files"
}

// Store global status
var ImportStatusVar ImportStatus

// each method returns errors in case something went wrong
// Each method in this interface (Import) with attached parameter is dbconn - the DB Connection
type ImportInterface interface {
	GetImportBatchInfo(batchID int) (*ImportBatch, error)
	InsertAccount(va_number string, bank_code string, batch_id int, batch_code string, parent_account_epay string) error
	ImportAccount(batch_id int) error
	ImportAccountLogic(batch_id int) error
}

// Make this interface global variable for application throughput. All package can access this interface
// With accessing this interface -> can call all methods in this interface (with attached dbconn - and they dont have to care about dbconn)
var ImportItf ImportInterface

func Reset() {
	ImportStatusVar.ResponseCode = 200
	ImportStatusVar.Message = ""
	ImportStatusVar.TotalSuccess = 0
	ImportStatusVar.TotalError = 0
	ImportStatusVar.TotalRecords = 0
	//
	ImportBatchVar.id = nil
}

func InitDB(cfg *config.Config) {
	// Please remember that sql.Open only validate argument without creating connection. Ping/PingContext do that.
	db, err := sql.Open(cfg.DATABASE_DRIVER, cfg.DATABASE_URL)
	if err != nil {
		log.Fatal(err)
		// log.Println(err) // continue for dev only. Should be Fatal in production mode
	}
	// Create database connection with context. Must in v1.2 for faster startup
	// ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	// defer cancel()
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		// log.Println(err) //should be Fatal in production mode
	}

	// Pass the database connection to the interfaces (many interfaces need many assignments)
	ImportItf = &DBConn{DB: db}

	// Close database connection after using
	// defer db.Close()
}
