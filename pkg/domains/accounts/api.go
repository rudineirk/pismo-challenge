package accounts

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
)

type httpHandler struct {
	service Service
}

func SetupHTTPRoutes(router *gin.Engine, service Service) {
	handler := httpHandler{
		service: service,
	}

	routeGroup := router.Group("/accounts")
	routeGroup.POST("", handler.CreateAccount)
	routeGroup.GET("/:account_id", handler.GetAccountByID)
}

func (handler *httpHandler) CreateAccount(ctx *gin.Context) {
	req := CreateAccountRequest{}
	if err := ctx.BindJSON(&req); err != nil {
		return
	}

	account, err := handler.service.CreateAccount(ctx, &req)

	if err != nil {
		if errors.Is(err, ErrInvalidDocumentNumber(nil)) || errors.Is(err, errorlib.ErrInvalidPayload(nil)) {
			ctx.JSON(http.StatusBadRequest, err)
		} else if errors.Is(err, errorlib.ErrDuplicated(nil)) {
			ctx.JSON(http.StatusConflict, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		return
	}

	ctx.JSON(http.StatusCreated, NewAPIResponseFromEntity(account))
}

func (handler *httpHandler) GetAccountByID(ctx *gin.Context) {
	accountIDRaw := ctx.Param("account_id")

	accountID, err := strconv.ParseInt(accountIDRaw, 10, 64)
	if err != nil {
		ctx.Status(http.StatusNotFound)
		return
	}

	account, err := handler.service.GetAccountByID(ctx, accountID)
	if err != nil {
		if errors.Is(err, errorlib.ErrNotFound(nil)) {
			ctx.Status(http.StatusNotFound)
			return
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	ctx.JSON(http.StatusOK, NewAPIResponseFromEntity(account))
}

type AccountAPIResponse struct {
	AccountID      int64  `json:"account_id"`
	DocumentNumber string `json:"document_number"`
}

func NewAPIResponseFromEntity(account *Account) *AccountAPIResponse {
	return &AccountAPIResponse{
		AccountID:      account.ID,
		DocumentNumber: account.DocumentNumber,
	}
}
