package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	beanimplement "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean/implement"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/HMAC_signature"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type PGPMiddleware struct {
	externalSearchMiddleware *ExternalSearchMiddleware
	keyLoader                *beanimplement.KeyLoader
}

func NewPGPMiddleware(externalSearchMiddleware *ExternalSearchMiddleware,
	keyLoader *beanimplement.KeyLoader) *PGPMiddleware {
	return &PGPMiddleware{
		externalSearchMiddleware: externalSearchMiddleware,
		keyLoader:                keyLoader,
	}
}

func (middleware *PGPMiddleware) loadPGPPrivateKey() (*crypto.Key, error) {
	privateKey, err := crypto.NewKeyFromArmored(string(middleware.keyLoader.PGPKey))
	if err != nil {
		return nil, errors.New("failed to parse private key")
	}
	return privateKey, nil
}

func loadPGPPublicKey(pemKey string) (*crypto.Key, error) {
	publicKey, err := crypto.NewKeyFromArmored(pemKey)
	if err != nil {
		return nil, errors.New("failed to parse public key")
	}
	return publicKey, nil
}

func (middleware *PGPMiddleware) SignDataPGP(data string) (string, error) {
	keyRing, err := middleware.loadPGPPrivateKey()
	if err != nil {
		return "", err
	}
	signer, err := crypto.PGP().Sign().SigningKey(keyRing).Detached().New()
	if err != nil {
		return "", err
	}
	signature, err := signer.Sign([]byte(data), crypto.Armor)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(signature), nil
}

func (middleware *PGPMiddleware) VerifyPGPSignature(partnerBank entity.PartnerBank, data model.ExternalTransactionRequest, signedData string) error {
	publicKey, err := loadPGPPublicKey(partnerBank.PublicKey)
	if err != nil {
		return err
	}

	signatureBytes, err := base64.StdEncoding.DecodeString(signedData)
	if err != nil {
		return errors.New("failed to decode signature")
	}

	verifier, err := crypto.PGP().Verify().VerificationKey(publicKey).New()
	if err != nil {
		return err
	}

	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	verifyResult, err := verifier.VerifyDetached(dataBytes, signatureBytes, crypto.Armor)
	if err != nil {
		return errors.New("signature verification failed")
	}
	if sigErr := verifyResult.SignatureError(); sigErr != nil {
		return sigErr
	}
	return nil
}

func (middleware *PGPMiddleware) Verify(c *gin.Context) {
	//get body
	var req model.ExternalTransactionData
	err := validation.BindJsonAndValidate(c, &req)
	if err != nil {
		return
	}
	//check exists bank
	partnerBank, err := middleware.externalSearchMiddleware.VerifyPartnerBank(c, req.SrcBankCode)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "Bank not registered",
				Code:    httpcommon.ErrorResponseCode.Unauthorized,
				Field:   "srcBankCode",
			}))
		return
	}
	//check time
	if req.Exp.Before(time.Now()) {
		c.AbortWithStatusJSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: "Expired information",
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
				Message: "Invalid package information",
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
	err = middleware.VerifyPGPSignature(*partnerBank, payloadRequest, req.SignedData)
	if err != nil {
		fmt.Println(err.Error())
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
	c.Set("secureType", "PGP")
	c.Next()
	return
}
