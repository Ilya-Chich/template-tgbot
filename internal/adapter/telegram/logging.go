package telegram

import (
	"encoding/json"
	"log"
	"time"

	"gopkg.in/telebot.v3"
)

// ErrorLog represents an error log entry
type ErrorLog struct {
	Type       string `json:"type"`
	Message    string `json:"message"`
	UserID     *int64 `json:"user_id,omitempty"`
	ServerName string `json:"server_name,omitempty"`
	SessionID  string `json:"session_id,omitempty"`
	Error      string `json:"error"`
	Timestamp  string `json:"timestamp"`
	Component  string `json:"component,omitempty"`
}

// LogError creates and logs a structured error entry
func LogError(message string, err error, opts ...ErrorOption) {
	errorLog := ErrorLog{
		Type:      "ERROR",
		Message:   message,
		Error:     err.Error(),
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	// Apply options
	for _, opt := range opts {
		opt(&errorLog)
	}

	if jsonData, jsonErr := json.Marshal(errorLog); jsonErr == nil {
		log.Printf("%s", jsonData)
	} else {
		log.Printf("Failed to marshal error log: %v", jsonErr)
	}
}

// ErrorOption is a functional option for configuring error logs
type ErrorOption func(*ErrorLog)

// WithUserID sets the user ID for the error log
func WithUserID(userID int64) ErrorOption {
	return func(el *ErrorLog) {
		el.UserID = &userID
	}
}

// WithServerName sets the server name for the error log
func WithServerName(serverName string) ErrorOption {
	return func(el *ErrorLog) {
		el.ServerName = serverName
	}
}

// WithSessionID sets the session ID for the error log
func WithSessionID(sessionID string) ErrorOption {
	return func(el *ErrorLog) {
		el.SessionID = sessionID
	}
}

// WithComponent sets the component name for the error log
func WithComponent(component string) ErrorOption {
	return func(el *ErrorLog) {
		el.Component = component
	}
}

// LogFatal logs a fatal error and exits the program
func LogFatal(message string, err error, opts ...ErrorOption) {
	LogError(message, err, opts...)
	log.Fatal("Fatal error occurred - exiting")
}

// MessageLog represents a message log entry
type MessageLog struct {
	Type      string `json:"type"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	ChatID    int64  `json:"chat_id"`
	Text      string `json:"text"`
	Timestamp string `json:"timestamp"`
}

// CallbackLog represents a callback log entry
type CallbackLog struct {
	Type      string `json:"type"`
	UserID    int64  `json:"user_id"`
	Username  string `json:"username"`
	Data      string `json:"data"`
	Timestamp string `json:"timestamp"`
}

// LoggingMiddleware logs all incoming updates from Telegram
func LoggingMiddleware(next telebot.HandlerFunc) telebot.HandlerFunc {
	return func(c telebot.Context) error {
		now := time.Now().UTC().Format(time.RFC3339)

		// Log message details if it's a message
		if c.Message() != nil {
			msg := c.Message()
			messageLog := MessageLog{
				Type:      "MESSAGE",
				UserID:    msg.Sender.ID,
				Username:  msg.Sender.Username,
				ChatID:    msg.Chat.ID,
				Text:      msg.Text,
				Timestamp: now,
			}

			if jsonData, err := json.Marshal(messageLog); err == nil {
				log.Printf("%s", jsonData)
			} else {
				log.Printf("Failed to marshal message log: %v", err)
			}
		}

		if c.Callback() != nil {
			cb := c.Callback()
			callbackLog := CallbackLog{
				Type:      "CALLBACK",
				UserID:    cb.Sender.ID,
				Username:  cb.Sender.Username,
				Data:      cb.Data,
				Timestamp: now,
			}

			if jsonData, err := json.Marshal(callbackLog); err == nil {
				log.Printf("%s", jsonData)
			} else {
				log.Printf("Failed to marshal callback log: %v", err)
			}
		}

		// Continue to the next handler
		return next(c)
	}
}

// BotBlockedErrorMsg is the exact error message returned by telebot when a user blocks the bot
const BotBlockedErrorMsg = "telegram: bot was blocked by the user (403)"
const BotDeactivatedErrorMsg = "telegram: user is deactivated (403)"

// IsBotBlockedError checks if the error indicates that a user has blocked the bot
func IsBotBlockedError(err error) bool {
	if err == nil {
		return false
	}
	return err.Error() == BotBlockedErrorMsg || err.Error() == BotDeactivatedErrorMsg
}
