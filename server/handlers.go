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

func ImportAccountHandlerv12(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	query := r.URL.Query()
	batch_id := query.Get("batch_id")
	if len(batch_id) < 1 {
		log.Println("Invalid batch_id")
		return
	}
	log.Println("WEGOT: ", batch_id)
	w.Write([]byte("hello con de"))
	// models.ImportItf.ImportAccountLogic("fbk_vir_001_20181206_001.dat")
	models.ImportItf.ImportAccountLogic("sample.dat")
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
