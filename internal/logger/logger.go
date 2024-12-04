package logger

import "github.com/razvanmarinn/chatroom/internal/cfg"

type Logger interface {
	Info(message string, args ...interface{})
    Warn(message string, args ...interface{})
    Error(message string, args ...interface{})
	Log(level string, message string, args ...interface{})
}

func NewLogger(config cfg.Config) Logger {
	switch config.LogType {
	case "local":
		return NewLocalLogger()
	case "centralized":
		return NewRabbitMQLogger()
	}
	return nil
}
