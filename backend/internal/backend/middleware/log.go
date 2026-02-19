package middleware

import (
	"fmt"
	"io"
	"log/slog"

	"github.com/dvalkoff/gomessenger/internal/backend/config"
)

const (
	JSONLoggingHandler = "json"
	TextLoggingHandler = "text"
)

const (
	serviceNameArgKey = "serviceName"
)

var logLevelMap map[string]slog.Level = map[string]slog.Level{
	"DEBUG": slog.LevelDebug,
	"INFO":  slog.LevelInfo,
	"WARN":  slog.LevelWarn,
	"ERROR": slog.LevelError,
}

func InitLogger(writer io.Writer, loggingConfig config.LoggingConfig) error {
	level, ok := logLevelMap[loggingConfig.Level]
	if !ok {
		return fmt.Errorf("logging level is not provided or incorrect: %s", loggingConfig.Level)
	}
	handlerOpts := &slog.HandlerOptions{
		AddSource: loggingConfig.AddSource,
		Level:     level,
	}
	var handler slog.Handler
	switch loggingConfig.Handler {
	case JSONLoggingHandler:
		handler = slog.NewJSONHandler(writer, handlerOpts)
	case TextLoggingHandler:
		handler = slog.NewTextHandler(writer, handlerOpts)
	default:
		return fmt.Errorf("logging handler is not provided or incorrect: %s", loggingConfig.Handler)
	}

	logger := slog.New(handler).With(serviceNameArgKey, loggingConfig.ServiceName)
	slog.SetDefault(logger)
	return nil
}
