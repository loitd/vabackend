package models

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	config "github.com/loitd/vabackend/config"
)

func (dbconn *DBConn) GetImportBatchInfo(batchID int) (*ImportBatch, error) {
	// get import batch info based on batchID input
	// Ping first
	err := dbconn.DB.Ping()
	if err != nil {
		log.Println("DBPing: ", err)
		return nil, err
	}
	// now do Query
	row := dbconn.DB.QueryRow("SELECT ID, BANK_CODE, QUANTITY, BATCH_CODE, FILE_NAME_ROOT, PARENT_ACCOUNT_EP FROM TBL_IMPORT_BATCH WHERE ID = :1", batchID)
	// rows, err := dbconn.DB.Query("SELECT ID, BANK_CODE, QUANTITY, BATCH_CODE, FILE_NAME_ROOT, PARENT_ACCOUNT_EP FROM TBL_IMPORT_BATCH")
	importbatch := &ImportBatch{}
	err = row.Scan(&importbatch.id, &importbatch.bank_code, &importbatch.quantity, &importbatch.batch_code, &importbatch.file_name_root, &importbatch.parent_account_epay)
	if err != nil {
		return nil, err
	}
	// Update totalimport
	ImportStatusVar.TotalImport, _ = strconv.Atoi(importbatch.quantity)
	// log.Println(*importbatch)
	log.Println("Got file location: ", importbatch.file_name_root)
	return importbatch, nil
}

func (dbconn *DBConn) InsertAccount(va_number string, bank_code string, batch_id int, batch_code string, parent_account_epay string) error {
	// Ping first
	err := dbconn.DB.Ping()
	if err != nil {
		log.Println("DBPing: ", err)
		return err
	}
	// now do Query
	// Insert account into TBL_VA_IMPORT
	sql := `INSERT INTO TBL_VA_IMPORT 
			(VA_NUMBER, VA_NAME, BANK_CODE, BATCH_ID, STATUS, CREATED_BY,      CREATED_AT, UPDATED_AT, BATCH_CODE, COPPY,PARENT_ACCOUNT_EP)
	VALUES  (:1,        'EPAY',  :2,        :3,       0,      'Administrator', SYSDATE ,   SYSDATE,    :4,         0,    :5)`
	// Execute the query now
	_, err = dbconn.DB.Query(sql, va_number, bank_code, batch_id, batch_code, parent_account_epay)
	if err != nil {
		log.Fatal(va_number, err)
		config.LogFile("./fatal.log", fmt.Sprintf(":0-:1", va_number, err))
		return err
	}
	log.Println("InsertAccount Done")
	return nil
}

func parseFile(filename string, jobs chan string, wg *sync.WaitGroup) {
	// read the file and push to channel
	f, err := os.Open(filename)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	reader := csv.NewReader(f)
	reader.Comma = ' '
	reader.Comment = '#'
	curline := 1
	for {
		line, err := reader.Read()
		// stop atEOF
		if err == io.EOF {
			log.Println("Its EOF!")
			break
		}
		// capture any other errors
		if err != nil {
			log.Println("Reading imported file: ", err)
			break
		}
		//ignore first line
		if curline == 1 {
			curline = curline + 1
			continue
		}
		// Remove the prefix
		account_number := line[0][1:]
		// Pushing to channel
		jobs <- account_number
		log.Println(fmt.Sprintf("Pushed %d - %s", curline, account_number))
		ImportStatusVar.NoOfImport = curline
		// increase curline afterall
		curline = curline + 1
	}
	// Notify waitgroup its done
	wg.Done()
	// Must close jobs channel
	close(jobs)
	log.Println("parseFile done!")
}

func workerDB(id int, jobs chan string, errs chan string, wg *sync.WaitGroup, ib *ImportBatch) {
	// read the channel for VA_NUMBER, this is range over channel
	// This range iterates over each element as it’s received from queue.
	// Because we CLOSED the channel above, the iteration terminates after receiving the 2 elements.
	for vaNumber := range jobs {
		// va_number := <-job
		log.Println(fmt.Sprintf("workerDB %d processing va_number: %s|%s|%s|%s", id, vaNumber, ib.bank_code, ib.id, ib.batch_code))
		// insert to database
		err := ImportItf.InsertAccount(vaNumber, ib.bank_code, ib.id, ib.batch_code, ib.parent_account_epay)
		if err != nil {
			errs <- fmt.Sprintf("%s | %s", vaNumber, err.Error)
		}
	}
	// Return done when ALL job finished
	wg.Done()
	log.Println(fmt.Sprintf("workerDB %d done!", id))
}

func writeErr(errs chan string, wg2 *sync.WaitGroup) {
	// write error only to file
	// Start only 1 routine for this task
	for err := range errs {
		config.LogFile("fatal.log", err)
		ImportStatusVar.NoOfFail = ImportStatusVar.NoOfFail + 1
	}
	wg2.Done()
}

// func (dbconn *DBConn) ImportStatus() (NoOfImport int, NoOfFail int, TotalImport int) {
// report status
// return
// }

func (dbconn *DBConn) ImportAccountLogic(batch_id int) error {
	// Must connect to the database to get the information
	// first get the file name
	importbatch, err := ImportItf.GetImportBatchInfo(batch_id)
	if err != nil {
		log.Println(err)
		return err
	}
	// startTime := time.Now()
	cfg, _ := config.LoadConfig()
	log.Println("importAccountLogic begin processing file: ", importbatch.file_name_root)
	// define bufferred channel
	jobs := make(chan string, cfg.JOBS_QUEUE_SIZE)
	errs := make(chan string, cfg.ERRS_QUEUE_SIZE)
	// define a waitgroup to wait for all workers to finish his job
	var wg, wg2 sync.WaitGroup
	// start new routine for read the file. ParseFile()
	wg.Add(1)
	go parseFile(importbatch.file_name_root, jobs, &wg)
	// Create worker routines
	for i := 1; i <= cfg.JOBS_WORKER_SIZE; i++ {
		wg.Add(1)
		go workerDB(i, jobs, errs, &wg, importbatch)
	}
	// Start routine for writting errors
	wg2.Add(1)
	go writeErr(errs, &wg2)
	// Now wait all of them.
	wg.Wait()
	log.Println("All jobs done.")
	// When all finihed, then close the result channel
	close(errs)
	// Now wait for log errs tobe done
	wg2.Wait()
	log.Println("All log done")
	// Calculate the time of processing
	// endTime := time.Now()
	// diff := endTime.Sub(startTime)
	// log.Println("total time taken ", diff.Seconds(), "seconds")
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
			log.Println("Its EOF!")
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
