package config

import (
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/danielpaulus/go-ios/ios"
)

var HomeDir string

var UDID, BundleID, TestRunnerBundleID, XCTestConfig, WdaPort, WdaMjpegPort, AppiumPort string
var DeviceOSVersion, DeviceName, ScreenSize, StreamPort, DeviceOS, ContainerServerPort, DevicesHost, DeviceModel, StreamSize, RemoteControl, SupervisionPassword string
var AndroidScreenWidth, AndroidScreenHeight string

func SetHomeDir() {
	HomeDir, _ = os.UserHomeDir()
}

func GetEnv() {
	UDID = os.Getenv("DEVICE_UDID")
	BundleID = os.Getenv("WDA_BUNDLEID")
	TestRunnerBundleID = BundleID
	XCTestConfig = "WebDriverAgentRunner.xctest"
	WdaPort = os.Getenv("WDA_PORT")
	WdaMjpegPort = os.Getenv("MJPEG_PORT")
	AppiumPort = os.Getenv("APPIUM_PORT")
	DeviceOSVersion = os.Getenv("DEVICE_OS_VERSION")
	DeviceName = os.Getenv("DEVICE_NAME")
	ScreenSize = os.Getenv("SCREEN_SIZE")
	StreamPort = os.Getenv("STREAM_PORT")
	DeviceOS = os.Getenv("DEVICE_OS")
	ContainerServerPort = os.Getenv("CONTAINER_SERVER_PORT")
	DevicesHost = os.Getenv("DEVICES_HOST")
	DeviceModel = os.Getenv("DEVICE_MODEL")
	StreamSize = os.Getenv("STREAM_SIZE")
	RemoteControl = os.Getenv("REMOTE_CONTROL")
	SupervisionPassword = os.Getenv("SUPERVISION_PASSWORD")
	AndroidScreenWidth = os.Getenv("SCREEN_WIDTH")
	AndroidScreenHeight = os.Getenv("SCREEN_HEIGHT")
}

var Device ios.DeviceEntry

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
		panic(err)
	}
	return nil
}
