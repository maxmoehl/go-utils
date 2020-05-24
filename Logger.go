// This package provides basic types and functions to support
// the use of github.com/maxmoehl/log-service.
package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"time"
)

const (
	Info    = "i"
	Warning = "w"
	Error   = "e"
)

var (
	application   string
	LogServiceUrl string
)

// This function returns a middleware that will write logs to the log service
func LoggingHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log("logger was used", Info)
		ctx.Next()
	}
}

type LogMessage struct {
	Id          uuid.UUID   `json:"_id" bson:"_id"`
	Severity    string      `json:"_severity" bson:"_severity"`
	TimeStamp   int64       `json:"_timestamp" bson:"_timestamp"`
	Application string      `json:"_application" bson:"_application"`
	Message     interface{} `json:"message" bson:"message"`
}

func log(message interface{}, severity string) {
	if application == "" {
		panic("application string not set")
	}
	if LogServiceUrl == "" {
		panic("log service url string not set")
	}
	res, _ := json.Marshal(LogMessage{
		Id:          uuid.New(),
		TimeStamp:   time.Now().Unix(),
		Application: application,
		Message:     message,
		Severity:    severity,
	})
	fmt.Println(string(res))
	req, _ := http.NewRequest("POST", LogServiceUrl, bytes.NewReader(res))
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	} else if resp.StatusCode != http.StatusOK {
		panic(resp.StatusCode)
	}
}

func LogInfo(message string) {
	log(message, Info)
}

func LogWarning(message string) {
	log(message, Warning)
}

func LogError(message string) {
	log(message, Error)
}

func SetApplication(newApplication string) {
	application = newApplication
}
