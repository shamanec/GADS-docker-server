package android_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	fmt.Println("Device setup")

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

func forwardMinicap() error {
	fmt.Println("INFO: Forwarding minicap connection to tcp:1313")

	err := sh.Command("adb", "forward", "tcp:1313", "localabstract:minicap").Run()
	if err != nil {
		return err
	}

	return nil
}

func startAppium() {
	fmt.Println("INFO: Starting Appium server")
	capabilities1 := appiumCapabilities{
		UDID:           config.UDID,
		AutomationName: "UiAutomator2",
		PlatformName:   "Android",
		DeviceName:     config.DeviceName,
	}
	capabilitiesJson, err := json.Marshal(capabilities1)
	if err != nil {
		panic(errors.New("Could not marshal Appium capabilities json, err: " + err.Error()))
	}

	// We are using /bin/bash -c here because os.exec does not invoke the system shell and `nvm` is not sourced in the container
	// should find a better solution in the future
	cmd := exec.Command("/bin/bash", "-c", "appium -p 4723 --log-timestamp --allow-cors --session-override --allow-insecure chromedriver_autodownload --default-capabilities '"+string(capabilitiesJson)+"'")
	fmt.Println(cmd)

	outfile, err := os.Create("/opt/logs/appium.log")
	if err != nil {
		panic(err)
	}
	defer outfile.Close()
	cmd.Stdout = outfile
	cmd.Stderr = outfile

	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}

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
