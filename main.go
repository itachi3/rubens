package main

import (
	//Our libraries
	"controllers"
	"utils"

	//Third party libraries
	"github.com/gorilla/handlers"
	"github.com/julienschmidt/httprouter"

	//Standard Library
	"libs"
	"models"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var config utils.Config

func myLoggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile(config.Logs.AccessLog, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	utils.PanicError(err, "Error opening access log")
	return handlers.CombinedLoggingHandler(logFile, h)
}

func main() {
	// Read config and initialize the connections
	file, e := ioutil.ReadFile("config.json")
	utils.PanicError(e, "Config read error")
	e = json.Unmarshal(file, &config)
	utils.PanicError(e, "Config read error")

	// Establish connection pool
	c := libs.InitRedis(&config)
	conn := &models.Connections{
		RedisConn: c,
	}
	defer c.Close()

	// Instantiate a new router
	r := httprouter.New()

	// Get a Universal controller instance & initialize path hooks
	uc := controllers.NewImageController(r, conn, &config)
	uc.InitializeHooks()

	/* Initialize middleware for request log mechanism
	   Redirect all logs to error log
	   Fire up the server
	*/
	errorLog := utils.InitializeErrorLog(&config)
	defer errorLog.Close()
	log.Fatal(http.ListenAndServe(":8000", myLoggingHandler(r)))
}
