package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

//=================//
//=====STRUCTS=====//

type JsonErrorResponse struct {
	EventName    string `json:"event"`
	ErrorMessage string `json:"error_message"`
}

type JsonResponse struct {
	Message string `json:"message"`
}

//=======================//
//=====API FUNCTIONS=====//

// Write to a ResponseWriter an event and message with a response code
func JSONError(w http.ResponseWriter, event string, error_string string, code int) {
	var errorMessage = JsonErrorResponse{
		EventName:    event,
		ErrorMessage: error_string}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(errorMessage)
}

// Write to a ResponseWriter an event and message with a response code
func SimpleJSONResponse(w http.ResponseWriter, response_message string, code int) {
	var message = JsonResponse{
		Message: response_message,
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(message)
}

// Prettify JSON with indentation and stuff
func PrettifyJSON(data string) string {
	var prettyJSON bytes.Buffer
	json.Indent(&prettyJSON, []byte(data), "", "  ")
	return prettyJSON.String()
}

// Convert interface into JSON string
func ConvertToJSONString(data interface{}) string {
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	return PrettifyJSON(string(b))
}

// Unmarshal provided JSON string into a struct
func UnmarshalJSONString(jsonString string, v interface{}) error {
	bs := []byte(jsonString)

	err := json.Unmarshal(bs, v)
	if err != nil {
		return err
	}

	return nil
}

// I've added the below structs in case we need to connect to Selenium Grid
// Needs to be handled though

type SeleniumGridJSON struct {
	GridCapabilities  []GridCapabilities `json:"capabilities"`
	GridConfiguration GridConfiguration  `json:"configuration"`
}

// BrowserName - device name or actual browser for the OS
// version - device OS version
// MaxInstances 1
// platform - Android or iOS
// deviceType phone
// platformName - Android or iOS
// platformVersion - device OS version
type GridCapabilities struct {
	BrowserName     string `json:"browserName"`
	Version         string `json:"version"`
	MaxInstances    string `json:"maxInstances"`
	Platform        string `json:"platform"`
	DeviceName      string `json:"deviceName"`
	DeviceType      string `json:"deviceType"`
	PlatformName    string `json:"platformName"`
	PlatformVersion string `json:"platformVersion"`
	UDID            string `json:"udid"`
}

// URL = http:// + devices host : + appium port /wd/hub
// timeout 180
// maxSession 1
// register true
/// registerCycle 5000
// automationName - UiAutomator2 for Android, XCUITest for iOS
// downPollingLimit 10
// hub protocol - http or https
type GridConfiguration struct {
	URL              string `json:"url"`
	AppiumPort       string `json:"port"`
	DeviceHost       string `json:"host"`
	GridPort         string `json:"hubPort"`
	GridHost         string `json:"hubHost"`
	Timeout          string `json:"timeout"`
	MaxSession       string `json:"maxSession"`
	Register         string `json:"register"`
	RegisterCycle    string `json:"registerCycle"`
	AutomationName   string `json:"automationName"`
	DownpollingLimit string `json:"downPollingLimit"`
	GridProtocol     string `json:"hubProtocol"`
}
