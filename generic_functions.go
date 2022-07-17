package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	android_server "github.com/shamanec/GADS-docker-server/android"
	"github.com/shamanec/GADS-docker-server/helpers"
	ios_server "github.com/shamanec/GADS-docker-server/ios"
)

type DeviceInfo struct {
	InstalledApps []string `json:"installed_apps"`
}

func GetInstalledApps(w http.ResponseWriter, r *http.Request) {
	var appIDs DeviceInfo

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
