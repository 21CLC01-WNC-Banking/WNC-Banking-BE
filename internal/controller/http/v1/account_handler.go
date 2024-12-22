package v1

import (
	"net/http"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	accountService       service.AccountService
	savedReceiverService service.SavedReceiverService
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

// @Summary Add Internal Receiver
// @Description Add a new internal receiver
// @Tags Receivers
// @Accept  json
// @Produce  json
// @Param receiver body model.InternalReceiver true "Internal Receiver Payload"
// @Success 204 "No Content"
// @Failure 500 {object} httpcommon.HttpResponse[any] "Internal Server Error"
// @Router /account/add-internal-receiver [post]
func (handler *AccountHandler) AddInternalReceiver(ctx *gin.Context) {
	var receiver model.InternalReceiver

	if err := validation.BindJsonAndValidate(ctx, &receiver); err != nil {
		return
	}
	err := handler.savedReceiverService.AddInternalReceiver(ctx, receiver)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(httpcommon.Error{
			Message: err.Error(), Field: "", Code: httpcommon.ErrorResponseCode.InternalServerError,
		}))
		return
	}
	ctx.AbortWithStatus(204)
}
