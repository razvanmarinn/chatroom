package logger

import (
	"fmt"
	"time"
)

type LocalLogger struct {
}

func (l *LocalLogger) Log(level string, message string, args ...interface{}) {
	logMessage := "[" + level + "] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	fmt.Println(logMessage)

}

func (l *LocalLogger) Info(message string, args ...interface{}) {
	logMessage := "[INFO] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	fmt.Println(logMessage)

}
func (l *LocalLogger) Warn(message string, args ...interface{}) {
	logMessage := "[WARN] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	fmt.Println(logMessage)

}
func (l *LocalLogger) Error(message string, args ...interface{}) {
	logMessage := "[ERROR] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	fmt.Println(logMessage)

}

func NewLocalLogger() *LocalLogger {
	return &LocalLogger{}
}
