package status_test

import (
	"io"
	"net/http"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
	"github.com/rudineirk/pismo-challenge/pkg/infra/database"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter/healthcheck"
	"github.com/rudineirk/pismo-challenge/pkg/infra/logger"
	"github.com/rudineirk/pismo-challenge/pkg/utils/testutils"
)

func TestStatusApis(t *testing.T) {
	err := testutils.SetRootCwd()
	assert.NoError(t, err)

	logger := logger.NewStubLogger()

	cfg, err := config.LoadConfig()
	assert.NoError(t, err)

	cfg.IsProduction = true

	testDB, err := testutils.NewTestDatabase(cfg.DatabaseURL)
	assert.NoError(t, err)

	defer testDB.Drop()

	sqlDB, bunDB, err := database.NewDatabase(testDB.URL)
	assert.NoError(t, err)

	err = database.RunMigrations(sqlDB)
	assert.NoError(t, err)

	router := httprouter.NewRouter(logger, cfg.IsProduction)
	healthcheck.SetupHealthCheck(router, bunDB)

	server, client := testutils.MakeTestHTTPServer(router)
	defer server.Close()

	t.Run("GET /status", func(t *testing.T) {
		t.Run("should return OK status if server is running", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/status")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"ok":true}`, string(body))
		})
	})

	t.Run("GET /healthcheck/liveliness", func(t *testing.T) {
		t.Run("should return OK status if server is running", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/healthcheck/liveliness")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"ok":true}`, string(body))
		})
	})

	t.Run("GET /healthcheck/readiness", func(t *testing.T) {
		t.Run("should return OK status if database connection is active", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/healthcheck/readiness")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"ok":true}`, string(body))
		})

		t.Run("should return OK status if database connection is closed", func(t *testing.T) {
			err := sqlDB.Close()
			assert.NoError(t, err)

			resp, err := client.Get(server.URL + "/healthcheck/readiness")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.JSONEq(t, `{"ok":false}`, string(body))
		})
	})
}
