package server

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	mux "github.com/gorilla/mux"
	"github.com/loitd/vabackend/config"
)

func InitServer(cfg *config.Config) {
	// Start webserver at configed port using mux in v1.2
	router := mux.NewRouter()
	// Routes consists of a path and a handler function
	router.HandleFunc("/api/v1.0/importaccount", ImportAccountHandlerv10)
	router.HandleFunc("/api/v1.2/importaccount", ImportAccountHandlerv12)
	// Bind port and pass router in. New server with timeout and graceful shutdown from v1.2 (go > 1.8 only)
	// define wait time flag and parse this flag after definition
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-webserver-timeout", time.Second*15, "the duration for graceful shutdown/close")
	flag.Parse()

	// define a server
	server := &http.Server{
		Addr:         cfg.LISTEN_PORT,
		WriteTimeout: time.Second * cfg.WRITETIMEOUTINSECONDS,
		ReadTimeout:  time.Second * cfg.READTIMEOUTINSECONDS,
		IdleTimeout:  time.Second * cfg.IDLETIMEOUTINSECONDS,
		Handler:      router,
	}
	// run server in a go routine so that's not blocked
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
			server.Close()
		}
	}()
	log.Println("Webserver started. Listening on: ", cfg.LISTEN_PORT)
	// define a channel for receiving signal with size=1
	ch := make(chan os.Signal, 1)
	// catch kill or interrupt signals
	signal.Notify(ch, os.Kill, os.Interrupt)
	// start new routine listen for signal
	// this function will wait for shutdown signal and react
	// block untilreceive signal
	<-ch
	// create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	log.Println("Please wait while shutting down.")
	server.Shutdown(ctx)
	os.Exit(0)
}
