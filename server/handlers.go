package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/loitd/vabackend/models"
)

func ImportStatusHandlerv12(w http.ResponseWriter, r *http.Request) {
	// Check the import status
	output, err := json.Marshal(models.ImportStatusVar)
	if err != nil {
		log.Println(err)
	}
	fmt.Fprintf(w, string(output))
}

func getStatus() ([]byte, error) {
	// Check the import status
	output, err := json.Marshal(models.ImportStatusVar)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	// fmt.Fprintf(w, string(output))
	return output, nil
}

func ImportAccountHandlerv12(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	// Reset counters for every new request
	models.ResetResult()
	// check for request method is POST or GET
	var bid string
	switch r.Method {
	case http.MethodGet:
		// handle GET
		log.Println("Handling GET request ...")
		query := r.URL.Query()
		bid = query.Get("batch_id")
	case http.MethodPost:
		// handle POST
		log.Println("Handling POST request ...")
		bid = r.FormValue("batch_id")
	default:
		// print error
		log.Println("Method not allowed")
	}
	// check len bid
	if len(bid) < 1 {
		log.Println("Invalid batch_id")
		return
	}
	batchid, err := strconv.Atoi(bid)
	if err != nil {
		log.Println("Invalid type of batch_id")
		return
	}
	log.Println("WEGOT: ", bid)
	// w.Write([]byte("hello con de"))
	// models.ImportItf.ImportAccountLogic("fbk_vir_001_20181206_001.dat")
	models.ImportItf.ImportAccountLogic(batchid)
	// return for caller
	// get result and return
	output, err := getStatus()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write(output)
	}
	// w.Write([]byte("{result: called successfully}"))
	// Calculate the time of processing
	endTime := time.Now()
	diff := endTime.Sub(startTime)
	log.Println("total time taken ", diff.Seconds(), "seconds")
}

func ImportAccountHandlerv10(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	batch_id := query.Get("batch_id")
	if len(batch_id) < 1 {
		log.Println("Invalid batch_id")
		return
	}
	log.Println("WEGOT: ", batch_id)
	//
	batchid, err := strconv.Atoi(batch_id)
	if err != nil {
		log.Println(err)
		return
	}
	// import now
	err = models.ImportItf.ImportAccount(batchid)
	if err != nil {
		log.Println(err)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write([]byte("hello con de"))
	//
}
