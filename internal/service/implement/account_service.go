package serviceimplement

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/HMAC_signature"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/generate_number_code"
	"github.com/gin-gonic/gin"
)

var (
	privateKeyRsaTeam, _ = env.GetEnv("SECRET_KEY_OF_RSA_TEAM")
	bankIdInRsaTeam, _   = env.GetEnv("BANK_ID_IN_RSA_TEAM")
)

type Response struct {
	Success bool     `json:"success"`
	Message []string `json:"message"`
	Data    struct {
		CustomerName string  `json:"customerName"`
		Balance      float64 `json:"balance"`
	} `json:"data"`
}

type AccountService struct {
	accountRepository  repository.AccountRepository
	customerRepository repository.CustomerRepository
	partnerBankService service.PartnerBankService
}

func NewAccountService(accountRepo repository.AccountRepository,
	customerRepo repository.CustomerRepository,
	partnerBankService service.PartnerBankService) service.AccountService {
	return &AccountService{accountRepository: accountRepo,
		customerRepository: customerRepo,
		partnerBankService: partnerBankService}
}

func (service *AccountService) AddNewAccount(ctx *gin.Context, customerId int64) error {

	newNumber := generate_number_code.GenerateRandomNumber(12)
	var balance int64 = 0
	err := service.accountRepository.CreateCommand(ctx, &entity.Account{
		CustomerID: customerId,
		Number:     newNumber,
		Balance:    &balance,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) GetCustomerByAccountNumber(ctx *gin.Context, accountNumber string) (*entity.User, error) {
	customer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (service *AccountService) UpdateBalanceByAccountNumber(ctx *gin.Context, amount int64, number string) (int64, error) {
	_, err := service.accountRepository.GetOneByNumberQuery(ctx, number)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return 0, errors.New("account not found")
		}
		return 0, err
	}
	newBalance, err := service.accountRepository.UpdateBalanceCommand(ctx, number, amount)
	if err != nil {
		return 0, err
	}
	return newBalance, nil
}

func (service *AccountService) GetAccountByCustomerId(ctx *gin.Context, customerId int64) (*entity.Account, error) {
	account, err := service.accountRepository.GetOneByCustomerIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

func (service *AccountService) GetAccountByNumber(ctx *gin.Context, number string) (*entity.Account, error) {
	account, err := service.accountRepository.GetOneByNumberQuery(ctx, number)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("account not found")
		}
		return nil, err
	}
	return account, nil
}

func (service *AccountService) GetExternalAccountName(ctx *gin.Context, detail model.GetExternalAccountNameRequest) (string, error) {
	//check exists external bank
	partnerBank, err := service.partnerBankService.GetBankById(ctx, detail.BankId)
	if err != nil {
		return "", err
	}
	bankIdInRsaTeamInt, err := strconv.ParseInt(bankIdInRsaTeam, 10, 64)
	if err != nil {
		return "", err
	}
	//setup payload
	req := &model.SearchExternalAccountRequest{
		BankId:        bankIdInRsaTeamInt,
		TimeStamp:     time.Now().Unix(),
		AccountNumber: detail.AccountNumber,
	}
	reqBytes, err := json.Marshal(req)
	//hash data
	hashData := HMAC_signature.GenerateHMAC(string(reqBytes), privateKeyRsaTeam)
	//setup and call to partner bank server
	request, err := http.NewRequest("POST", partnerBank.ResearchApi, bytes.NewBuffer(reqBytes))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("HMAC", hashData)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	//handler response
	var resp Response
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", err
	}
	if response.StatusCode == http.StatusOK {
		return resp.Data.CustomerName, nil
	} else {
		return "", errors.New(resp.Message[0])
	}
}
