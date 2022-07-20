package ios_server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/testmanagerd"
	"github.com/shamanec/GADS-docker-server/config"
	"github.com/shamanec/GADS-docker-server/helpers"
	log "github.com/sirupsen/logrus"
)

func setLogging() {
	log.SetFormatter(&log.JSONFormatter{})
	project_log_file, err := os.OpenFile(config.HomeDir+"/logs/wda-sync.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		panic(err)
	}
	log.SetOutput(project_log_file)
}

func StartWDA(w http.ResponseWriter, r *http.Request) {
	setLogging()
	err := StartWDAInternal()
	if err != nil {
		helpers.JSONError(w, "", err.Error(), 500)
	}
	fmt.Fprintf(w, "Attempting to start WDA on port: "+config.WdaPort)
}

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
