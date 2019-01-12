package server

import (
	"log"
	"net/http"
	"strconv"

	"github.com/loitd/vabackend/models"
)

func ImportAccountHandlerv12(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	batch_id := query.Get("batch_id")
	if len(batch_id) < 1 {
		log.Println("Invalid batch_id")
		return
	}
	log.Println("WEGOT: ", batch_id)
	w.Write([]byte("hello con de"))
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
