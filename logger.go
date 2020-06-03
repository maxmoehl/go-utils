package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"net/url"
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
)

// This can be used by go http servers that are following the standard signature to
// log some statistics about the requests
func RouterMiddleWare(inner http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log(fmt.Sprintf("%s\t%s\t%s", r.Method, r.RequestURI, time.Since(start)), SeverityInfo)
	})
}

type LogMessage struct {
	Id          uuid.UUID   `json:"_id" bson:"_id"`
	Severity    string      `json:"_severity" bson:"_severity"`
	TimeStamp   int64       `json:"_timestamp" bson:"_timestamp"`
	Application string      `json:"_application" bson:"_application"`
	Message     interface{} `json:"message" bson:"message"`
}

// Writes a log to the console but also tries to send the log to a logging server
// if the LogServiceUrl is set. If the url is not set no data will be sent anywhere.
func log(message interface{}, severity string) {
	if application == "" {
		panic("application string not set")
	}
	res, _ := json.Marshal(LogMessage{
		Id:          uuid.New(),
		TimeStamp:   time.Now().Unix(),
		Application: application,
		Message:     message,
		Severity:    severity,
	})
	fmt.Println(string(res))
	if parsedUrl, err := url.Parse(LogServiceUrl); err == nil && LogServiceUrl != "" {
		req, _ := http.NewRequest("POST", parsedUrl.String(), bytes.NewReader(res))
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			panic(err)
		} else if resp.StatusCode != http.StatusOK {
			panic(resp.StatusCode)
		}
	}
}

// Writes a log entry with an info severity. See log() for more details
func LogInfo(message string) {
	log(message, SeverityInfo)
}

// Writes a log entry with an warning severity. See log() for more details
func LogWarning(message string) {
	log(message, SeverityWarning)
}

// Writes a log entry with an error severity. See log() for more details
func LogError(message string) {
	log(message, SeverityError)
}

// Sets the application string. This must be done before any log function is
// called as the application string is part of the log message to see where
// a log entry came from.
func SetApplication(newApplication string) {
	application = newApplication
}
