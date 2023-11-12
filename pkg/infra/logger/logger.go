package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
)

func NewStubLogger() *zerolog.Logger {
	return NewLogger("", "disabled")
}

func FromCfg(cfg *config.Config) *zerolog.Logger {
	return NewLogger(cfg.LogFormat, cfg.LogLevel)
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
