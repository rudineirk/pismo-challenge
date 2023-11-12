package status_test

import (
	"fmt"
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
	assert.Nil(t, err)

	logger := logger.NewStubLogger()

	cfg, err := config.LoadConfig()
	assert.Nil(t, err)

	cfg.IsProduction = true

	testDB, err := testutils.NewTestDatabase(cfg.DatabaseURL)
	assert.Nil(t, err)

	defer testDB.Drop()

	fmt.Println(testDB.URL)

	db, bunDB, err := database.NewDatabase(testDB.URL)
	assert.Nil(t, err)

	err = database.RunMigrations(db)
	assert.Nil(t, err)

	router := httprouter.NewRouter(logger, cfg.IsProduction)
	healthcheck.SetupHealthCheck(router, bunDB)

	server, client := testutils.MakeTestHTTPServer(router)
	defer server.Close()

	t.Run("/status", func(t *testing.T) {
		t.Run("should return OK status if server is running", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/status")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"ok":true}`, string(body))
		})
	})

	t.Run("/healthcheck/liveliness", func(t *testing.T) {
		t.Run("should return OK status if server is running", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/status")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"ok":true}`, string(body))
		})
	})

	t.Run("/healthcheck/readiness", func(t *testing.T) {
		t.Run("should return OK status if database connection is active", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/status")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"ok":true}`, string(body))
		})

		t.Run("should return OK status if database connection is closed", func(t *testing.T) {
			err := db.Close()
			assert.Nil(t, err)

			resp, err := client.Get(server.URL + "/status")
			assert.Nil(t, err)
			assert.Equal(t, http.StatusServiceUnavailable, resp.StatusCode)

			body, err := io.ReadAll(resp.Body)
			assert.Nil(t, err)
			assert.Equal(t, `{"ok":false}`, string(body))
		})
	})
}
