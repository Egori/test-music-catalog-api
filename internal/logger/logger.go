// Package logger provides logging functionality
package logger

import (
	"log"
)

// Logger — интерфейс для логирования
type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Error(args ...interface{})
}

// LogLevel type for log levels
type LogLevel int

const (
	DEBUG LogLevel = iota
	INFO
	ERROR
	FATAL
)

// StandardLogger реализация логгера с использованием стандартного логирования
type StandardLogger struct {
	level LogLevel
}

// NewLogger creates a new Logger instance based on log level from environment
func NewLogger(levelStr string) *StandardLogger {
	level := getLogLevel(levelStr)
	return &StandardLogger{level: level}
}

// getLogLevel fetches log level from environment or defaults to INFO
func getLogLevel(levelStr string) LogLevel {
	switch levelStr {
	case "debug":
		return DEBUG
	case "error":
		return ERROR
	case "fatal":
		return FATAL
	default:
		return INFO
	}
}

// Debug logs
func (l *StandardLogger) Debug(v ...interface{}) {
	if l.level <= DEBUG {
		log.Println("[DEBUG]", v)
	}
}

// Info logs
func (l *StandardLogger) Info(v ...interface{}) {
	if l.level <= INFO {
		log.Println("[INFO]", v)
	}
}

// Error logs
func (l *StandardLogger) Error(v ...interface{}) {
	log.Println("[ERROR]", v)
}

// Fatal logs and exits
func (l *StandardLogger) Fatal(v ...interface{}) {
	log.Fatal("[FATAL]", v)
}
