package status_test

import (
	"io"
	"net/http"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter"
	"github.com/rudineirk/pismo-challenge/pkg/infra/logger"
	"github.com/rudineirk/pismo-challenge/pkg/utils/testutils"
)

func TestStatusApi(t *testing.T) {
	cfg := &config.Config{IsProduction: true}
	logger := logger.NewStubLogger()
	router := httprouter.NewRouter(cfg, logger)

	server, client := testutils.MakeTestHTTPServer(router)
	defer server.Close()

	t.Run("should return OK status if server is running", func(t *testing.T) {
		resp, err := client.Get(server.URL + "/status")
		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		body, err := io.ReadAll(resp.Body)
		assert.Nil(t, err)
		assert.Equal(t, `{"ok":true}`, string(body))
	})
}
