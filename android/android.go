package android_server

import (
	"bufio"
	"bytes"
	"errors"
	"os/exec"
	"strings"

	log "github.com/sirupsen/logrus"
)

var udid = "WCR7N18B14002300"

func startAppiumAndroid() {

}

func startMinicap() {

}

func GetInstalledApps() ([]string, error) {
	commandString := "adb -s " + udid + " shell cmd package list packages -3"
	cmd := exec.Command("bash", "-c", commandString)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "",
		}).Error("Could not get third party package names. Error: " + err.Error())
		return nil, errors.New("Could not get third party package names")
	}

	var packageNames []string

	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	for scanner.Scan() {
		packageName := strings.SplitAfter(scanner.Text(), "package:")[1]
		packageNames = append(packageNames, packageName)
	}

	return packageNames, nil
}

func LaunchApp(packageName string) error {
	commandString := "adb -s " + udid + " shell monkey -p " + packageName + " 1"
	cmd := exec.Command("bash", "-c", commandString)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		log.WithFields(log.Fields{
			"event": "",
		}).Error("Could not start app with packageName: " + packageName + ". Error: " + err.Error())
		return errors.New("Could not start app with packageName" + packageName + " ")
	}

	return nil
}
