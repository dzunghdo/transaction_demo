package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"transaction_demo/app/apperr"
	"transaction_demo/app/usecase"
	"transaction_demo/app/usecase/dto"
)

type AccountHandler struct {
	BaseHandler
	accountUC usecase.AccountUC
}

func NewAccountHandler(accountUC usecase.AccountUC) *AccountHandler {
	return &AccountHandler{
		accountUC: accountUC,
	}
}

// CreateAccount creates a new account
// @Summary Create a new account
// @Description  Create a new account with the provided details.
// @Tags Account
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} apperr.AppError
// @Failure 404 {object} apperr.AppError
// @Failure 500 {object} apperr.AppError
// @Router /accounts [POST]
func (hdl *AccountHandler) CreateAccount(ctx *gin.Context) {
	var (
		req dto.AccountDTO
		err error
	)
	defer func() {
		if err != nil {
			hdl.RenderError(ctx, err)
		} else {
			hdl.RenderResponse(ctx, http.StatusCreated, nil, nil)
		}
	}()

	if err := ctx.ShouldBindJSON(&req); err != nil {
		hdl.RenderError(ctx, apperr.ErrInvalidInput.WithError(err))
		return
	}
	_, err = hdl.accountUC.Create(ctx, req)
}

// GetAccountBalance retrieves the balance of an account
// @Summary Get Account Balance
// @Description Retrieve the balance of an account by its ID.
// @Tags Account
// @Accept json
// @Produce json
// @Success 200 {object} dto.AccountDTO
// @Failure 400 {object} apperr.AppError
// @Failure 404 {object} apperr.AppError
// @Failure 500 {object} apperr.AppError
// @Router /accounts/{account_id} [GET]
func (hdl *AccountHandler) GetAccountBalance(ctx *gin.Context) {
	var (
		accountID uint64
		res       dto.AccountDTO
		err       error
	)
	defer func() {
		if err != nil {
			hdl.RenderError(ctx, err)
		} else {
			hdl.RenderResponse(ctx, http.StatusOK, res, nil)
		}
	}()

	accountIDStr := ctx.Param("account_id")
	accountID, err = strconv.ParseUint(accountIDStr, 10, 64)
	if err != nil {
		fmt.Println("Invalid account_id format", accountIDStr)
		err = apperr.ErrInvalidInput.WithError(err).WithMessage("Invalid account ID format")
		return
	}
	if accountID <= 0 {
		fmt.Println("Invalid account_id", accountID)
		err = apperr.ErrInvalidInput.WithMessage("Account ID must be a positive integer")
		return
	}

	res, err = hdl.accountUC.GetBalance(ctx, accountID)
}

// MakeTransaction  performs a transaction on an account
// @Summary Make a transaction
// @Description  Perform a transaction on an account, updating its balance.
// @Tags Transaction
// @Accept json
// @Produce json
// @Success 200
// @Failure 400 {object} apperr.AppError
// @Failure 404 {object} apperr.AppError
// @Failure 500 {object} apperr.AppError
// @Router /transaction [POST]
func (hdl *AccountHandler) MakeTransaction(context *gin.Context) {
	var (
		req dto.TransactionDTO
		err error
	)

	defer func() {
		if err != nil {
			hdl.RenderError(context, err)
		} else {
			hdl.RenderResponse(context, http.StatusCreated, nil, nil)
		}
	}()

	if err = context.ShouldBindJSON(&req); err != nil {
		err = apperr.ErrInvalidInput.WithError(err).WithMessage("Invalid request body")
		return
	}

	err = hdl.accountUC.MakeTransaction(context, req)
}
