package android_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/avast/retry-go"
	"github.com/codeskyblue/go-sh"
	"github.com/shamanec/GADS-docker-server/config"
)

type appiumCapabilities struct {
	UDID           string `json:"appium:udid"`
	AutomationName string `json:"appium:automationName"`
	PlatformName   string `json:"platformName"`
	DeviceName     string `json:"appium:deviceName"`
}

func SetupDevice() {
	fmt.Println("INFO: Device setup")

	// Check if device is available to adb
	err := retry.Do(
		func() error {
			err := checkDeviceAvailable()
			if err != nil {
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	// Start minicap and wait for it to be up 5 seconds
	go startMinicap()
	time.Sleep(5 * time.Second)

	//Try to forward minicap to host container
	err = retry.Do(
		func() error {
			err := forwardMinicap()
			if err != nil {
				fmt.Println("This is error from forward minicap")
				return err
			}
			return nil
		},
		retry.Attempts(3),
		retry.Delay(3*time.Second),
	)
	if err != nil {
		panic(err)
	}

	// Start getting minicap stream after service was started and forwarded to host container
	go GetTCPStream(conn, imageChan)

	// Start the Appium server
	go startAppium()
}

// Check if the Android device is available to adb
func checkDeviceAvailable() error {
	fmt.Println("INFO: Checking if device is available to adb")

	output, err := sh.Command("adb", "devices").Output()
	if err != nil {
		return errors.New("Could not execute `adb devices`, err: " + err.Error())
	}

	// Check if we got the device UDID in the list of `adb devices`
	if strings.Contains(string(output), config.UDID) {
		return nil
	}

	return errors.New("Device with UDID=" + config.UDID + " was not available to adb")
}

// Forward minicap socket to the host container
func forwardMinicap() error {
	fmt.Println("INFO: Forwarding minicap connection to tcp:1313")

	err := sh.Command("adb", "forward", "tcp:1313", "localabstract:minicap").Run()
	if err != nil {
		return err
	}

	return nil
}

// Starts the Appium server on the device
func startAppium() {
	fmt.Println("INFO: Starting Appium server")

	// Create the Appium capabilities
	capabilities := appiumCapabilities{
		UDID:           config.UDID,
		AutomationName: "UiAutomator2",
		PlatformName:   "Android",
		DeviceName:     config.DeviceName,
	}
	// Marshal the capabilities into a json
	capabilitiesJson, err := json.Marshal(capabilities)
	if err != nil {
		panic(errors.New("Could not marshal Appium capabilities json, err: " + err.Error()))
	}

	// Create a json file for the capabilities
	capabilitiesFile, err := os.Create("/opt/capabilities.json")
	if err != nil {
		panic(err)
	}

	// Wrute the json byte slice to the json file created above
	_, err = capabilitiesFile.Write(capabilitiesJson)
	if err != nil {
		panic(err)
	}

	// Create file for the Appium logs
	outfile, err := os.Create("/opt/logs/appium.log")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()

	// Create new shell session and redirect Stdout and Stderr to the Appium logs file
	session := sh.NewSession()
	session.Stdout = outfile
	session.Stderr = outfile

	// Start the Appium server with default cli arguments and using default capabilities from the file created above
	err = session.Command("appium", "-p", "4723", "--log-timestamp", "--allow-cors", "--allow-insecure", "chromedriver_autodownload", "--default-capabilities", "/opt/capabilities.json").Run()
	if err != nil {
		panic(err)
	}
}

// Starts minicap service on the device
func startMinicap() {
	fmt.Println("INFO: Starting minicap")

	if config.RemoteControl == "true" {
		session := sh.NewSession()
		session.SetDir("/root/minicap")

		if config.MinicapHalfResolution == "true" {
			height, err := strconv.Atoi(config.AndroidScreenHeight)
			width, err := strconv.Atoi(config.AndroidScreenWidth)
			if err != nil {
				panic(err)
			}

			config.AndroidScreenHeight = strconv.Itoa(height / 2)
			config.AndroidScreenWidth = strconv.Itoa(width / 2)
		}

		// Discard Stdout so we don't constantly write to the container-server.log (if needed)
		//session.Stdout = io.Discard

		err := session.Command("./run.sh", "-P", config.ScreenSize+"@"+config.AndroidScreenWidth+"x"+config.AndroidScreenHeight+"/0").Start()
		if err != nil {
			panic(err)
		}

		err = session.Wait()
		if err != nil {
			panic(err)
		}
	}
}
