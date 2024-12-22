package serviceimplement

import (
	"errors"
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
	"math/rand"
)

type AccountService struct {
	accountRepository  repository.AccountRepository
	customerRepository repository.CustomerRepository
	CoreService        service.CoreService
}

func NewAccountService(accountRepo repository.AccountRepository, customerRepo repository.CustomerRepository, coreService service.CoreService) service.AccountService {
	return &AccountService{accountRepository: accountRepo, customerRepository: customerRepo, CoreService: coreService}
}

func (service *AccountService) generateAccountNumber(ctx *gin.Context) (string, error) {
	for {
		accountNumber := fmt.Sprintf("%012d", rand.Int63n(1000000000000)) // 12-digit number with leading zeros

		// Check for unique
		existsAccount, err := service.accountRepository.GetOneByNumberQuery(ctx, accountNumber)
		if err != nil {
			return "", err
		}
		if existsAccount == nil {
			return accountNumber, nil
		}
	}
}

func (service *AccountService) AddNewAccount(ctx *gin.Context, customerId int64) error {

	newNumber, err := service.generateAccountNumber(ctx)
	if err != nil {
		return err
	}
	err = service.accountRepository.CreateCommand(ctx, &entity.Account{
		CustomerID: customerId,
		Number:     newNumber,
		Balance:    0,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *AccountService) InternalTransfer(ctx *gin.Context, transferReq model.InternalTransferRequest) error {
	//get customer and check info
	customerId, exists := ctx.Get("customerId")
	if !exists {
		return errors.New("customer not exists")
	}

	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId.(int64))
	if err != nil {
		return err
	}
	if sourceCustomer == nil {
		return errors.New("customer not exists")
	}

	//get account by customerId and check sourceNumber
	sourceAccount, err := service.accountRepository.GetOneByCustomerIdQuery(ctx, customerId.(int64))
	if err != nil {
		return err
	}
	if sourceAccount == nil {
		return errors.New("source account not exists")
	}
	if sourceAccount.Number != transferReq.SourceAccountNumber {
		return errors.New("source account not match")
	}

	//check targetNumber
	targetAccount, err := service.accountRepository.GetOneByNumberQuery(ctx, transferReq.TargetAccountNumber)
	if err != nil {
		return err
	}
	if targetAccount == nil {
		return errors.New("target account not exists")
	}
	//estimate fee
	fee, err := service.CoreService.EstimateTransferFee(ctx, transferReq.Amount)
	if err != nil {
		return err
	}

	//check is source fee and change balance
	checkFee := *transferReq.IsSourceFee
	if checkFee {
		totalDeduction := transferReq.Amount + fee
		if sourceAccount.Balance < totalDeduction {
			return errors.New("insufficient balance in source account")
		}
		sourceAccount.Balance -= totalDeduction
		targetAccount.Balance += transferReq.Amount
	} else {
		if sourceAccount.Balance < transferReq.Amount {
			return errors.New("insufficient balance in source account")
		}
		sourceAccount.Balance -= transferReq.Amount
		targetAccount.Balance += transferReq.Amount - fee
	}

	//update to DB
	err = service.accountRepository.UpdateCommand(ctx, *sourceAccount)
	if err != nil {
		return err
	}
	err = service.accountRepository.UpdateCommand(ctx, *targetAccount)
	if err != nil {
		return err
	}

	//Save to history, notify, response history

	return nil
}

func (service *AccountService) GetCustomerByAccountNumber(ctx *gin.Context, accountNumber string) (*entity.User, error) {
	customer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, accountNumber)
	if err != nil {
		return nil, err
	}
	return customer, nil
}
