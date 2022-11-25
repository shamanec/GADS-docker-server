package ios_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/forward"
	"github.com/danielpaulus/go-ios/ios/imagemounter"
	"github.com/danielpaulus/go-ios/ios/testmanagerd"
	"github.com/shamanec/GADS-docker-server/config"
	log "github.com/sirupsen/logrus"
)

// Start usbmuxd service after starting the container
func startUsbmuxd() {
	prg := "usbmuxd"
	arg1 := "-U"
	arg2 := "usbmux"
	arg3 := "-f"

	// Build the usbmuxd command
	cmd := exec.Command(prg, arg1, arg2, arg3)

	// Run the command to start usbmuxd
	err := cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

// Mount the developer disk images downloading them automatically
func mountDeveloperImage(device ios.DeviceEntry) error {
	err := imagemounter.FixDevImage(device, "/opt/DeveloperDiskImages")
	if err != nil {
		return err
	}

	return nil
}

// Pair the device which is expected to be supervised
func pairDevice(device ios.DeviceEntry) error {
	p12, err := os.ReadFile("/opt/supervision.p12")
	if err != nil {
		return err
	}

	err = ios.PairSupervised(device, p12, config.SupervisionPassword)
	if err != nil {
		return err
	}

	return nil
}

// Start the Appium server for the device
func startAppium() {
	prg := "appium"
	arg1 := "-p 4723"
	arg2 := "--log-timestamp"
	arg3 := "--allow-cors"
	arg4 := "--session-override"
	arg5 := `--default-capabilities '{"appium:udid": "'` + config.UDID + `'", "appium:mjpegServerPort": "9100", "appium:clearSystemFiles": "false", "appium:webDriverAgentUrl":"http://localhost:8100", "appium:preventWDAAttachments": "true", "appium:simpleIsVisibleCheck": "false", "appium:wdaLocalPort": "8100", "appium:platformVersion": "'` + config.DeviceOSVersion + `'", "appium:automationName":"XCUITest", "platformName": "iOS", "appium:deviceName": "'` + config.DeviceName + `'", "appium:wdaLaunchTimeout": "120000", "appium:wdaConnectionTimeout": "240000"}'`
	cmd := exec.Command(prg, arg1, arg2, arg3, arg4, arg5)

	outfile, err := os.Create("/opt/logs/appium.log")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

// Install WebDriverAgent on the device
func prepareWDA(device ios.DeviceEntry) error {
	// fmt.Println("Installing WebDriverAgent")
	// prg := "ios"
	// arg1 := "install"
	// arg2 := "--path=/opt/WebDriverAgent.ipa"
	// arg3 := "--udid=" + config.UDID

	// cmd := exec.Command(prg, arg1, arg2, arg3)

	// err := cmd.Run()
	// if err != nil {
	// 	return err
	// }

	err := InstallAppWithDevice(device, "WebDriverAgent.ipa")
	if err != nil {
		return err
	}

	go startWDA()
	return nil
}

// Start the WebDriverAgent on the device
func startWDA() {
	fmt.Println("Starting WDA")
	prg := "ios"
	arg1 := "runwda"
	arg2 := "--bundleid=" + config.BundleID
	arg3 := "--testrunnerbundleid=" + config.BundleID
	arg4 := "--xctestconfig=WebDriverAgentRunner.xctest"
	arg5 := "--udid=" + config.UDID

	cmd := exec.Command(prg, arg1, arg2, arg3, arg4, arg5)

	outfile, err := os.Create("/opt/logs/wda.log")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	cmd.Stdout = outfile
	cmd.Stderr = outfile

	err = cmd.Run()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
}

// Start WebDriverAgent directly using go-ios modules
func StartWDAInternal() error {
	device, err := ios.GetDevice(config.UDID)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "run_wda",
		}).Error("Could not get device when installing app. Error: " + err.Error())
		return err
	}

	go func() {
		err := testmanagerd.RunXCUIWithBundleIdsCtx(nil, config.BundleID,
			config.TestRunnerBundleID,
			config.XCTestConfig,
			device,
			[]string{},
			[]string{"USE_PORT=" + config.WdaPort, "MJPEG_SERVER_PORT=" + config.WdaMjpegPort})

		log.WithFields(log.Fields{
			"event": "run_wda",
		}).Error("Failed running wda. Error: " + err.Error())
		fmt.Println(err.Error())
	}()

	return nil
}

func forwardPort(device ios.DeviceEntry, hostPort uint16, devicePort uint16) error {
	err := forward.Forward(device, hostPort, devicePort)
	if err != nil {
		return err
	}

	return nil
}

func updateWDA() error {
	fmt.Println("Updating WDA session")
	sessionID, err := createWDASession()
	if err != nil {
		return err
	}

	err = updateWdaStreamSettings(sessionID)
	if err != nil {
		return err
	}

	return nil
}

func updateWdaStreamSettings(sessionID string) error {
	requestString := `{"settings": {"mjpegServerFramerate": 30, "mjpegServerScreenshotQuality": 50, "mjpegScalingFactor": 100}}`

	response, err := http.Post("http://localhost:8100/session/"+sessionID+"/appium/settings", "application/json", strings.NewReader(requestString))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Could not successfully update WDA stream settings")
	}

	return nil
}

func createWDASession() (string, error) {
	requestString := `{
		"capabilities": {
			"firstMatch": [
				{
					"arguments": [],
					"environment": {},
					"eventloopIdleDelaySec": 0,
					"shouldWaitForQuiescence": true,
					"shouldUseTestManagerForVisibilityDetection": false,
					"maxTypingFrequency": 60,
					"shouldUseSingletonTestManager": true,
					"shouldTerminateApp": true,
					"forceAppLaunch": true,
					"useNativeCachingStrategy": true,
					"forceSimulatorSoftwareKeyboardPresence": false
				}
			],
			"alwaysMatch": {}
		}
	}`

	response, err := http.Post("http://localhost:8100/session", "application/json", strings.NewReader(requestString))
	if err != nil {
		return "", err
	}

	responseBody, _ := io.ReadAll(response.Body)

	var responseJson map[string]interface{}
	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		return "", err
	}

	if responseJson["sessionId"] == "" {
		if err != nil {
			return "", errors.New("Could not get `sessionId` while creating a new WebDriverAgent session")
		}
	}

	return fmt.Sprintf("%v", responseJson["sessionId"]), nil
}
