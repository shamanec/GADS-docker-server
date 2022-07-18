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

var UDID = "ccec159ba0219c9fa0d0fc3d85451ab0dcfebd16"
var BundleID = "com.shamanec.WebDriverAgentRunner.xctrunner"
var TestRunnerBundleID = BundleID
var XCTestConfig = "WebDriverAgentRunner.xctest"
var WdaPort = "20004"
var WdaMjpegPort = "20104"
var AppiumPort = "4723"
var DeviceOSVersion = "15.4"
var DeviceName = "Device_name"
var ScreenSize = "375x667"
var StreamPort = "1000"

func SetHomeDir() {
	HomeDir, _ = os.UserHomeDir()
}
