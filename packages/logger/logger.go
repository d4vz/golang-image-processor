package logger

import (
	"log/slog"
	"os"
)

func Info(message string) {
	slog.Info(message)
}

func Error(message string) {
	slog.Error(message)
}

func Debug(message string) {
	slog.Debug(message)
}

func Fatal(message string) {
	slog.Error(message)
	os.Exit(1)
}
