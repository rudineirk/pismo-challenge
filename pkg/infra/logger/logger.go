package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
)

func NewLogger(cfg *config.Config) *zerolog.Logger {
	var output io.Writer

	switch cfg.LogFormat {
	case "cli":
		output = zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
	default:
		output = os.Stdout
	}

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		level = zerolog.InfoLevel
	}

	logger := zerolog.New(output).With().Timestamp().Logger().Level(level)
	if level == zerolog.DebugLevel {
		logger = logger.With().Caller().Logger()
	}

	return &logger
}
