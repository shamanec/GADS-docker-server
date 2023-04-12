package config

import (
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/danielpaulus/go-ios/ios"
)

var HomeDir string

// Generic vars
var UDID, AppiumPort, DeviceOSVersion, DeviceName, ScreenSize, StreamPort, DeviceOS, ContainerServerPort, DevicesHost, DeviceModel string

// iOS vars
var BundleID, TestRunnerBundleID, XCTestConfig, WdaPort, WdaMjpegPort, SupervisionPassword string
var Device ios.DeviceEntry

// Android vars
var AndroidScreenWidth, AndroidScreenHeight, StreamSize string

func SetHomeDir() {
	HomeDir, _ = os.UserHomeDir()
}

func GetEnv() {
	// Generic vars
	UDID = os.Getenv("DEVICE_UDID")
	AppiumPort = os.Getenv("APPIUM_PORT")
	DeviceOSVersion = os.Getenv("DEVICE_OS_VERSION")
	DeviceName = os.Getenv("DEVICE_NAME")
	ScreenSize = os.Getenv("SCREEN_SIZE")
	StreamPort = os.Getenv("STREAM_PORT")
	DeviceOS = os.Getenv("DEVICE_OS")
	ContainerServerPort = os.Getenv("CONTAINER_SERVER_PORT")
	DevicesHost = os.Getenv("DEVICES_HOST")
	DeviceModel = os.Getenv("DEVICE_MODEL")

	// iOS vars
	BundleID = os.Getenv("WDA_BUNDLEID")
	TestRunnerBundleID = BundleID
	XCTestConfig = "WebDriverAgentRunner.xctest"
	WdaPort = os.Getenv("WDA_PORT")
	WdaMjpegPort = os.Getenv("MJPEG_PORT")
	SupervisionPassword = os.Getenv("SUPERVISION_PASSWORD")

	// Android vars
	StreamSize = os.Getenv("STREAM_SIZE")
	AndroidScreenWidth = os.Getenv("SCREEN_WIDTH")
	AndroidScreenHeight = os.Getenv("SCREEN_HEIGHT")
}

// Get ios.DeviceEntry for go-ios functions on container start
func GetDevice() error {
	err := retry.Do(
		func() error {
			availableDevice, err := ios.GetDevice(UDID)
			if err != nil {
				return err
			}
			Device = availableDevice
			return nil
		},
		retry.Attempts(3),
		retry.Delay(3*time.Second),
	)
	if err != nil {
		return err
	}

	return nil
}
