package middleware

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/jwt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
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
	//get hashed data in header
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

	//decode token
	claim, err := jwt.VerifyToken(hashedData, secretString)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.TokenExpired {
			c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: "thông tin cũ đã quá hạn",
					Code:    httpcommon.ErrorResponseCode.TimeoutRequest,
					Field:   "hashedData",
				}))
			return
		}
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "hashedData",
			}))
		return
	}
	//getPayloadinfo
	payload, ok := claim.Payload.(map[string]interface{})
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "payload",
			}))
		return
	}
	partnerBankCode := payload["srcBankCode"].(string)
	_, err = middleware.VerifyPartnerBank(c, partnerBankCode)
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
	}

	//rehash token
	var req model.AccountNumberInfoRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "body",
			}))
		return
	}
	claim.Payload = req
	expectedHash, err := jwt.GenerateTokenByClaims(*claim, secretString)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
				Field:   "rehash",
			}))
		return
	}
	//check token
	if expectedHash != hashedData {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "map token",
			}))
		return
	}
	c.Set("request", req)
	c.Next()
	return
}
