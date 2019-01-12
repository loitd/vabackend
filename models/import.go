package models

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
)

// The data strcuture get from database
type ImportBatch struct {
	id                  int
	bank_code           string
	quantity            string
	batch_code          string
	file_name_root      string
	parent_account_epay string
}

// each method returns errors in case something went wrong
// Each method in this interface (Import) with attached parameter is dbconn - the DB Connection
type ImportInterface interface {
	GetImportBatchInfo(batchID int) (*ImportBatch, error)
	InsertAccount(va_number string, bank_code string, batch_id int, batch_code string, parent_account_epay string) error
	ImportAccount(batch_id int) error
}

// Make this interface global variable for application throughput. All package can access this interface
// With accessing this interface -> can call all methods in this interface (with attached dbconn - and they dont have to care about dbconn)
var ImportItf ImportInterface

func (dbconn *DBConn) GetImportBatchInfo(batchID int) (*ImportBatch, error) {
	// get import batch info based on batchID input
	row := dbconn.DB.QueryRow("SELECT ID, BANK_CODE, QUANTITY, BATCH_CODE, FILE_NAME_ROOT, PARENT_ACCOUNT_EP FROM TBL_IMPORT_BATCH WHERE ID = :1", batchID)
	// rows, err := dbconn.DB.Query("SELECT ID, BANK_CODE, QUANTITY, BATCH_CODE, FILE_NAME_ROOT, PARENT_ACCOUNT_EP FROM TBL_IMPORT_BATCH")
	importbatch := &ImportBatch{}
	err := row.Scan(&importbatch.id, &importbatch.bank_code, &importbatch.quantity, &importbatch.batch_code, &importbatch.file_name_root, &importbatch.parent_account_epay)
	if err != nil {
		return nil, err
	}
	log.Println(*importbatch)
	log.Println(importbatch.file_name_root)
	return importbatch, nil
}

func (dbconn *DBConn) InsertAccount(va_number string, bank_code string, batch_id int, batch_code string, parent_account_epay string) error {
	// curtime := time.Time()

	sql := `INSERT INTO TBL_VA_IMPORT 
			(VA_NUMBER, VA_NAME, BANK_CODE, BATCH_ID, STATUS, CREATED_BY,      CREATED_AT, UPDATED_AT, BATCH_CODE, COPPY,PARENT_ACCOUNT_EP)
	VALUES  (:1,        'EPAY',  :2,        :3,       0,      'Administrator', SYSDATE ,   SYSDATE,    :4,         0,    :5)`
	_, err := dbconn.DB.Query(sql, va_number, bank_code, batch_id, batch_code, parent_account_epay)
	if err != nil {
		log.Fatal(va_number, err)
		return err
	}
	log.Println("Done")
	return nil
}

func (dbconn *DBConn) ImportAccount(batch_id int) error {
	// first get the file name
	importbatch, err := ImportItf.GetImportBatchInfo(batch_id)
	if err != nil {
		log.Println(err)
		return err
	}
	//first read the file
	f, err := os.Open(importbatch.file_name_root)
	if err != nil {
		log.Println(err)
		return err
	}
	defer f.Close() //after error check, close this file
	lines := csv.NewReader(bufio.NewReader(f))
	lines.Comma = ' '
	lines.Comment = '#'
	curline := 1
	for {
		line, err := lines.Read()
		// stop atEOF
		if err == io.EOF {
			break
		}
		//ignore first line
		if curline == 1 {
			curline = curline + 1
			continue
		}
		// Remove the prefix
		account_number := line[0][1:]
		log.Println("Processing :1 - :2", curline, account_number)
		// insert to datbase
		err = ImportItf.InsertAccount(account_number, importbatch.bank_code, batch_id, importbatch.batch_code, importbatch.parent_account_epay)
		if err != nil {
			log.Println(err)
			return err
		}
		// increase curline afterall
		curline = curline + 1
	}
	return nil

}
