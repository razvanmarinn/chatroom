package logger

import (
	"log"
	"time"

	"github.com/razvanmarinn/chatroom/internal/queue"
)

type RabbitMQLogger struct {
	queue *queue.RabbitMQ
}

func (l *RabbitMQLogger) Log(level string, message string, args ...interface{}) {
	logMessage := "[" + level + "] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	err := l.queue.Publish(logMessage)
	if err != nil {
		log.Printf("Failed to publish log: %v", err)
	}
}

func (l *RabbitMQLogger) Info(message string, args ...interface{}) {
	logMessage := "[INFO] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	err := l.queue.Publish(logMessage)
	if err != nil {
		log.Printf("Failed to publish log: %v", err)
	}
}
func (l *RabbitMQLogger) Warn(message string, args ...interface{}) {
	logMessage := "[WARN] " + "[" + time.Now().Format(time.RFC3339) + "]" + message
	err := l.queue.Publish(logMessage)
	if err != nil {
		log.Printf("Failed to publish log: %v", err)
	}
}
func (l *RabbitMQLogger) Error(message string, args ...interface{}) {
	logMessage := "[ERROR] " + "[" + time.Now().Format(time.RFC3339) + "] " + message
	for _, arg := range args {
		if err, ok := arg.(error); ok {
			logMessage += " | Error: " + err.Error()
			break
		}
	}

	err := l.queue.Publish(logMessage)
	if err != nil {
		log.Printf("Failed to publish log: %v", err)
	}
}

func NewRabbitMQLogger() *RabbitMQLogger {
	return &RabbitMQLogger{queue: queue.NewRabbitMQ()}
}
