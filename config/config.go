package config

import "os"

var HomeDir string

// var UDID = os.Getenv("DEVICE_UDID")
// var BundleID = os.Getenv("WDA_BUNDLEID")
// var TestRunnerBundleID = BundleID
// var XCTestConfig = "WebDriverAgentRunner.xctest"
// var WDA_PORT = os.Getenv("WDA_PORT")
// var WDA_MJPEG_PORT = os.Getenv("MJPEG_PORT")
// var APPIUM_PORT = "4723"
// var DEVICE_OS_VERSION = os.Getenv("DEVICE_OS_VERSION")
// var DEVICE_NAME = os.Getenv("DEVICE_NAME")

// var UDID = os.Getenv("DEVICE_UDID")
// var BundleID = os.Getenv("WDA_BUNDLEID")
// var TestRunnerBundleID = BundleID
// var XCTestConfig = "WebDriverAgentRunner.xctest"
// var WDA_PORT = os.Getenv("WDA_PORT")
// var WDA_MJPEG_PORT = os.Getenv("MJPEG_PORT")
// var APPIUM_PORT = "4723"
// var DEVICE_OS_VERSION = os.Getenv("DEVICE_OS_VERSION")
// var DEVICE_NAME = os.Getenv("DEVICE_NAME")

// var UDID = "ccec159ba0219c9fa0d0fc3d85451ab0dcfebd16"
// var BundleID = "com.shamanec.WebDriverAgentRunner.xctrunner"
// var TestRunnerBundleID = BundleID
// var XCTestConfig = "WebDriverAgentRunner.xctest"
// var WdaPort = "20004"
// var WdaMjpegPort = "20104"
// var AppiumPort = "4723"
// var DeviceOSVersion = "15.4"
// var DeviceName = "Device_name"
// var ScreenSize = "375x667"
// var StreamPort = "1000"
// var DeviceOS = "ios"

var UDID, BundleID, TestRunnerBundleID, XCTestConfig, WdaPort, WdaMjpegPort, AppiumPort, DeviceOSVersion, DeviceName, ScreenSize, StreamPort, DeviceOS, ContainerServerPort, DevicesHost, DeviceModel, StreamSize string

func SetHomeDir() {
	HomeDir, _ = os.UserHomeDir()
}

func GetEnv() {
	// os.Setenv("DEVICE_UDID", "ccec159ba0219c9fa0d0fc3d85451ab0dcfebd16")
	// os.Setenv("WDA_BUNDLEID", "com.shamanec.WebDriverAgentRunner.xctrunner")
	// os.Setenv("WDA_PORT", "20004")
	// os.Setenv("MJPEG_PORT", "20104")
	// os.Setenv("DEVICE_OS_VERSION", "15.4")
	// os.Setenv("DEVICE_NAME", "Device_name")
	// os.Setenv("SCREEN_SIZE", "375x667")
	// os.Setenv("STREAM_PORT", "1000")
	// os.Setenv("DEVICE_OS", "ios")
	// os.Setenv("APPIUM_PORT", "4723")

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
}
