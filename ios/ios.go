package ios_server

import (
	"errors"
	"strconv"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/forward"
	"github.com/danielpaulus/go-ios/ios/installationproxy"
	"github.com/danielpaulus/go-ios/ios/instruments"
	"github.com/danielpaulus/go-ios/ios/zipconduit"
	"github.com/shamanec/GADS-docker-server/config"
	"github.com/shamanec/GADS-docker-server/helpers"
	log "github.com/sirupsen/logrus"
)

type IOSDevice struct {
	InstalledApps []string        `json:"installed_apps"`
	DeviceConfig  IOSDeviceConfig `json:"device_config"`
}

type IOSDeviceConfig struct {
	AppiumPort          string `json:"appium_port"`
	DeviceName          string `json:"device_name"`
	DeviceOSVersion     string `json:"device_os_version"`
	DeviceUDID          string `json:"device_udid"`
	WdaMjpegPort        string `json:"wda_mjpeg_port"`
	WdaPort             string `json:"wda_port"`
	DeviceScreenSize    string `json:"screen_size"`
	DeviceHost          string `json:"device_host"`
	DeviceModel         string `json:"device_model"`
	ContainerServerPort string `json:"container_server_port"`
	DeviceOS            string `json:"device_os"`
}

func GetDeviceInfo() (string, error) {
	bundleIDs, err := GetInstalledApps()
	if err != nil {
		return "", err
	}

	config := IOSDeviceConfig{
		AppiumPort:          config.AppiumPort,
		DeviceName:          config.DeviceName,
		DeviceUDID:          config.UDID,
		DeviceOSVersion:     config.DeviceOSVersion,
		WdaMjpegPort:        config.WdaMjpegPort,
		WdaPort:             config.WdaPort,
		DeviceScreenSize:    config.ScreenSize,
		DeviceHost:          config.DevicesHost,
		DeviceModel:         config.DeviceModel,
		ContainerServerPort: config.ContainerServerPort,
		DeviceOS:            config.DeviceOS,
	}

	deviceInfo := IOSDevice{
		InstalledApps: bundleIDs,
		DeviceConfig:  config,
	}

	return helpers.ConvertToJSONString(deviceInfo), nil
}

func InstallApp(fileName string) error {
	filePath := "/opt/" + fileName

	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "install_app",
		}).Error("Could not get device when installing app. Error: " + err.Error())
		return errors.New("Failed installing application:" + err.Error())
	}

	conn, err := zipconduit.New(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "install_app",
		}).Error("Could not create zipconduit when installing app. Error: " + err.Error())
		return errors.New("Failed installing application:" + err.Error())
	}

	err = conn.SendFile(filePath)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "install_app",
		}).Error("Could not install app. Error: " + err.Error())
		return errors.New("Failed installing application:" + err.Error())
	}
	return nil
}

func UninstallApp(bundle_id string) error {
	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "uninstall_ios_app",
		}).Error("Could not get device when uninstalling app with bundleID:'" + bundle_id + "'. Error: " + err.Error())
		return errors.New("Error")
	}

	svc, err := installationproxy.New(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "uninstall_ios_app",
		}).Error("Failed connecting installationproxy when uninstalling app with bundleID:'" + bundle_id + "'. Error: " + err.Error())
		return errors.New("Error")
	}

	err = svc.Uninstall(bundle_id)

	if err != nil {
		log.WithFields(log.Fields{
			"event": "uninstall_ios_app",
		}).Error("Failed uninstalling app with bundleID:'" + bundle_id + "'. Error: " + err.Error())
		return errors.New("Error")
	}
	return nil
}

type goIOSAppList []struct {
	BundleID string `json:"CFBundleIdentifier"`
}

func GetInstalledApps() ([]string, error) {
	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not get device with UDID: '" + config.UDID + "'. Error: " + err.Error())
		return nil, errors.New("Could not get device with UDID: '" + config.UDID + "'. Error: " + err.Error())
	}

	svc, err := installationproxy.New(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not create installation proxy for device with UDID: '" + config.UDID + "'. Error: " + err.Error())
		return nil, errors.New("Could not create installation proxy for device with UDID: '" + config.UDID + "'. Error: " + err.Error())
	}

	user_apps, err := svc.BrowseUserApps()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not get user apps for device with UDID: '" + config.UDID + "'. Error: " + err.Error())
		return nil, errors.New("Could not get user apps for device with UDID: '" + config.UDID + "'. Error: " + err.Error())
	}

	var data goIOSAppList

	err = helpers.UnmarshalJSONString(helpers.ConvertToJSONString(user_apps), &data)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "device_container_create",
		}).Error("Could not unmarshal request body when uninstalling iOS app")
		return nil, errors.New("Could not unmarshal user apps json")
	}

	var bundleIDs []string

	for _, dataObject := range data {
		bundleIDs = append(bundleIDs, dataObject.BundleID)
	}

	return bundleIDs, nil
}

func LaunchApp(bundleID string) (uint64, error) {

	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not get device with UDID: '" + config.UDID + "'. Error: " + err.Error())
		return 0, errors.New("Could not get device with UDID: '" + config.UDID + "'. Error: " + err.Error())
	}

	pControl, err := instruments.NewProcessControl(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not create process control for device with UDID: " + config.UDID + ". Error: " + err.Error())
		return 0, errors.New("Could not create process control for device with UDID: '" + config.UDID + "'. Error: " + err.Error())
	}

	pid, err := pControl.LaunchApp(bundleID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not launch app for device with UDID: " + config.UDID + ". Error: " + err.Error())
		return 0, errors.New("Could not launch app for device with UDID: '" + config.UDID + "'. Error: " + err.Error())
	}

	return pid, nil
}

func ForwardWDA() error {
	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not get device with UDID: '" + config.UDID + "'. Error: " + err.Error())
		return errors.New("Error")
	}

	wda_port, err := strconv.ParseUint(config.WdaPort, 10, 32)
	if err != nil {
		return err
	}

	forward.Forward(device, uint16(wda_port), uint16(wda_port))

	return nil
}

func ForwardWDAStream() error {
	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not get device with UDID: '" + config.UDID + "'. Error: " + err.Error())
		return errors.New("Error")
	}

	wda_mjpeg_port, err := strconv.ParseUint(config.WdaMjpegPort, 10, 32)
	if err != nil {
		return err
	}

	forward.Forward(device, uint16(wda_mjpeg_port), uint16(wda_mjpeg_port))

	return nil
}
