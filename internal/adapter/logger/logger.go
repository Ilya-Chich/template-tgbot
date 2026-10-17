package logger

import (
	"encoding/json"
	"log"
	"time"

	"tgbot/internal/core/port"
)

type Logger struct{}

func New() *Logger { return &Logger{} }

type errorEntry struct {
	Type      string `json:"type"`
	Message   string `json:"message"`
	Error     string `json:"error"`
	Component string `json:"component,omitempty"`
	Timestamp string `json:"timestamp"`
}

func (l *Logger) Error(message string, err error, component string) {
	entry := errorEntry{"ERROR", message, err.Error(), component, time.Now().UTC().Format(time.RFC3339)}
	if data, e := json.Marshal(entry); e == nil {
		log.Printf("%s", data)
	} else {
		log.Printf("%s: %v", message, err)
	}
}

func (l *Logger) Fatal(message string, err error, component string) {
	l.Error(message, err, component)
	log.Fatal("fatal error")
}

var _ port.Logger = (*Logger)(nil)
