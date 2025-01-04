package v1

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type PartnerBankHandler struct {
	accountService     service.AccountService
	transactionService service.TransactionService
	partnerBankService service.PartnerBankService
	middlewareRSA      *middleware.RSAMiddleware
}

func NewPartnerBankHandler(
	accountService service.AccountService,
	transactionService service.TransactionService,
	partnerBankService service.PartnerBankService,
	middlewareRSA *middleware.RSAMiddleware) *PartnerBankHandler {
	return &PartnerBankHandler{
		accountService:     accountService,
		transactionService: transactionService,
		partnerBankService: partnerBankService,
		middlewareRSA:      middlewareRSA}
}

// @Summary Get account name
// @Description Get account name in our bank by account number from external bank
// @Tags Partner bank
// @Accept json
// @Param request body model.AccountNumberInfoRequest true "PartnerBank payload"
// @Produce  json
// @Router /partner-bank/get-account-information [post]
// @Success 200 {object} httpcommon.HttpResponse[model.AccountNumberInfoResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *PartnerBankHandler) GetAccountNumberInfo(c *gin.Context) {
	req, _ := c.Get("request")
	data := req.(model.AccountNumberInfoRequest)
	user, err := handler.accountService.GetCustomerByAccountNumber(c, data.DesAccountNumber)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			c.AbortWithStatusJSON(http.StatusNotFound, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: "không tìm thấy thông tin tài khoản",
					Code:    httpcommon.ErrorResponseCode.RecordNotFound,
					Field:   "account number",
				}))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.AccountNumberInfoResponse](&model.AccountNumberInfoResponse{
		DesAccountNumber: data.DesAccountNumber,
		DesAccountName:   user.Name,
	}))
}

// @Summary Partner bank
// @Description get list partner banks
// @Tags Partner bank
// @Accept json
// @Produce  json
// @Router /partner-bank/ [GET]
// @Success 200 {object} httpcommon.HttpResponse[[]entity.PartnerBank]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *PartnerBankHandler) GetListPartnerBank(c *gin.Context) {
	listBank, err := handler.partnerBankService.GetListPartnerBank(c)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			}))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse[[]entity.PartnerBank](&listBank))
}

// @Summary Partner bank
// @Description receive external transfer from partner banks
// @Tags Partner bank
// @Accept json
// @Produce  json
// @Router /partner-bank/external-transfer-rsa [POST]
// @Success 200 {object} httpcommon.HttpResponse[model.ExternalTransactionResponse]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *PartnerBankHandler) ReceiveExternalTransfer(c *gin.Context) {
	req, _ := c.Get("request")
	externalTransaction := req.(model.ExternalTransactionRequest)

	partnerBankId, _ := c.Get("partnerBankId")
	//save
	err := handler.transactionService.ReceiveExternalTransfer(c, externalTransaction, partnerBankId.(int64))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
				Field:   "",
			}))
		return
	}
	//encode response
	responseString := "transfer success"
	signedResponse, err := handler.middlewareRSA.SignDataRSA(responseString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
				Field:   "",
			}))
		return
	}
	c.JSON(http.StatusOK, httpcommon.NewSuccessResponse[model.ExternalTransactionResponse](&model.ExternalTransactionResponse{
		Data: responseString, SignedData: signedResponse}))
}
