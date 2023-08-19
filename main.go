package main

import (
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	android_server "github.com/shamanec/GADS-docker-server/android"
	"github.com/shamanec/GADS-docker-server/config"
	ios_server "github.com/shamanec/GADS-docker-server/ios"
	log "github.com/sirupsen/logrus"
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
	// // Create a new instance of the mux router
	// myRouter := mux.NewRouter().StrictSlash(true)
	// myRouter.HandleFunc("/installed-apps", GetInstalledApps)
	// myRouter.HandleFunc("/launch-app/{app}", LaunchApp)
	// myRouter.HandleFunc("/install-app/{app}", InstallApp)
	// myRouter.HandleFunc("/device-info", GetDeviceInfo)
	// myRouter.Handle("/stream", android_server.JpegStreamHandler())

	router := gin.Default()
	router.GET("/android-stream", android_server.StreamProxy)

	log.Fatal(http.ListenAndServe(":"+config.ContainerServerPort, router))
}

func main() {
	config.SetHomeDir()
	config.GetEnv()

	if config.DeviceOS == "ios" {
		go ios_server.SetupDevice()
	}

	if config.DeviceOS == "android" {
		go android_server.SetupDevice()
	}

	setLogging()
	handleRequests()
}
