package main

import (
	"fmt"
	"net/http"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/installationproxy"
	"github.com/shamanec/GADS-docker-server/helpers"

	log "github.com/sirupsen/logrus"
)

type InstalledApps struct {
	InstalledApps []string `json:"installed_apps"`
}

type goIOSAppList []struct {
	BundleID string `json:"CFBundleIdentifier"`
}

func GetInstalledApps(w http.ResponseWriter, r *http.Request) {
	device, err := ios.GetDevice(udid)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
		helpers.JSONError(w, "", "Could not get device with UDID: '"+udid+"'. Error: "+err.Error(), 500)
		return
	}

	svc, err := installationproxy.New(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not create installation proxy for device with UDID: '" + udid + "'. Error: " + err.Error())
		helpers.JSONError(w, "", "Could not create installation proxy for device with UDID: '"+udid+"'. Error: "+err.Error(), 500)
		return
	}

	user_apps, err := svc.BrowseUserApps()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not get user apps for device with UDID: '" + udid + "'. Error: " + err.Error())
		helpers.JSONError(w, "", "Could not get user apps for device with UDID: '"+udid+"'. Error: "+err.Error(), 500)
		return
	}

	var data goIOSAppList

	err = helpers.UnmarshalJSONString(helpers.ConvertToJSONString(user_apps), &data)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "device_container_create",
		}).Error("Could not unmarshal request body when uninstalling iOS app")
		helpers.JSONError(w, "", "Could not unmarshal user apps json", 500)
		return
	}

	var bundleIDs InstalledApps

	for _, dataObject := range data {
		bundleIDs.InstalledApps = append(bundleIDs.InstalledApps, dataObject.BundleID)
	}

	fmt.Fprintf(w, helpers.ConvertToJSONString(bundleIDs))
}
