package accounts

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/rudineirk/pismo-challenge/pkg/utils/errorlib"
)

type httpHandler struct {
	service  Service
	validate *validator.Validate
}

func SetupHTTPRoutes(router *gin.Engine, service Service) {
	handler := httpHandler{
		service:  service,
		validate: validator.New(validator.WithRequiredStructEnabled()),
	}

	routeGroup := router.Group("/accounts")
	routeGroup.POST("", handler.CreateAccount)
	routeGroup.GET("/:account_id", handler.GetAccountByID)
}

func (handler *httpHandler) CreateAccount(ctx *gin.Context) {
	req := CreateAccountAPIRequest{}

	if err := ctx.BindJSON(req); err != nil {
		return
	}

	if err := handler.validate.Struct(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorlib.ErrInvalidPayload(nil))
	}

	account := Account{
		DocumentNumber: req.DocumentNumber,
	}

	if err := handler.service.CreateAccount(ctx, &account); err != nil {
		if errors.Is(err, ErrInvalidDocumentNumber(nil)) {
			ctx.JSON(http.StatusBadRequest, err)
		} else if errors.Is(err, errorlib.ErrDuplicated(nil)) {
			ctx.JSON(http.StatusConflict, err)
		} else {
			_ = ctx.AbortWithError(http.StatusInternalServerError, err)
		}

		return
	}

	ctx.JSON(http.StatusCreated, NewAPIResponseFromEntity(&account))
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
