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
	beanimplement "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean/implement"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/HMAC_signature"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"
)

var rsaSecretKey = os.Getenv("RSA_SECRET_KEY")

type RSAMiddleware struct {
	externalSearchMiddleware *ExternalSearchMiddleware
	accountService           service.AccountService
	keyLoader                *beanimplement.KeyLoader
}

func NewRSAMiddleware(externalSearchMiddleware *ExternalSearchMiddleware,
	accountService service.AccountService,
	keyLoader *beanimplement.KeyLoader) *RSAMiddleware {
	return &RSAMiddleware{externalSearchMiddleware: externalSearchMiddleware,
		accountService: accountService,
		keyLoader:      keyLoader}
}

func (middleware *RSAMiddleware) loadRSAPrivateKey() (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(middleware.keyLoader.RSAKey)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse private key")
	}
	return key, nil
}

func loadRSAPublicKey(pemKey string) (*rsa.PublicKey, error) {
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

func (middleware *RSAMiddleware) SignDataRSA(data string) (string, error) {
	// Hash the data using SHA256
	hashed := sha256.Sum256([]byte(data))
	//get private key
	privateKey, err := middleware.loadRSAPrivateKey()
	if err != nil {
		return "", err
	}
	// Sign the hashed data with RSA
	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, hashed[:])
	if err != nil {
		return "", err
	}
	// Encode the signature to Base64
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (middleware *RSAMiddleware) VerifySignature(partnerBank entity.PartnerBank, data model.ExternalTransactionRequest, signedData string) error {
	//load public key
	rsaPublicKey, err := loadRSAPublicKey(partnerBank.PublicKey)
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
		return errors.New("failed to decode signature: %v")
	}
	//create hash from data
	hash := sha256.Sum256(dataBytes)
	//verify signature
	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signatureBytes)
}

func (middleware *RSAMiddleware) Verify(c *gin.Context) {
	//get body
	var req model.ExternalTransactionData
	err := validation.BindJsonAndValidate(c, &req)
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
	partnerBank, err := middleware.externalSearchMiddleware.VerifyPartnerBank(c, req.SrcBankCode)
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
	//check source account number
	/* call service to check name */
	//check target account number
	err = middleware.VerifyAccountNumber(c, req.DesAccountNumber)
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
	//check amount
	if req.Amount < 1 {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "số tiền chuyển không hợp lệ",
				Code:    httpcommon.ErrorResponseCode.TimeoutRequest,
				Field:   "Amount",
			}))
		return
	}
	//check time
	if req.Exp.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "thông tin cũ đã quá hạn",
				Code:    httpcommon.ErrorResponseCode.TimeoutRequest,
				Field:   "Exp",
			}))
		return
	}

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
	//check hashed data
	payloadRequest := model.ExternalTransactionRequest{
		SrcAccountNumber: req.SrcAccountNumber,
		SrcBankCode:      req.SrcBankCode,
		DesAccountNumber: req.DesAccountNumber,
		Amount:           req.Amount,
		Description:      req.Description,
		IsSourceFee:      req.IsSourceFee,
		Exp:              req.Exp,
	}
	data, err := json.Marshal(payloadRequest)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
				Field:   "encoding payload",
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

	//check signature
	err = middleware.VerifySignature(*partnerBank, payloadRequest, req.SignedData)
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
	c.Set("request", payloadRequest)
	c.Next()
	return
}
