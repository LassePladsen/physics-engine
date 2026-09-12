package logger

import (
	"fmt"
	"log/slog"
	"os"
	"strings"
)

// Logger has a usable default so packages can log safely when they are used
// outside the main executable, such as in tests.
var Logger = slog.Default()
var LoggerLevel = slog.LevelInfo

func Init() {
	// Show debug logging with `LOG_LEVEL=DEBUG go run ...`
	if strings.EqualFold(os.Getenv("LOG_LEVEL"), "debug") {
		LoggerLevel = slog.LevelDebug
	}

	Logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: LoggerLevel,
	}))
}

func Debug(msg string) {
	Logger.Debug(msg)
}

func Info(msg string) {
	Logger.Info(msg)
}

func Error(msg string) {
	Logger.Error(msg)
}

func Debugf(msg string, vars ...any) {
	Logger.Debug(fmt.Sprintf(msg, vars...))
}

func Infof(msg string, vars ...any) {
	Logger.Info(fmt.Sprintf(msg, vars...))
}

func Errorf(msg string, vars ...any) {
	Logger.Error(fmt.Sprintf(msg, vars...))
}
