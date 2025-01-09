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
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/HMAC_signature"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

var (
	rsa_rpivate_key, _ = env.GetEnv("RSA_PRIVATE_KEY")
)

type RSAMiddleware struct {
	externalSearchMiddleware *ExternalSearchMiddleware
}

func NewRSAMiddleware(externalSearchMiddleware *ExternalSearchMiddleware) *RSAMiddleware {
	return &RSAMiddleware{externalSearchMiddleware: externalSearchMiddleware}
}

func (middleware *RSAMiddleware) loadRSAPrivateKey() (*rsa.PrivateKey, error) {
	rsaKey, err := base64.StdEncoding.DecodeString(rsa_rpivate_key)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(rsaKey)
	key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse private key")
	}
	return key, nil
}

func loadRSAPublicKey(pemKey string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pemKey))
	if block == nil {
		return nil, errors.New("failed to decode PEM block, invalid format")
	}
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

func (middleware *RSAMiddleware) VerifyRSASignature(partnerBank entity.PartnerBank, data []byte, signedData string) error {
	//load public key
	rsaPublicKey, err := loadRSAPublicKey(partnerBank.PublicKey)
	if err != nil {
		return err
	}
	//decode signature
	signatureBytes, err := base64.StdEncoding.DecodeString(signedData)
	if err != nil {
		return errors.New("failed to decode signature: %v")
	}
	//create hash from data
	hash := sha256.Sum256(data)
	//verify signature
	return rsa.VerifyPKCS1v15(rsaPublicKey, crypto.SHA256, hash[:], signatureBytes)
}

func (middleware *RSAMiddleware) Verify(c *gin.Context) {
	//get body
	var req model.ExternalTransferRequest
	err := validation.BindJsonAndValidate(c, &req)
	if err != nil {
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
	//check time
	expTime, err := time.Parse(time.RFC3339Nano, req.Exp)
	if expTime.Before(time.Now()) {
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
	payloadRequest := model.ExternalPayload{
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
	err = middleware.VerifyRSASignature(*partnerBank, data, req.SignedData)
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
	c.Set("secureType", "RSA")
	c.Next()
	return
}
