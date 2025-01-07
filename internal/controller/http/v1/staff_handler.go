package v1

import (
	"net/http"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
)

type StaffHandler struct {
	staffService service.StaffService
}

func NewStaffHandler(staffService service.StaffService) *StaffHandler {
	return &StaffHandler{staffService: staffService}
}

// @Summary Register customer
// @Description Staff register customer
// @Tags Staff
// @Accept json
// @Param request body model.RegisterRequest true "Auth payload"
// @Produce  json
// @Router /staff/register-customer [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *StaffHandler) RegisterCustomer(ctx *gin.Context) {
	var registerRequest model.RegisterRequest

	if err := validation.BindJsonAndValidate(ctx, &registerRequest); err != nil {
		return
	}

	err := handler.staffService.RegisterCustomer(ctx, registerRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Add amount to account
// @Description Add amount to account
// @Tags Staff
// @Accept json
// @Param request body model.AddAmountToAccountRequest true "AddAmount payload"
// @Produce  json
// @Router /staff/add-amount [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *StaffHandler) AddAmountToAccount(ctx *gin.Context) {
	var request model.AddAmountToAccountRequest

	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	err := handler.staffService.AddAmountToAccount(ctx, &request)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Get transactions by account number
// @Description Get transactions by account number
// @Tags Staff
// @Accept json
// @Param accountNumber query string true "Account payload"
// @Produce  json
// @Router /staff/transactions-by-account [get]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[[]model.GetTransactionsByCustomerResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *StaffHandler) GetTransactionsByAccountNumber(ctx *gin.Context) {
	accountNumber := ctx.Query("accountNumber")
	transactions, err := handler.staffService.GetTransactionsByAccountNumber(ctx, accountNumber)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InvalidRequest,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&transactions))
}
