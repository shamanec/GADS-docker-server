package android_server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
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
	time.Sleep(10 * time.Second)

	// Try to forward minicap to host container
	// err = retry.Do(
	// 	func() error {
	// 		err := forwardMinicap()
	// 		if err != nil {
	// 			return err
	// 		}
	// 		return nil
	// 	},
	// 	retry.Attempts(3),
	// 	retry.Delay(3*time.Second),
	// )
	// if err != nil {
	// 	panic(err)
	// }

	go forwardMinicap()

	go startAppium()
}

func checkDeviceAvailable() error {
	output, err := exec.Command("adb", "devices").Output()
	if err != nil {
		return errors.New("Could not execute `adb devices`, err: " + err.Error())
	}

	if strings.Contains(string(output), config.UDID) {
		return nil
	}

	return errors.New("Device with UDID=" + config.UDID + " was not available to adb")
}

func forwardMinicap() error {
	// fmt.Println("Forwarding minicap")
	// cmd := exec.Command("adb", "forward", "tcp:1313", "localabstract:minicap")

	// err := cmd.Run()
	// if err != nil {
	// 	return errors.New("Could not forward minicap socket, err: " + err.Error())
	// }

	// return nil

	session := sh.NewSession()
	err := session.Command("adb forward tcp:1313 localabstract:minicap").Run()
	if err != nil {
		return err
	}

	return nil
}

func startAppium() {
	fmt.Println("Starting appium")
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

// cd /root/minicap/ && ./run.sh -P ${SCREEN_SIZE}@${STREAM_WIDTH}x${STREAM_HEIGHT}/0 >>/opt/logs/minicap.log 2>&1 &
func startMinicap() {
	fmt.Println("Starting minicap")
	if config.RemoteControl == "true" {
		session := sh.NewSession()
		session.SetDir("/root/minicap")

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
