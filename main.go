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

func setLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	server_log_file, err := os.OpenFile(config.HomeDir+"/container-server.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0777)
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

	log.Fatal(http.ListenAndServe(":"+config.ContainerServerPort, myRouter))
}

func main() {
	config.SetHomeDir()
	config.GetEnv()

	if config.DeviceOS == "ios" {
		ios_server.ForwardWDA()
		ios_server.ForwardWDAStream()
	}

	setLogging()
	handleRequests()
}
