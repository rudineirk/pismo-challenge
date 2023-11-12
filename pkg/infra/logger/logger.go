package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func NewStubLogger() *zerolog.Logger {
	return NewLogger("", "disabled")
}

func NewLogger(logFormat string, logLevel string) *zerolog.Logger {
	var output io.Writer

	switch logFormat {
	case "cli":
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	default:
		output = os.Stdout
	}

	level, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(output).With().Timestamp().Logger().Level(level)
	if level == zerolog.DebugLevel {
		logger = logger.With().Caller().Logger()
	}

	return &logger
}
