package stackify

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var stackify_hostname, _ = os.Hostname()

var StackifyApiKey = flag.String("stackify-apikey", "", "Stackify api key")
var StackifyEnv = flag.String("stackify-env", "dev", "Environment name (default: dev)")
var StackifyServer = flag.String("stackify-server", stackify_hostname, "Server name (default: hostname)")
var StackifyAppName = flag.String("stackify-appname", "", "The name of the application")
var StackifyAppLocation = flag.String("stackify-applocation", "", "The full directory path for the application (optional)")
var StackifyLogger = flag.String("stackify-logger", "stackify-client", "The name and version of the logging project generating this request (default: stackify-client)")
var StackifyPlatform = flag.String("stackify-platform", "go", "The logging language (default: go)")
var StackifyTimeout = flag.Duration("stackify-timeout", time.Duration(5) * time.Second, "Stackify POST timeout (default: 5s)")

type API struct {
	APIKey   string
	Client *http.Client
}

type Report struct {
	Env string             `json:"Env"`
	ServerName string      `json:"ServerName"`
	AppName string         `json:"AppName"`
	AppLocation string     `json:"AppLoc,omitempty"`
	Logger string          `json:"Logger"`
	Platform string        `json:"Platform"`
	Messages []*Event      `json:"Msgs"`
}

type Event struct {
	Message string      `json:"Msg"`                 // The log message
	Level EventType     `json:"Level"`               // The log level
	Time int64          `json:"EpochMs"`             // Unix/POSIX/Epoch time with millisecond precision

	ThreadName string   `json:"Th,omitempty"`        // The thread name
	SourceMethod string `json:"SrcMethod,omitempty"` // Fully qualified method name
	SourceLine int64    `json:"SrcLine,omitempty"`   // Line number
	TransactionId int64 `json:"TransID,omitempty"`   // Transaction identifier
	Data string         `json:"Data,omitempty"`      // Additional JSON metadata about the log event. Special characters need to be escaped.

	// Msgs[*]/Ex/* omitted
}

type EventType string

const (
	StackifyEndpoint = "https://api.stackify.com/Log/Save"
	StackifyVersion = "V1"

	Info            EventType = "INFO"
	Warning         EventType = "WARN"
	Error           EventType = "ERROR"
)

type Response struct {
	Result bool `json:"success"`     // True indicates success, false otherwise
	ResponseTime int64 `json:"took"` // Elapsed time in milliseconds
}

// https://github.com/stackify/stackify-api/blob/master/endpoints/POST_Log_Save.md
func NewClient() *API {
	a := &API{
		APIKey: *StackifyApiKey,
		Client: http.DefaultClient,
	}

	a.Client.Timeout = *StackifyTimeout

	return a
}

func CreateReport() *Report {
	return CreateReportFromMessages([]*Event{})
}

func CreateReportFromMessages(events []*Event) *Report {
	return &Report{
		Env: *StackifyEnv,
		ServerName: *StackifyServer,
		AppName: *StackifyAppName,
		AppLocation: *StackifyAppLocation,
		Logger: *StackifyLogger,
		Platform: *StackifyPlatform,
		Messages: events,
	}
}

func CreateEvent(level EventType, msg string) *Event {
	return &Event{
		Message: msg,
		Level: level,
		Time: time.Now().UnixNano() / int64(time.Millisecond),
	}
}

func (a *API) Send(report *Report) (*Response, error) {
	js, err := json.Marshal(report)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", StackifyEndpoint, bytes.NewBuffer(js))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Stackify-PV", StackifyVersion)
	req.Header.Set("X-Stackify-Key", a.APIKey)

	resp, err := a.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Bad response %v", resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result Response

	if err = json.Unmarshal(body, &result); err != nil {
		return nil, err
	}

	if !result.Result {
		return &result, fmt.Errorf("Failed to POST %v to Stackify", report)
	}

	return &result, nil
}
