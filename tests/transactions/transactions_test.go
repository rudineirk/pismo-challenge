package transactions_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	assert "github.com/stretchr/testify/require"

	"github.com/rudineirk/pismo-challenge/pkg/domains/accounts"
	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
	"github.com/rudineirk/pismo-challenge/pkg/domains/transactions"
	"github.com/rudineirk/pismo-challenge/pkg/infra/config"
	"github.com/rudineirk/pismo-challenge/pkg/infra/database"
	"github.com/rudineirk/pismo-challenge/pkg/infra/httprouter"
	"github.com/rudineirk/pismo-challenge/pkg/infra/logger"
	"github.com/rudineirk/pismo-challenge/pkg/utils/testutils"
)

const ContentTypeJSON = "application/json"

func TestTransactionsAPIs(t *testing.T) {
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

	accountsRepo := accounts.NewRepository(bunDB)
	accountsSvc := accounts.NewService(accountsRepo)
	accounts.SetupHTTPRoutes(router, accountsSvc)

	transactionsRepo := transactions.NewRepository(bunDB)
	transactionsSvc := transactions.NewService(transactionsRepo, accountsSvc)
	transactions.SetupHTTPRoutes(router, transactionsSvc)

	server, client := testutils.MakeTestHTTPServer(router)
	defer server.Close()

	accountID, err := CreateAccount(server, client)
	assert.NoError(t, err)

	t.Run("POST /transactions", func(t *testing.T) {
		reqs := []map[string]any{
			{
				"account_id":        accountID,
				"operation_type_id": operationtypes.CashPurchaseType,
				"amount":            -1.15,
			},
			{
				"account_id":        accountID,
				"operation_type_id": operationtypes.InstallmentType,
				"amount":            -0.01,
			},
			{
				"account_id":        accountID,
				"operation_type_id": operationtypes.WithdrawType,
				"amount":            -0.09,
			},
			{
				"account_id":        accountID,
				"operation_type_id": operationtypes.PaymentType,
				"amount":            float64(1),
			},
		}

		for _, req := range reqs {
			t.Run("should create a new transaction", func(t *testing.T) {
				jsonPayload, err := json.Marshal(req)
				assert.NoError(t, err)

				resp, err := client.Post(server.URL+"/transactions", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
				assert.NoError(t, err)
				assert.Equal(t, http.StatusCreated, resp.StatusCode)

				respData := transactions.TransactionAPIResponse{}
				err = json.NewDecoder(resp.Body).Decode(&respData)
				assert.NoError(t, err)

				assert.NotEqual(t, int64(0), respData.TransactionID)
				assert.Equal(t, accountID, respData.AccountID)
				assert.Equal(t, req["operation_type_id"], respData.OperationTypeID)
				assert.Equal(t, req["amount"], respData.Amount)
				assert.WithinDuration(t, time.Now(), respData.EventDate, 20*time.Millisecond)
			})
		}

		t.Run("should return error if payload format is invalid", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"account_id":        nil,
				"operation_type_id": 123,
				"amount":            123,
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/transactions", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

			resp, err = client.Post(server.URL+"/transactions", ContentTypeJSON, strings.NewReader(`{"invalid":"json"`))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return error if account is invalid", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"account_id":        678,
				"operation_type_id": operationtypes.PaymentType,
				"amount":            1.1,
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/transactions", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return error if operation type is invalid", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"account_id":        accountID,
				"operation_type_id": 5,
				"amount":            1.1,
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/transactions", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})

		t.Run("should return error if amount is invalid", func(t *testing.T) {
			jsonPayload, err := json.Marshal(map[string]any{
				"account_id":        accountID,
				"operation_type_id": operationtypes.PaymentType,
				"amount":            -1.1,
			})
			assert.NoError(t, err)

			resp, err := client.Post(server.URL+"/transactions", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
			assert.NoError(t, err)
			assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		})
	})
}

func CreateAccount(server *httptest.Server, client *http.Client) (int64, error) {
	jsonPayload, err := json.Marshal(map[string]any{
		"document_number": "27935572003",
	})
	if err != nil {
		return 0, err
	}

	resp, err := client.Post(server.URL+"/accounts", ContentTypeJSON, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return 0, err
	}

	respData := accounts.AccountAPIResponse{}
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		return 0, err
	}

	return respData.AccountID, nil
}
