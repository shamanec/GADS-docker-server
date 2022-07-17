package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var server_log_file *os.File

//var udid = os.Getenv("DEVICE_UDID")
var udid = "00008030000418C136FB802E"
var DeviceOS = "android"

func setLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	server_log_file, err := os.OpenFile("/home/shamanec/logs/project.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
	if err != nil {
		panic(err)
	}
	log.SetOutput(server_log_file)
}

func handleRequests() {
	// Create a new instance of the mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/installed-apps", GetInstalledApps)
	myRouter.HandleFunc("/launch-app/{app}", LaunchApp)
	myRouter.HandleFunc("/install-app/{app}", InstallApp)

	log.Fatal(http.ListenAndServe(":10001", myRouter))
}

func main() {
	setLogging()
	handleRequests()
}
