package v1

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

// @Summary Transaction
// @Description Pre Transaction from internal account to internal account
// @Tags Transaction
// @Accept json
// @Param request body model.PreInternalTransferRequest true "Transaction payload"
// @Produce  json
// @Router /transaction/pre-internal-transfer [post]
// @Success 200 {object} httpcommon.HttpResponse[string]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TransactionHandler) PreInternalTransfer(ctx *gin.Context) {
	var transfer model.PreInternalTransferRequest

	if err := validation.BindJsonAndValidate(ctx, &transfer); err != nil {
		return
	}
	transactionId, err := handler.transactionService.PreInternalTransfer(ctx, transfer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[string](&transactionId))
}

// @Summary Transaction
// @Description Verify OTP and transaction from internal account to internal account
// @Tags Transaction
// @Accept json
// @Param request body model.InternalTransferRequest true "Transaction payload"
// @Produce  json
// @Router /transaction/internal-transfer [post]
// @Success 200 {object} httpcommon.HttpResponse[entity.Transaction]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TransactionHandler) InternalTransfer(ctx *gin.Context) {
	var transfer model.InternalTransferRequest

	if err := validation.BindJsonAndValidate(ctx, &transfer); err != nil {
		return
	}
	transaction, err := handler.transactionService.InternalTransfer(ctx, transfer)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[*entity.Transaction](&transaction))
}

// @Summary Transaction
// @Description Add new Debt reminder
// @Tags Transaction
// @Accept json
// @Param request body model.DebtReminderRequest true "Transaction payload"
// @Produce  json
// @Router /transaction/debt-reminder [post]
// @Success 200 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TransactionHandler) AddDebtReminder(ctx *gin.Context) {
	var req model.DebtReminderRequest
	if err := validation.BindJsonAndValidate(ctx, &req); err != nil {
		return
	}
	err := handler.transactionService.AddDebtReminder(ctx, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(200)
}

// @Summary Debt reminder
// @Description cancel a debt reminder from source or target user
// @Tags Transaction
// @Accept json
// @Param id query string true "Id of debt reminder"
// @Param request body model.DebtReminderReplyRequest true "Transaction payload"
// @Produce  json
// @Router /transaction/cancel-debt-reminder/:id [PUT]
// @Success 200 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *TransactionHandler) CancelDebtReminder(ctx *gin.Context) {
	debtReminderId := ctx.Param(`id`)
	if debtReminderId == "" {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "Missing id Parameter",
				Code:    httpcommon.ErrorResponseCode.MissingIdParameter,
				Field:   "id",
			},
		))
		return
	}
	var req model.DebtReminderReplyRequest
	if err := validation.BindJsonAndValidate(ctx, &req); err != nil {
		return
	}
	err := handler.transactionService.CancelDebtReminder(ctx, debtReminderId, req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(200)
}
