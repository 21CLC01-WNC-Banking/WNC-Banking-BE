package v1

import (
	"net/http"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService       service.AccountService
	savedReceiverService service.SavedReceiverService
	authService          service.AuthService
}

func NewAccountHandler(accountService service.AccountService, savedReceiverService service.SavedReceiverService) *AccountHandler {
	return &AccountHandler{
		accountService:       accountService,
		savedReceiverService: savedReceiverService,
	}
}

// @Summary Get Customer Name by Account Number
// @Description Get Customer Name by Account Number
// @Tags Accounts
// @Param accountNumber query string true "Account payload"
// @Produce  json
// @Router /account/customer-name [get]
// @Success 200 {object} httpcommon.HttpResponse[model.GetCustomerNameByAccountNumberResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AccountHandler) GetCustomerNameByAccountNumber(ctx *gin.Context) {
	accountNumber := ctx.Query("accountNumber")
	customer, err := handler.accountService.GetCustomerByAccountNumber(ctx, accountNumber)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse(&model.GetCustomerNameByAccountNumberResponse{
		Name: customer.Name,
	}))
}

// @Summary Get Account by Customer ID
// @Description Get Account by Customer ID
// @Tags Accounts
// @Produce  json
// @Router /account/ [get]
// @Success 200 {object} httpcommon.HttpResponse[model.AccountResponse]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AccountHandler) GetAccountByCustomerId(ctx *gin.Context) {
	//get customer and check info
	customerId, exists := ctx.Get("userId")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: "Can not find CustomerId", Field: "token", Code: httpcommon.ErrorResponseCode.InvalidUserInfo,
		}))
		return
	}
	account, err := handler.accountService.GetAccountByCustomerId(ctx, customerId.(int64))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
	}
	customer, err := handler.authService.GetUserById(ctx, customerId.(int64))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InvalidUserInfo,
		}))
		return
	}
	ctx.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.AccountResponse](&model.AccountResponse{Account: account, Name: customer.Name}))
}
