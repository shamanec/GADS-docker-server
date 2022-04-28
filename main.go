package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var server_log_file *os.File

func setLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	server_log_file, err := os.OpenFile("./opt/logs/project.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(server_log_file)
}

func handleRequests() {
	// Create a new instance of the mux router
	myRouter := mux.NewRouter().StrictSlash(true)

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func setupIOSDevice() {

}

func setupAndroidDevice() {

}

func main() {
	handleRequests()
	if os.Getenv("WDA_BUNDLEID") != "" {
		setupIOSDevice()
	} else {
		setupAndroidDevice()
	}
}
