package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"
)

const (
	SeverityInfo    = "i"
	SeverityWarning = "w"
	SeverityError   = "e"
)

var (
	application   string
	LogServiceUrl string
	debug         bool
	loggerInfo    chan interface{}
	loggerWarn    chan interface{}
	loggerError   chan interface{}
)

type LogMessage struct {
	Id          uuid.UUID   `json:"id" bson:"id"`
	Severity    string      `json:"severity" bson:"severity"`
	TimeStamp   int64       `json:"timestamp" bson:"timestamp"`
	Application string      `json:"application" bson:"application"`
	Content     interface{} `json:"content" bson:"content"`
}

// Instead of writing logs using the goroutine that is currently executing
// we will have a dedicated go routine that will handle logs that are being
// sent to a channel. This should be a major performance improvement since
// sending a request takes rather long.
//
// For backwards compatibility reasons the old functions
// will stay the same for now.
func init() {
	debug = os.Getenv("LOG_LEVEL") == "DEBUG"
	loggerInfo = make(chan interface{}, 100)
	loggerWarn = make(chan interface{}, 100)
	loggerError = make(chan interface{}, 100)
	go logger()
}

// This function listens to the three different channels where log messages
// can be sent. When a message is received the message will be logged and
// if loggingUrl is set the log will be sent there.
func logger() {
	for {
		var (
			s string
			c interface{}
		)
		select {
		case c = <-loggerInfo:
			s = SeverityInfo
		case c = <-loggerWarn:
			s = SeverityWarning
		case c = <-loggerError:
			s = SeverityError
		}
		log(c, s)
	}
}

// This can be used by go http servers that are following the standard signature to
// log some statistics about the requests
func RouterMiddleWare(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		LogInfo(map[string]interface{}{
			"method":        r.Method,
			"requestUri":    r.RequestURI,
			"remoteAddress": r.RemoteAddr,
			"duration":      time.Since(start),
		})
	})
}

// Writes a log to the console but also tries to send the log to a logging server
// if the LogServiceUrl is set. If the url is not set no data will be sent anywhere.
func log(content interface{}, severity string) {
	if application == "" {
		panic("application string not set")
	}
	res, _ := json.Marshal(LogMessage{
		Id:          uuid.New(),
		TimeStamp:   time.Now().Unix(),
		Application: application,
		Content:     content,
		Severity:    severity,
	})
	if debug {
		fmt.Println(string(res))
	}
	if parsedUrl, err := url.Parse(LogServiceUrl); err == nil && LogServiceUrl != "" {
		req, _ := http.NewRequest("POST", parsedUrl.String(), bytes.NewReader(res))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println("Error sending request: " + err.Error())
		} else if resp.StatusCode != http.StatusOK {
			fmt.Println("Received non 200 status code: " + strconv.Itoa(resp.StatusCode))
		}
	}
}

// Writes a log entry with severity: info. See log() for more details
func LogInfo(content interface{}) {
	loggerInfo <- content
}

// Writes a log entry with severity: warning. See log() for more details
func LogWarning(content interface{}) {
	loggerWarn <- content
}

// Writes a log entry with severity: error. See log() for more details
func LogError(content interface{}) {
	loggerError <- content
}

// Sets the application string. This must be done before any log function is
// called as the application string is part of the log message to see where
// a log entry came from.
func SetApplication(newApplication string) {
	application = newApplication
}
