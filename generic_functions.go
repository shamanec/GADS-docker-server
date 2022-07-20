package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	android_server "github.com/shamanec/GADS-docker-server/android"
	"github.com/shamanec/GADS-docker-server/config"
	"github.com/shamanec/GADS-docker-server/helpers"
	ios_server "github.com/shamanec/GADS-docker-server/ios"
)

type AndroidDevice struct {
	InstalledApps InstalledApps
	DeviceConfig  AndroidDeviceConfig `json:"device_config"`
}

type AndroidDeviceConfig struct {
	AppiumPort      string `json:"appium_port"`
	DeviceName      string `json:"device_name"`
	DeviceOSVersion string `json:"device_os_version"`
	DeviceUDID      string `json:"device_udid"`
	StreamSize      string `json:"stream_size"`
	StreamPort      string `json:"stream_port"`
}

type InstalledApps struct {
	InstalledApps []string `json:"installed_apps"`
}

func GetDeviceInfo(w http.ResponseWriter, r *http.Request) {
	var info string
	var err error

	if config.DeviceOS == "ios" {
		info, err = ios_server.GetDeviceInfo()
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}

	} else {
		info = "test"
	}

	fmt.Fprintf(w, info)
}

func GetInstalledApps(w http.ResponseWriter, r *http.Request) {
	var appIDs InstalledApps

	if config.DeviceOS == "ios" {
		bundleIDs, err := ios_server.GetInstalledApps()
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}
		appIDs.InstalledApps = bundleIDs
		fmt.Fprintf(w, helpers.ConvertToJSONString(appIDs))

	} else {
		packageNames, err := android_server.GetInstalledApps()
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}

		appIDs.InstalledApps = packageNames
		fmt.Fprintf(w, helpers.ConvertToJSONString(appIDs))
	}
}

func LaunchApp(w http.ResponseWriter, r *http.Request) {
	// Get the request path vars
	vars := mux.Vars(r)
	app := vars["app"]

	if config.DeviceOS == "ios" {
		_, err := ios_server.LaunchApp(app)
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}
	} else {
		err := android_server.LaunchApp(app)
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}
	}
	fmt.Fprintf(w, "App '"+app+"' is started.")
}

func InstallApp(w http.ResponseWriter, r *http.Request) {
	// Get the request path vars
	vars := mux.Vars(r)
	appName := vars["app"]

	if config.DeviceOS == "ios" {
		err := ios_server.InstallApp(appName)
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "App '"+appName+"' installed.")
	} else {
		err := android_server.InstallApp(appName)
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}
		fmt.Fprintf(w, "App '"+appName+"' installed.")
	}
}
