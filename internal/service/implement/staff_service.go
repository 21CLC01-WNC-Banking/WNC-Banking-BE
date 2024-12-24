package serviceimplement

import (
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type StaffService struct {
	customerRepository    repository.CustomerRepository
	passwordEncoder       bean.PasswordEncoder
	accountService        service.AccountService
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
}

func NewStaffService(
	customerRepository repository.CustomerRepository,
	passwordEncoder bean.PasswordEncoder,
	accountService service.AccountService,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
) service.StaffService {
	return &StaffService{
		customerRepository:    customerRepository,
		passwordEncoder:       passwordEncoder,
		accountService:        accountService,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
	}
}

func (service *StaffService) RegisterCustomer(ctx *gin.Context, registerRequest model.RegisterRequest) error {
	existsCustomer, err := service.customerRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil && err.Error() != httpcommon.ErrorMessage.SqlxNoRow {
		return err
	}
	if existsCustomer != nil {
		return errors.New("Email have already registered")
	}
	hashPW, err := service.passwordEncoder.Encrypt(registerRequest.Password)
	if err != nil {
		return err
	}
	newCustomer := &entity.User{
		Email:       registerRequest.Email,
		Name:        registerRequest.Name,
		PhoneNumber: registerRequest.PhoneNumber,
		Password:    string(hashPW),
		RoleId:      1,
	}
	err = service.customerRepository.CreateCommand(ctx, newCustomer)
	if err != nil {
		return err
	}

	// auto create an account
	currentCustomer, err := service.customerRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil {
		return err
	}
	err = service.accountService.AddNewAccount(ctx, currentCustomer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (service *StaffService) AddAmountToAccount(ctx *gin.Context, request *model.AddAmountToAccountRequest) error {
	balance, err := service.accountRepository.UpdateBalanceCommand(ctx, request.AccountNumber, request.Amount)
	if err != nil {
		return err
	}

	isSourceFee := false
	transaction := entity.Transaction{
		TargetAccountNumber: request.AccountNumber,
		TargetBalance:       balance,
		Type:                "internal",
		Status:              "success",
		Description:         "staff add amount to account",
		IsSourceFee:         &isSourceFee,
	}

	_, err = service.transactionRepository.CreateCommand(ctx, &transaction)
	return err
}
