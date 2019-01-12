package main

import (
	"log"

	"github.com/loitd/vabackend/config"
	"github.com/loitd/vabackend/models"
	"github.com/loitd/vabackend/server"
)

func main() {
	// first load config file
	config, err := config.LoadConfig()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Config loaded. Database connection initializing ...")
	// config load done, now open connection to database and pass the connection to the interface
	models.InitDB(config)
	log.Println("Database connection initialized. Starting webserver ...")
	// Start server
	server.InitServer(config)
	// importitf.ImportAccount(121)
	// importitf.InsertAccount("900000011", "WOORIBANK", "1", "WRB000011", "3")

}
