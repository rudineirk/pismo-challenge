package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter"
	"github.com/rudineirk/pismo-challenge/pkg/infra/logger"
	"github.com/rudineirk/pismo-challenge/pkg/infra/signalhandler"
)

func main() {
	cfg, err := config.LoadConfig()
	logger := logger.NewLogger(cfg)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load env config")
	}

	sighandler := signalhandler.NewSignalHandler(logger)
	router := httprouter.NewRouter(cfg, logger)
	httpserver := httprouter.NewServer(cfg, router)

	go sighandler.Listen(func(ctx context.Context) {
		if err := httpserver.Shutdown(ctx); err != nil {
			logger.Err(err).Msg("Server shutdown error")
		}
	})

	logger.Info().Msg(fmt.Sprintf("Starting server on http://%s", httpserver.Addr))

	if err = httpserver.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server stopped")
}
