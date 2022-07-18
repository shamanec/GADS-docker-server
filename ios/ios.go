package ios_server

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strconv"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/forward"
	"github.com/danielpaulus/go-ios/ios/imagemounter"
	"github.com/danielpaulus/go-ios/ios/installationproxy"
	"github.com/danielpaulus/go-ios/ios/instruments"
	"github.com/danielpaulus/go-ios/ios/testmanagerd"
	"github.com/danielpaulus/go-ios/ios/zipconduit"
	"github.com/shamanec/GADS-docker-server/config"
	"github.com/shamanec/GADS-docker-server/helpers"
	log "github.com/sirupsen/logrus"
)

// var udid = os.Getenv("DEVICE_UDID")
// var bundleid = os.Getenv("WDA_BUNDLEID")
// var testrunnerbundleid = bundleid
// var xctestconfig = "WebDriverAgentRunner.xctest"
// var wda_port = os.Getenv("WDA_PORT")
// var wda_mjpeg_port = os.Getenv("MJPEG_PORT")
// var appium_port = "4723"
// var device_os_version = os.Getenv("DEVICE_OS_VERSION")
// var device_name = os.Getenv("DEVICE_NAME")

var udid = "ccec159ba0219c9fa0d0fc3d85451ab0dcfebd16"
var bundleid = "com.shamanec.WebDriverAgentRunner.xctrunner"
var testrunnerbundleid = bundleid
var xctestconfig = "WebDriverAgentRunner.xctest"
var wda_port = "20004"
var wda_mjpeg_port = "20104"
var appium_port = "4723"
var device_os_version = "15.4"
var device_name = "Device_name"

func GetDeviceInfo() {

}

func StartAppiumIOS() {

	capabilities := `{"mjpegServerPort": ` + wda_mjpeg_port +
		`, "clearSystemFiles": "false",` +
		`"webDriverAgentUrl":"http://192.168.1.6:` + wda_port + `",` +
		`"preventWDAAttachments": "true",` +
		`"simpleIsVisibleCheck": "false",` +
		`"wdaLocalPort": "` + wda_port + `",` +
		`"platformVersion": "` + device_os_version + `",` +
		`"automationName":"XCUITest",` +
		`"platformName": "iOS",` +
		`"deviceName": "` + device_name + `",` +
		`"wdaLaunchTimeout": "120000",` +
		`"wdaConnectionTimeout": "240000"}`

	commandString := "appium -p " + appium_port + " --udid=" + udid + " --log-timestamp --default-capabilities '" + capabilities + "'"
	cmd := exec.Command("bash", "-c", commandString)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	fmt.Println("command is: " + commandString)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "start_appium_ios",
		}).Error("test: " + err.Error())
		return
	}
	log.WithFields(log.Fields{
		"event": "start_appium_ios",
	}).Info("test")
}

// func StartWDA(w http.ResponseWriter, r *http.Request) {

// 	device, err := ios.GetDevice(udid)
// 	if err != nil {
// 		log.WithFields(log.Fields{
// 			"event": "run_wda",
// 		}).Error("Could not get device when installing app. Error: " + err.Error())
// 	}

// 	go func() {
// 		err := testmanagerd.RunXCUIWithBundleIds(bundleid,
// 			testrunnerbundleid,
// 			xctestconfig,
// 			device,
// 			[]string{},
// 			[]string{"USE_PORT=" + wda_port, "MJPEG_SERVER_PORT=" + wda_mjpeg_port})

// 		log.WithFields(log.Fields{
// 			"event": "run_wda",
// 		}).Error("Failed running wda. Error: " + err.Error())
// 		fmt.Println(err.Error())
// 	}()

// 	fmt.Fprintf(w, "Started WDA on port: "+wda_port)
// }

func StopWDA() {
	err := testmanagerd.CloseXCUITestRunner()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "stop_wda",
		}).Error("Failed closing wda runner. Error: " + err.Error())
	}
}

func InstallWDA() error {
	err := InstallApp("WebDriverAgent.ipa")
	return err
}

func InstallApp(fileName string) error {
	filePath := "/opt/" + fileName

	device, err := ios.GetDevice(udid)
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

func MountDiskImages() error {
	device, err := ios.GetDevice(udid)

	if err != nil {
		log.WithFields(log.Fields{
			"event": "mount_dev_images",
		}).Error("Could not get device when mounting dev images. Error: " + err.Error())
		return errors.New("Failed mounting disk images")
	}

	mountConn, err := imagemounter.New(device)
	signatures, err := mountConn.ListImages()

	if len(signatures) == 0 {
		basedir := "/opt/devimages"

		err = imagemounter.FixDevImage(device, basedir)
		log.WithFields(log.Fields{
			"event": "mount_dev_images",
		}).Error("Could not get device when mounting dev images. Error: " + err.Error())
		return errors.New("Failed mounting disk images")
	} else {
		log.WithFields(log.Fields{
			"event": "mount_dev_images",
		}).Info("DevImages are mounted on device with UDID: '" + udid)
		return nil
	}
}

func UninstallApp(bundle_id string) error {
	device, err := ios.GetDevice(udid)
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
	device, err := ios.GetDevice(udid)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
		return nil, errors.New("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
	}

	svc, err := installationproxy.New(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not create installation proxy for device with UDID: '" + udid + "'. Error: " + err.Error())
		return nil, errors.New("Could not create installation proxy for device with UDID: '" + udid + "'. Error: " + err.Error())
	}

	user_apps, err := svc.BrowseUserApps()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_device_apps",
		}).Error("Could not get user apps for device with UDID: '" + udid + "'. Error: " + err.Error())
		return nil, errors.New("Could not get user apps for device with UDID: '" + udid + "'. Error: " + err.Error())
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

	device, err := ios.GetDevice(udid)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
		return 0, errors.New("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
	}

	pControl, err := instruments.NewProcessControl(device)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not create process control for device with UDID: " + udid + ". Error: " + err.Error())
		return 0, errors.New("Could not create process control for device with UDID: '" + udid + "'. Error: " + err.Error())
	}

	pid, err := pControl.LaunchApp(bundleID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not launch app for device with UDID: " + udid + ". Error: " + err.Error())
		return 0, errors.New("Could not launch app for device with UDID: '" + udid + "'. Error: " + err.Error())
	}

	return pid, nil
}

func ForwardWDA() error {
	device, err := ios.GetDevice(udid)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
		return errors.New("Error")
	}

	wda_port, err := strconv.ParseUint(config.WdaPort, 10, 32)
	if err != nil {
		return err
	}

	// wda_mjpeg_port, err := strconv.ParseUint(config.WdaMjpegPort, 10, 32)
	// if err != nil {
	// 	return err
	// }

	forward.Forward(device, uint16(wda_port), uint16(wda_port))

	return nil
}

func ForwardWDAStream() error {
	device, err := ios.GetDevice(udid)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "ios_launch_app",
		}).Error("Could not get device with UDID: '" + udid + "'. Error: " + err.Error())
		return errors.New("Error")
	}

	wda_mjpeg_port, err := strconv.ParseUint(config.WdaMjpegPort, 10, 32)
	if err != nil {
		return err
	}

	forward.Forward(device, uint16(wda_mjpeg_port), uint16(wda_mjpeg_port))

	return nil
}
