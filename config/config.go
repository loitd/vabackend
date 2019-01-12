package config

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type Config struct {
	DATABASE_DRIVER       string
	DATABASE_URL          string
	LISTEN_PORT           string
	WRITETIMEOUTINSECONDS time.Duration
	READTIMEOUTINSECONDS  time.Duration
	IDLETIMEOUTINSECONDS  time.Duration
}

func LoadConfig() (*Config, error) {
	// Open the json file
	f, _ := os.Open("config.json")
	defer f.Close()
	// Parse json file into variables
	decoder := json.NewDecoder(f)
	conf := Config{}
	err := decoder.Decode(&conf)
	if err != nil {
		log.Println("LoadConfig:", err)
		return nil, err
	}
	return &conf, nil
}
