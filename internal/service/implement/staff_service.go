package serviceimplement

import (
	"errors"
	"time"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/rand"
)

type StaffService struct {
	customerRepository    repository.CustomerRepository
	passwordEncoder       bean.PasswordEncoder
	accountService        service.AccountService
	accountRepository     repository.AccountRepository
	transactionRepository repository.TransactionRepository
	mailCLient            bean.MailClient
}

func NewStaffService(
	customerRepository repository.CustomerRepository,
	passwordEncoder bean.PasswordEncoder,
	accountService service.AccountService,
	accountRepository repository.AccountRepository,
	transactionRepository repository.TransactionRepository,
	mailCLient bean.MailClient,
) service.StaffService {
	return &StaffService{
		customerRepository:    customerRepository,
		passwordEncoder:       passwordEncoder,
		accountService:        accountService,
		accountRepository:     accountRepository,
		transactionRepository: transactionRepository,
		mailCLient:            mailCLient,
	}
}

func generateRandomPassword(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"

	rand.Seed(uint64(time.Now().UnixNano()))

	password := make([]byte, length)
	for i := range password {
		password[i] = charset[rand.Intn(len(charset))]
	}

	return string(password)
}

func (service *StaffService) RegisterCustomer(ctx *gin.Context, registerRequest model.RegisterRequest) error {
	existsCustomer, err := service.customerRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil && err.Error() != httpcommon.ErrorMessage.SqlxNoRow {
		return err
	}
	if existsCustomer != nil {
		return errors.New("Email have already registered")
	}

	// generate random password
	randomPassword := generateRandomPassword(8)
	hashPW, err := service.passwordEncoder.Encrypt(randomPassword)
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

	// send random password to email
	emailBody := service.mailCLient.GenerateRandomPasswordBody(registerRequest.Email, randomPassword)
	err = service.mailCLient.SendEmail(ctx, registerRequest.Email, "Generate random password", emailBody)
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
	if request.Description == "" {
		request.Description = "staff add amount to account"
	}

	isSourceFee := false
	transaction := entity.Transaction{
		TargetAccountNumber: request.AccountNumber,
		TargetBalance:       balance,
		Amount:              request.Amount,
		Type:                "internal",
		Status:              "success",
		Description:         request.Description,
		IsSourceFee:         &isSourceFee,
	}

	_, err = service.transactionRepository.CreateCommand(ctx, &transaction)
	return err
}

func (service *StaffService) GetTransactionsByAccountNumber(ctx *gin.Context, accountNumber string) (*model.GetTransactionsByCustomerResponse, error) {
	account, err := service.accountRepository.GetOneByNumberQuery(ctx, accountNumber)
	if err != nil || account == nil {
		return nil, errors.New("account not found")
	}

	transactions, err := service.transactionRepository.GetTransactionByAccountNumber(ctx, accountNumber)
	if err != nil {
		return nil, err
	}

	userId, exists := ctx.Get("userId")
	if !exists {
		return nil, errors.New("customer not exists")
	}
	customer, err := service.customerRepository.GetOneByIdQuery(ctx, userId.(int64))
	if err != nil {
		return nil, err
	}

	transactionResp := make([]model.GetTransactionsResponse, 0)

	for _, transaction := range transactions {
		var amount int64
		var balance int64
		if transaction.TargetAccountNumber == accountNumber {
			amount = transaction.Amount
			balance = transaction.TargetBalance
		} else {
			amount = transaction.Amount * -1
			balance = transaction.SourceBalance
		}
		transactionResp = append(transactionResp, model.GetTransactionsResponse{
			Id:                  transaction.Id,
			Amount:              amount,
			CreatedAt:           transaction.CreatedAt,
			Description:         transaction.Description,
			Type:                transaction.Type,
			Balance:             balance,
			SourceAccountNumber: transaction.SourceAccountNumber,
			TargetAccountNumber: transaction.TargetAccountNumber,
		})
	}

	return &model.GetTransactionsByCustomerResponse{
		CustomerName: customer.Name,
		Transactions: transactionResp,
	}, nil
}
