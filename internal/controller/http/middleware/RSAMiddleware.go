package middleware

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/jwt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strings"
)

var rsaSecretKey = os.Getenv("RSA_SECRET_KEY")

type RSAMiddleware struct {
	externalSearchMiddleware *ExternalSearchMiddleware
	accountService           service.AccountService
}

func NewRSAMiddleware(externalSearchMiddleware *ExternalSearchMiddleware, accountService service.AccountService) *RSAMiddleware {
	return &RSAMiddleware{externalSearchMiddleware: externalSearchMiddleware, accountService: accountService}
}

func loadPrivateKey() (*rsa.PrivateKey, error) {
	rsaSecretKey = strings.ReplaceAll(rsaSecretKey, `\n`, "\n")
	block, _ := pem.Decode([]byte(rsaSecretKey))
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse private key")
	}
	return key, nil
}

func SignWithRSAPrivateKey(data string) (string, error) {
	privateKey, err := loadPrivateKey()
	if err != nil {
		return "", err
	}
	// Hash the data using SHA-256
	hasher := sha256.New()
	_, err = hasher.Write([]byte(data))
	if err != nil {
		return "", errors.New("failed to hash data")
	}
	hashed := hasher.Sum(nil)

	// Sign the hashed data using RSA and PKCS1v15
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed)
	if err != nil {
		return "", errors.New("failed to sign data")
	}

	// Return the signature as a base64 encoded string
	return base64.StdEncoding.EncodeToString(signature), nil
}

func loadPublicKey(pemKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse public key")
	}
	rsaKey, ok := key.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("not an RSA public key")
	}
	return rsaKey, nil
}

func (middleware *RSAMiddleware) VerifyAccountNumber(c *gin.Context, accountNumber string) error {
	_, err := middleware.accountService.GetAccountByNumber(c, accountNumber)
	if err != nil {
		return err
	}
	return nil
}

func (middleware *RSAMiddleware) VerifySignature(partnerBank entity.PartnerBank, data model.ExternalTransactionRequest, signedData string) error {
	//load public key
	rsaPublicKey, err := loadPublicKey(partnerBank.PublicKey)
	if err != nil {
		return err
	}
	//change struct to json
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	//decode signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signedData)
	if err != nil {
		return fmt.Errorf("failed to decode signature: %v", err)
	}
	//create hash from data
	hash := sha256.Sum256(dataBytes)
	//verify signature
	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signatureBytes)
}

func (middleware *RSAMiddleware) Verify(c *gin.Context) {
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
	partnerBank, err := middleware.externalSearchMiddleware.VerifyPartnerBank(c, partnerBankCode)
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
	//check account number
	desAccountNumber := payload["desAccountNumber"].(string)
	err = middleware.VerifyAccountNumber(c, desAccountNumber)
	if err != nil {
		if err.Error() == "account not found" {
			c.AbortWithStatusJSON(http.StatusNotFound, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: "không tìm thấy thông tin tài khoản",
					Code:    httpcommon.ErrorResponseCode.RecordNotFound,
					Field:   "desAccountNumber",
				}))
			return
		}
	}

	//rehash token
	var req model.ExternalTransactionData
	//err = c.ShouldBindJSON(&req)
	err = validation.BindJsonAndValidate(c, &req)
	if err != nil {
		fmt.Print(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "body",
			}))
		return
	}
	transactionRequest := model.ExternalTransactionRequest{
		SrcAccountNumber: req.SrcAccountNumber,
		SrcBankCode:      req.SrcBankCode,
		DesAccountNumber: req.DesAccountNumber,
		Amount:           req.Amount,
		Description:      req.Description,
		IsSourceFee:      req.IsSourceFee,
	}
	claim.Payload = transactionRequest
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
	//check signature
	err = middleware.VerifySignature(*partnerBank, transactionRequest, req.SignedData)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin gói không chính xác hoặc bị chỉnh sửa",
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
				Field:   "signature",
			}))
		return
	}
	c.Set("partnerBankId", partnerBank.ID)
	c.Set("request", transactionRequest)
	c.Set("middlewareType", "RSA")
	c.Next()
	return
}
