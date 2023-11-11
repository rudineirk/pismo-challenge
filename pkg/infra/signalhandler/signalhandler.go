package signalhandler

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rs/zerolog"
)

type SignalHandler struct {
	signalChan chan os.Signal
	logger     *zerolog.Logger
}

func NewSignalHandler(logger *zerolog.Logger) *SignalHandler {
	return &SignalHandler{
		signalChan: make(chan os.Signal, 1),
		logger:     logger,
	}
}

func (s *SignalHandler) Listen(handler func(context.Context)) {
	signal.Notify(s.signalChan, syscall.SIGINT, syscall.SIGTERM)

	<-s.signalChan
	s.logger.Info().Msg("Shutting down server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	handler(ctx)
}
