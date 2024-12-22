package serviceimplement

import (
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/generate_number_code"
	"github.com/gin-gonic/gin"
)

type AccountService struct {
	accountRepository  repository.AccountRepository
	customerRepository repository.CustomerRepository
}

func NewAccountService(accountRepo repository.AccountRepository, customerRepo repository.CustomerRepository) service.AccountService {
	return &AccountService{accountRepository: accountRepo, customerRepository: customerRepo}
}

func (service *AccountService) AddNewAccount(ctx *gin.Context, customerId int64) error {

	newNumber := generate_number_code.GenerateRandomNumber(12)

	err := service.accountRepository.CreateCommand(ctx, &entity.Account{
		CustomerID: customerId,
		Number:     newNumber,
		Balance:    0,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) GetCustomerByAccountNumber(ctx *gin.Context, accountNumber string) (*entity.Customer, error) {
	customer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	return customer, nil
}

func (service *AccountService) UpdateBalanceByAccountNumber(ctx *gin.Context, balance int64, number string) error {
	account, err := service.accountRepository.GetOneByNumberQuery(ctx, number)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("account not found")
		}
		return err
	}
	account.Balance = balance
	err = service.accountRepository.UpdateCommand(ctx, *account)
	if err != nil {
		return err
	}
	return nil
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
