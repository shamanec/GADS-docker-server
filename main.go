package main

import (
	"net/http"
	"os"

	ios_server "github.com/shamanec/GADS-docker-server/ios"
	log "github.com/sirupsen/logrus"

	"github.com/gorilla/mux"
)

var server_log_file *os.File

//var udid = os.Getenv("DEVICE_UDID")
var udid = "00008030000418C136FB802E"

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
	ios_server.InstallWDA()
	ios_server.MountDiskImages()
	ios_server.StartWDA()
	ios_server.StartAppiumIOS()

	log.Fatal(http.ListenAndServe(":10001", myRouter))
}

func setupIOSDevice() {

}

func setupAndroidDevice() {

}

func main() {
	setLogging()
	handleRequests()
	if os.Getenv("WDA_BUNDLEID") != "" {
		setupIOSDevice()
	} else {
		setupAndroidDevice()
	}
}
