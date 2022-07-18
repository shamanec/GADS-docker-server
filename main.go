package main

import (
	"net/http"
	"os"

	"github.com/shamanec/GADS-docker-server/config"
	ios_server "github.com/shamanec/GADS-docker-server/ios"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var server_log_file *os.File

//var udid = os.Getenv("DEVICE_UDID")
var udid = "ccec159ba0219c9fa0d0fc3d85451ab0dcfebd16"
var DeviceOS = "ios"

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
	myRouter.HandleFunc("/start-wda", ios_server.StartWDA)
	myRouter.HandleFunc("/device-info", GetDeviceInfo)

	log.Fatal(http.ListenAndServe(":10001", myRouter))
}

func main() {
	config.SetHomeDir()
	ios_server.ForwardWDA()
	ios_server.ForwardWDAStream()
	setLogging()
	handleRequests()
}
