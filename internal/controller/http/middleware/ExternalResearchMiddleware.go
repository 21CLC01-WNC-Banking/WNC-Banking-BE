package middleware

import (
	"encoding/json"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/HMAC_signature"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var secretString = os.Getenv("SECRET_KEY_FOR_EXTERNAL_BANK")

type ExternalSearchMiddleware struct {
	partnerBankService service.PartnerBankService
}

func NewExternalSearchMiddleware(partnerBankService service.PartnerBankService) *ExternalSearchMiddleware {
	return &ExternalSearchMiddleware{
		partnerBankService: partnerBankService,
	}
}

func (middleware *ExternalSearchMiddleware) VerifyPartnerBank(c *gin.Context, bankCode string) (*entity.PartnerBank, error) {
	bank, err := middleware.partnerBankService.GetPartnerBankByBankCode(c, bankCode)
	if err != nil {
		return nil, err
	}
	return bank, nil
}

func (middleware *ExternalSearchMiddleware) VerifyAPI(c *gin.Context) {
	//get body
	var req model.AccountNumberInfoRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "body",
			}))
		return
	}
	//check exists bank
	_, err = middleware.VerifyPartnerBank(c, req.SrcBankCode)
	if err != nil {
		if err.Error() == "partner bank not found" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: "ngân hàng chưa được đăng ký",
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
					Field:   "srcBankCode",
				}))
			return
		}
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
				Field:   "srcBankCode",
			}))
		return
	}
	//check time
	expTime, _ := time.Parse(time.RFC3339Nano, req.Exp)
	if expTime.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin cũ đã quá hạn",
				Code:    httpcommon.ErrorResponseCode.TimeoutRequest,
				Field:   "Exp",
			}))
		return
	}

	//get hashed data in header and check valid
	hashedData := c.GetHeader("hashedData")
	if hashedData == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "hashedData",
			}))
		return
	}
	data, err := json.Marshal(req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
				Field:   "encoding body",
			}))
		return
	}
	valid := HMAC_signature.VerifyHMAC(string(data), hashedData, secretString)
	if !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "verify hashed data",
			}))
		return
	}
	c.Set("request", req)
	c.Next()
	return
}
