package config

import "os"

var HomeDir string

var UDID, BundleID, TestRunnerBundleID, XCTestConfig, WdaPort, WdaMjpegPort, AppiumPort, DeviceOSVersion, DeviceName, ScreenSize, StreamPort, DeviceOS, ContainerServerPort, DevicesHost, DeviceModel, StreamSize, RemoteControl string

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
}
