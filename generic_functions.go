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

type IOSDevice struct {
	InstalledApps []string        `json:"installed_apps"`
	DeviceConfig  IOSDeviceConfig `json:"device_config"`
}

type AndroidDevice struct {
	InstalledApps InstalledApps
	DeviceConfig  AndroidDeviceConfig `json:"device_config"`
}

type IOSDeviceConfig struct {
	AppiumPort       string `json:"appium_port"`
	DeviceName       string `json:"device_name"`
	DeviceOSVersion  string `json:"device_os_version"`
	DeviceUDID       string `json:"device_udid"`
	WdaMjpegPort     string `json:"wda_mjpeg_port"`
	WdaPort          string `json:"wda_port"`
	WdaURL           string `json:"wda_url"`
	WdaMjpegURL      string `json:"wda_stream_url"`
	DeviceScreenSize string `json:"screen_size"`
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
	if DeviceOS == "ios" {
		bundleIDs, err := ios_server.GetInstalledApps()
		if err != nil {
			helpers.JSONError(w, "", err.Error(), 500)
			return
		}

		config := IOSDeviceConfig{
			AppiumPort:       config.AppiumPort,
			DeviceName:       config.DeviceName,
			DeviceUDID:       config.UDID,
			DeviceOSVersion:  config.DeviceOSVersion,
			WdaMjpegPort:     config.WdaMjpegPort,
			WdaPort:          config.WdaPort,
			WdaURL:           "http://192.168.1.6:20004",
			WdaMjpegURL:      "http://192.168.1.6:20104",
			DeviceScreenSize: config.ScreenSize,
		}

		deviceInfo := IOSDevice{
			InstalledApps: bundleIDs,
			DeviceConfig:  config,
		}

		fmt.Fprintf(w, helpers.ConvertToJSONString(deviceInfo))
	}
}

func GetInstalledApps(w http.ResponseWriter, r *http.Request) {
	var appIDs InstalledApps

	if DeviceOS == "ios" {
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

	if DeviceOS == "ios" {
		_, err := ios_server.LaunchApp(app)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	} else {
		err := android_server.LaunchApp(app)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
	}
	fmt.Fprintf(w, "App '"+app+"' is started.")
}

func InstallApp(w http.ResponseWriter, r *http.Request) {
	// Get the request path vars
	vars := mux.Vars(r)
	appName := vars["app"]

	if DeviceOS == "ios" {
		err := ios_server.InstallApp(appName)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, "App '"+appName+"' installed.")
	} else {
		err := android_server.InstallApp(appName)
		if err != nil {
			fmt.Fprintf(w, err.Error())
			return
		}
		fmt.Fprintf(w, "App '"+appName+"' installed.")
	}
}
