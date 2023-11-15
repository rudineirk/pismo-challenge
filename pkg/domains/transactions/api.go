package transactions

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rudineirk/pismo-challenge/pkg/domains/operationtypes"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
)

type httpHandler struct {
	service Service
}

func SetupHTTPRoutes(router *gin.Engine, service Service) {
	handler := httpHandler{
		service: service,
	}

	routeGroup := router.Group("/transactions")
	routeGroup.POST("", handler.CreateTransaction)
}

func (handler *httpHandler) CreateTransaction(ctx *gin.Context) {
	req := CreateTransactionRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	account, err := handler.service.CreateTransaction(ctx, &req)

	if err != nil {
		isBadRequest := errors.Is(err, ErrAccountIDNotFound(nil)) ||
			errors.Is(err, ErrInvalidOperationTypeID(nil)) ||
			errors.Is(err, ErrInvalidAmount(nil)) ||
			errors.Is(err, errorlib.ErrInvalidPayload(nil))

		if isBadRequest {
			ctx.JSON(http.StatusBadRequest, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		return
	}

	ctx.JSON(http.StatusCreated, NewAPIResponseFromEntity(account))
}

type TransactionAPIResponse struct {
	TransactionID   int64               `json:"transaction_id"`
	AccountID       int64               `json:"account_id"`
	OperationTypeID operationtypes.Type `json:"operation_type_id"`
	Amount          float64             `json:"amount"`
	EventDate       time.Time           `json:"event_date"`
}

func NewAPIResponseFromEntity(transaction *Transaction) *TransactionAPIResponse {
	return &TransactionAPIResponse{
		TransactionID:   transaction.ID,
		AccountID:       transaction.AccountID,
		OperationTypeID: transaction.OperationTypeID,
		Amount:          transaction.Amount,
		EventDate:       transaction.EventDate,
	}
}
