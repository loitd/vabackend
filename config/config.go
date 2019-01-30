package config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
)

type Config struct {
	DATABASE_DRIVER       string
	DATABASE_URL          string
	LISTEN_PORT           string
	RESULT_FILEPATH       string
	WRITETIMEOUTINSECONDS time.Duration
	READTIMEOUTINSECONDS  time.Duration
	IDLETIMEOUTINSECONDS  time.Duration
	JOBS_QUEUE_SIZE       int
	ERRS_QUEUE_SIZE       int
	JOBS_WORKER_SIZE      int
}

func LoadConfig() (*Config, error) {
	// Open the json file
	f, err := os.Open("config.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// Parse json file into variables
	decoder := json.NewDecoder(f)
	conf := Config{}
	err = decoder.Decode(&conf)
	if err != nil {
		log.Println("LoadConfig:", err)
		return nil, err
	}
	return &conf, nil
}

func LogFile(filename string, msg string) error {
	// Log to files
	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Println("Error while LogFile.Open: ", err.Error())
		return err
	}
	defer f.Close()
	//
	startTime := time.Now()
	fmtmsg := fmt.Sprintf("[%d-%s-%d %d:%d:%d] %s\r\n", startTime.Year(), startTime.Month(), startTime.Day(), startTime.Hour(), startTime.Minute(), startTime.Second(), msg)
	_, err = f.WriteString(fmtmsg)
	if err != nil {
		log.Println("Error while LogFile.Write", err.Error())
		return err
	}
	f.Close()
	return nil
}
