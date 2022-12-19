package ios_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/codeskyblue/go-sh"
	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/forward"
	"github.com/danielpaulus/go-ios/ios/imagemounter"
	"github.com/danielpaulus/go-ios/ios/testmanagerd"
	"github.com/shamanec/GADS-docker-server/config"
	log "github.com/sirupsen/logrus"
)

func SetupDevice() {
	fmt.Println("INFO: Device setup")

	// Start usbmuxd service since it does not start automatically in the container
	// in a goroutine because its a long-running process
	go startUsbmuxd()

	// Try to get ios.DeviceEntry after usbmuxd is started 3 times in 10 seconds
	// stop the server if no device available
	err := config.GetDevice()
	if err != nil {
		panic("Could not get device with go-ios after 3 attempts, err:" + err.Error())
	}

	// Try to pair the device (assumes it is supervised)
	err = retry.Do(
		func() error {
			err := pairDevice()
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(5*time.Second),
	)
	if err != nil {
		panic("Could not pair supervised device after 3 attempts, err:\n" + err.Error())
	}

	// Try to mount the developer disk images on the device
	err = retry.Do(
		func() error {
			err := mountDeveloperImage()
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		panic("Could not mount developer disk images after 3 attempts, err:\n" + err.Error())
	}

	// Try to install and start WebDriverAgent
	err = retry.Do(
		func() error {
			err := installAndStartWebDriverAgent()
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		panic("Could not install WebDriverAgent after 3 attempts, err:\n" + err.Error())
	}

	// TODO check how to filter through WebDriverAgent startup output to see when it is up instead of hardcoded sleep
	time.Sleep(15 * time.Second)

	// Try to forward WebDriverAgent to host container
	err = retry.Do(
		func() error {
			err := forwardPort(8100, 8100)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		panic("Could not forward WebDriverAgent port, err:\n" + err.Error())
	}

	// Try to forward WebDriverAgent mjpeg stream to host container
	err = retry.Do(
		func() error {
			err := forwardPort(9100, 9100)
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		panic("Could not forward WebDriverAgent mjpeg stream port, err:\n" + err.Error())
	}

	// TODO think if we should panic here or just leave it be
	// Try to update WebDriverAgent streaming settings
	err = retry.Do(
		func() error {
			err := updateWebDriverAgent()
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(2*time.Second),
	)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "update_webdriveragent",
		}).Error("Could not create WebDriverAgent session or update stream settings, err:" + err.Error())
	}

	// Start Appium server
	go startAppium()
}

// Start usbmuxd service after starting the container
func startUsbmuxd() {
	fmt.Println("INFO: Starting usbmuxd service")

	// Create a new shell session and discard Stdout so we dont spam container-server.log
	// Should only see Stderr in case something happens
	session := sh.NewSession()
	session.Stdout = io.Discard

	// Create a usbmuxd.log file for Stderr
	usbmuxdLog, err := os.Create("/opt/logs/usbmuxd.log")
	if err != nil {
		panic("Could not create /opt/logs/usbmuxd.log file, err:" + err.Error())
	}
	// Redirect session Stderr to the usbmuxd.log file created above
	session.Stderr = usbmuxdLog

	err = session.Command("usbmuxd", "-U", "usbmux", "-f").Run()
	if err != nil {
		panic("usbmuxd service failed, err:" + err.Error())
	}
}

// Pair the device which is expected to be supervised
func pairDevice() error {
	p12, err := os.ReadFile("/opt/supervision.p12")
	if err != nil {
		log.Warning("Could not read /opt/supervision.p12 file, err:" + err.Error())
		log.Info("Falling back to unsupervised pairing")
		err = ios.Pair(config.Device)
		if err != nil {
			return errors.New("Could not pair successfully, err:" + err.Error())
		}
	}

	err = ios.PairSupervised(config.Device, p12, config.SupervisionPassword)
	if err != nil {
		return errors.New("Could not pair successfully, err:" + err.Error())
	}

	return nil
}

// Mount the developer disk images downloading them automatically in /opt/DeveloperDiskImages
func mountDeveloperImage() error {
	err := imagemounter.FixDevImage(config.Device, "/opt/DeveloperDiskImages")
	if err != nil {
		return err
	}

	return nil
}

// Install WebDriverAgent on the device
func installAndStartWebDriverAgent() error {
	err := InstallAppWithDevice(config.Device, "WebDriverAgent.ipa")
	if err != nil {
		return err
	}

	go startWebDriverAgent()
	return nil
}

// Start the WebDriverAgent on the device
func startWebDriverAgent() {
	fmt.Println("INFO: Starting WebDriverAgent")
	outfile, err := os.Create("/opt/logs/wda.log")
	if err != nil {
		panic("Could not create /opt/logs/wda.log file, err:" + err.Error())
	}
	defer outfile.Close()

	session := sh.NewSession()
	session.Stdout = outfile
	session.Stderr = outfile

	// Lazy way to do this using go-ios binary, should some day update to use go-ios modules instead!!!
	err = session.Command("ios", "runwda", "--bundleid="+config.BundleID, "--testrunnerbundleid="+config.BundleID, "--xctestconfig=WebDriverAgentRunner.xctest", "--udid="+config.UDID).Run()
	if err != nil {
		panic("Running WebDriverAgent using go-ios binary failed, err:" + err.Error())
	}
}

// Forward a port from device to container using go-ios
func forwardPort(hostPort uint16, devicePort uint16) error {
	fmt.Printf("INFO: Forwarding port=%v to host port=%v \n", devicePort, hostPort)
	err := forward.Forward(config.Device, hostPort, devicePort)
	if err != nil {
		return err
	}

	return nil
}

// Create a new WebDriverAgent session and update stream settings
func updateWebDriverAgent() error {
	fmt.Println("INFO: Updating WebDriverAgent session and mjpeg stream settings")

	sessionID, err := createWebDriverAgentSession()
	if err != nil {
		return err
	}

	err = updateWebDriverAgentStreamSettings(sessionID)
	if err != nil {
		return err
	}

	return nil
}

func updateWebDriverAgentStreamSettings(sessionID string) error {
	// Set 30 frames per second, without any scaling, half the original screenshot quality
	// TODO should make this configurable in some way, although can be easily updated the same way
	requestString := `{"settings": {"mjpegServerFramerate": 30, "mjpegServerScreenshotQuality": 50, "mjpegScalingFactor": 100}}`

	// Post the mjpeg server settings
	response, err := http.Post("http://localhost:8100/session/"+sessionID+"/appium/settings", "application/json", strings.NewReader(requestString))
	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		return errors.New("Could not successfully update WDA stream settings, status code=" + strconv.Itoa(response.StatusCode))
	}

	return nil
}

// Create a new WebDriverAgent session
func createWebDriverAgentSession() (string, error) {
	// TODO see if this JSON can be simplified
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

	// Post to create new session
	response, err := http.Post("http://localhost:8100/session", "application/json", strings.NewReader(requestString))
	if err != nil {
		return "", err
	}

	// Get the response into a byte slice
	responseBody, _ := io.ReadAll(response.Body)

	// Unmarshal response into a basic map
	var responseJson map[string]interface{}
	err = json.Unmarshal(responseBody, &responseJson)
	if err != nil {
		return "", err
	}

	// Check the session ID from the map
	if responseJson["sessionId"] == "" {
		if err != nil {
			return "", errors.New("Could not get `sessionId` while creating a new WebDriverAgent session")
		}
	}

	return fmt.Sprintf("%v", responseJson["sessionId"]), nil
}

// Start WebDriverAgent directly using go-ios modules
// TODO see if this can be reworked to log somewhere else, currently unused
func startWDAInternal() error {
	go func() {
		err := testmanagerd.RunXCUIWithBundleIdsCtx(nil, config.BundleID,
			config.TestRunnerBundleID,
			config.XCTestConfig,
			config.Device,
			[]string{},
			[]string{"USE_PORT=" + config.WdaPort, "MJPEG_SERVER_PORT=" + config.WdaMjpegPort})

		if err != nil {
			log.WithFields(log.Fields{
				"event": "run_wda",
			}).Error("Failed running wda. Error: " + err.Error())
			panic(err)
		}
	}()

	return nil
}

type appiumCapabilities struct {
	UDID                  string `json:"appium:udid"`
	WdaMjpegPort          string `json:"appium:mjpegServerPort"`
	ClearSystemFiles      string `json:"appium:clearSystemFiles"`
	WdaURL                string `json:"appium:webDriverAgentUrl"`
	PreventWdaAttachments string `json:"appium:preventWDAAttachments"`
	SimpleIsVisibleCheck  string `json:"appium:simpleIsVisibleCheck"`
	WdaLocalPort          string `json:"appium:wdaLocalPort"`
	PlatformVersion       string `json:"appium:platformVersion"`
	AutomationName        string `json:"appium:automationName"`
	PlatformName          string `json:"platformName"`
	DeviceName            string `json:"appium:deviceName"`
	WdaLaunchTimeout      string `json:"appium:wdaLaunchTimeout"`
	WdaConnectionTimeout  string `json:"appium:wdaConnectionTimeout"`
}

// Start the Appium server for the device
func startAppium() {
	fmt.Println("INFO: Starting Appium server")

	capabilities1 := appiumCapabilities{
		UDID:                  config.UDID,
		WdaURL:                "http://localhost:8100",
		WdaMjpegPort:          "9100",
		WdaLocalPort:          "8100",
		WdaLaunchTimeout:      "120000",
		WdaConnectionTimeout:  "240000",
		ClearSystemFiles:      "false",
		PreventWdaAttachments: "true",
		SimpleIsVisibleCheck:  "false",
		AutomationName:        "XCUITest",
		PlatformName:          "iOS",
		DeviceName:            config.DeviceName,
	}
	capabilitiesJson, err := json.Marshal(capabilities1)
	if err != nil {
		panic("Could not marshal Appium capabilities json, err: " + err.Error())
	}

	// Create a json file for the capabilities
	capabilitiesFile, err := os.Create("/opt/capabilities.json")
	if err != nil {
		panic("Could not create /opt/capabilities.json file, err:" + err.Error())
	}

	// Wrute the json byte slice to the json file created above
	_, err = capabilitiesFile.Write(capabilitiesJson)
	if err != nil {
		panic("Could not write capabilities to /opt/capabilities.json, err:" + err.Error())
	}

	// Create file for the Appium logs
	outfile, err := os.Create("/opt/logs/appium.log")
	if err != nil {
		panic("Could not create /opt/logs/appium.log file, err:" + err.Error())
	}
	defer outfile.Close()

	// Create new shell session and redirect Stdout and Stderr to the Appium logs file
	session := sh.NewSession()
	session.Stdout = outfile
	session.Stderr = outfile

	// Start the Appium server with default cli arguments and using default capabilities from the file created above
	err = session.Command("appium", "-p", "4723", "--log-timestamp", "--allow-cors", "--default-capabilities", "/opt/capabilities.json").Run()
	if err != nil {
		panic("Appium server failed, err:" + err.Error())
	}
}
