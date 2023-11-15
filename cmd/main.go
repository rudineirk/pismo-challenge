package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/rudineirk/pismo-challenge/pkg/domains/accounts"
	"github.com/rudineirk/pismo-challenge/pkg/domains/transactions"
	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
	"github.com/rudineirk/pismo-challenge/pkg/infra/database"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter/healthcheck"
	"github.com/rudineirk/pismo-challenge/pkg/infra/logger"
	"github.com/rudineirk/pismo-challenge/pkg/infra/signalhandler"
)

func main() {
	cfg, err := config.LoadConfig()
	logger := logger.NewLogger(cfg.LogFormat, cfg.LogLevel)

	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to load env config")
	}

	sqlDB, bunDB, err := database.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to connect to database")
	}

	defer sqlDB.Close()

	if err = database.RunMigrations(sqlDB); err != nil {
		logger.Fatal().Err(err).Msg("Failed to run database migrations")
	}

	sighandler := signalhandler.NewSignalHandler(logger)
	router := httprouter.NewRouter(logger, cfg.IsProduction)
	httpserver := httprouter.NewServer(cfg.HTTPPort, router)

	healthcheck.SetupHealthCheck(router, bunDB)

	go sighandler.Listen(func(ctx context.Context) {
		if err := httpserver.Shutdown(ctx); err != nil {
			logger.Err(err).Msg("Server shutdown error")
		}
	})

	accountsRepo := accounts.NewRepository(bunDB)
	accountsSvc := accounts.NewService(accountsRepo)
	accounts.SetupHTTPRoutes(router, accountsSvc)

	transactionsRepo := transactions.NewRepository(bunDB)
	transactionsSvc := transactions.NewService(transactionsRepo, accountsSvc)
	transactions.SetupHTTPRoutes(router, transactionsSvc)

	logger.Info().Msg(fmt.Sprintf("Starting server on http://localhost%s", httpserver.Addr))

	if err = httpserver.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	logger.Info().Msg("Server stopped")
}
