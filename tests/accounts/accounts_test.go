package accounts_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"

	assert "github.com/stretchr/testify/require"

	"github.com/rudineirk/pismo-challenge/pkg/domains/accounts"
	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
	"github.com/rudineirk/pismo-challenge/pkg/infra/database"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter"
	"github.com/rudineirk/pismo-challenge/pkg/infra/logger"
	"github.com/rudineirk/pismo-challenge/pkg/utils/testutils"
)

const ContentTypeJSON = "application/json"

func TestAccountsAPIs(t *testing.T) {
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

	repo := accounts.NewRepository(bunDB)
	svc := accounts.NewService(repo)

	router := httprouter.NewRouter(logger, cfg.IsProduction)
	accounts.SetupHTTPRoutes(router, svc)

	server, client := testutils.MakeTestHTTPServer(router)
	defer server.Close()

	t.Run("POST /accounts", func(t *testing.T) {
		t.Run("should create a new account", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"document_number": "66895932070",
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/accounts", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			respData := accounts.AccountAPIResponse{}
			err = json.NewDecoder(resp.Body).Decode(&respData)
			assert.NoError(t, err)

			assert.NotEqual(t, int64(0), respData.AccountID)
			assert.Equal(t, "66895932070", respData.DocumentNumber)
		})

		t.Run("should return error if payload format is invalid", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"document_number": nil,
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/accounts", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			resp, err = client.Post(server.URL+"/accounts", ContentTypeJSON, strings.NewReader(`{"invalid":"json"`))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return error if document is duplicated", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"document_number": "24.885.962/0001-33",
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/accounts", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			respData := accounts.AccountAPIResponse{}
			err = json.NewDecoder(resp.Body).Decode(&respData)
			assert.NoError(t, err)

			assert.NotEqual(t, int64(0), respData.AccountID)
			assert.Equal(t, "24885962000133", respData.DocumentNumber)

			resp, err = client.Post(server.URL+"/accounts", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusConflict, resp.StatusCode)
		})
	})

	t.Run("GET /accounts/{id}", func(t *testing.T) {
		t.Run("should get account by ID", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"document_number": "52.987.490/0001-65",
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/accounts", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, resp.StatusCode)

			account := accounts.AccountAPIResponse{}
			err = json.NewDecoder(resp.Body).Decode(&account)
			assert.NoError(t, err)

			assert.NotEqual(t, int64(0), account.AccountID)

			resp, err = client.Get(fmt.Sprintf("%s/accounts/%d", server.URL, account.AccountID))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)

			respData := accounts.AccountAPIResponse{}
			err = json.NewDecoder(resp.Body).Decode(&respData)
			assert.NoError(t, err)

			assert.Equal(t, account.AccountID, respData.AccountID)
			assert.Equal(t, account.DocumentNumber, respData.DocumentNumber)
		})

		t.Run("should return not found if can't find account", func(t *testing.T) {
			resp, err := client.Get(server.URL + "/accounts/987")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)

			resp, err = client.Get(server.URL + "/accounts/abc-123")
			assert.NoError(t, err)
			assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		})
	})
}
