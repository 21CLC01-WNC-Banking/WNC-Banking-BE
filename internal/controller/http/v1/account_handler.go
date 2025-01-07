package v1

import (
	"net/http"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"

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

func NewAccountHandler(accountService service.AccountService, savedReceiverService service.SavedReceiverService, authService service.AuthService) *AccountHandler {
	return &AccountHandler{
		accountService:       accountService,
		savedReceiverService: savedReceiverService,
		authService:          authService,
	}
}

// @Summary Get Customer Name by Account Number
// @Description Get Customer Name by Account Number
// @Tags Accounts
// @Param accountNumber query string true "Account payload"
// @Produce  json
// @Router /account/customer-name [get]
// @Param  Authorization header string true "Authorization: Bearer"
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

// @Summary Get Account by Customer Id
// @Description Get Account by Customer Id
// @Tags Accounts
// @Produce  json
// @Router /account/ [get]
// @Param  Authorization header string true "Authorization: Bearer"
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
		return
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

// @Summary Get external account name
// @Description Get external account name by account number from external bank
// @Tags Account
// @Accept json
// @Param request body model.GetExternalAccountNameRequest true "PartnerBank payload"
// @Produce  json
// @Router /account/get-external-account-name [post]
// @Param  Authorization header string true "Authorization: Bearer"
// @Success 200 {object} httpcommon.HttpResponse[string]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AccountHandler) GetExternalAccountName(c *gin.Context) {
	var req model.GetExternalAccountNameRequest
	err := validation.BindJsonAndValidate(c, &req)
	if err != nil {
		return
	}
	if req.AccountNumber != "987654321098" || req.BankId != 1 {
		c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	accountName := "NGUYEN NHAT NAM"
	//accountName, err := handler.accountService.GetExternalAccountName(c, req)
	//if err != nil {
	//	c.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
	//		Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
	//	}))
	//	return
	//}
	//c.JSON(http.StatusOK, httpcommon.NewSuccessResponse[string](&accountName))
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse[string](&accountName))
}
