package ios_server

import (
	"fmt"
	"net/http"
	"os"

	"github.com/danielpaulus/go-ios/ios"
	"github.com/danielpaulus/go-ios/ios/testmanagerd"
	"github.com/shamanec/GADS-docker-server/config"
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
	device, err := ios.GetDevice(udid)
	if err != nil {
		log.WithFields(log.Fields{
			"event": "run_wda",
		}).Error("Could not get device when installing app. Error: " + err.Error())
	}

	go func() {
		err := testmanagerd.RunXCUIWithBundleIdsCtx(nil, bundleid,
			testrunnerbundleid,
			xctestconfig,
			device,
			[]string{},
			[]string{"USE_PORT=" + wda_port, "MJPEG_SERVER_PORT=" + wda_mjpeg_port})

		log.WithFields(log.Fields{
			"event": "run_wda",
		}).Error("Failed running wda. Error: " + err.Error())
		fmt.Println(err.Error())
	}()

	fmt.Fprintf(w, "Attempting to start WDA on port: "+wda_port)
}
