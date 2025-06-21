package logger

import (
	"log"
	"os"
)

var (
	Logger *log.Logger
)

func init() {
	file, err := os.OpenFile("debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
}

// Info logs an informational message
func Info(msg string) {
	Logger.Println("[INFO] " + msg)
}

// Error logs an error message
func Error(msg string) {
	Logger.Println("[ERROR] " + msg)
}

// Debug logs a debug message
func Debug(msg string) {
	Logger.Println("[DEBUG] " + msg)
}
