package serviceimplement

import (
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"strconv"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/constants"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/mail"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/redis"
	"github.com/gin-gonic/gin"
)

type TransactionService struct {
	transactionRepository repository.TransactionRepository
	customerRepository    repository.CustomerRepository
	accountService        service.AccountService
	coreService           service.CoreService
	redisClient           bean.RedisClient
	mailClient            bean.MailClient
}

func NewTransactionService(transactionRepository repository.TransactionRepository,
	customerRepository repository.CustomerRepository,
	accountService service.AccountService,
	coreService service.CoreService,
	redisClient bean.RedisClient,
	mailClient bean.MailClient) service.TransactionService {
	return &TransactionService{
		transactionRepository: transactionRepository,
		customerRepository:    customerRepository,
		accountService:        accountService,
		coreService:           coreService,
		redisClient:           redisClient,
		mailClient:            mailClient,
	}
}

func (service *TransactionService) PreInternalTransfer(ctx *gin.Context, transferReq model.PreInternalTransferRequest) (string, error) {
	//verify info
	sourceCustomer, sourceAccount, targetAccount, err := service.verifyTransactionInfo(ctx, transferReq.SourceAccountNumber, transferReq.TargetAccountNumber)
	if err != nil {
		return "", err
	}

	//estimate fee
	fee, err := service.coreService.EstimateTransferFee(ctx, transferReq.Amount)
	if err != nil {
		return "", err
	}

	//check is source fee and change balance
	checkFee := *transferReq.IsSourceFee
	if checkFee {
		totalDeduction := transferReq.Amount + fee
		if *sourceAccount.Balance < totalDeduction {
			return "", errors.New("insufficient balance in source account")
		}
		*sourceAccount.Balance = -(totalDeduction)
		*targetAccount.Balance = transferReq.Amount
	} else {
		if *sourceAccount.Balance < transferReq.Amount {
			return "", errors.New("insufficient balance in source account")
		}
		*sourceAccount.Balance = -(transferReq.Amount)
		*targetAccount.Balance = transferReq.Amount - fee
	}

	//store transaction
	transaction := &entity.Transaction{
		SourceAccountNumber: sourceAccount.Number,
		TargetAccountNumber: targetAccount.Number,
		Amount:              transferReq.Amount,
		BankId:              nil,
		Type:                transferReq.Type,
		Description:         transferReq.Description,
		Status:              "pending",
		IsSourceFee:         transferReq.IsSourceFee,
		SourceBalance:       *sourceAccount.Balance,
		TargetBalance:       *targetAccount.Balance,
	}

	//save transaction
	transactionId, err := service.transactionRepository.CreateCommand(ctx, transaction)
	if err != nil {
		return "", err
	}

	//send OTP
	err = service.SendOTPToEmail(ctx, sourceCustomer.Email, transactionId)
	if err != nil {
		return "", err
	}
	return transactionId, nil
}

func (service *TransactionService) SendOTPToEmail(ctx *gin.Context, email string, transactionId string) error {
	// generate otp
	otp := mail.GenerateOTP(6)

	// store otp in redis
	baseKey := constants.VERIFY_TRANSFER_KEY
	number, err := strconv.ParseInt(transactionId, 10, 64)
	if err != nil {
		return err
	}
	key := redis.Concat(baseKey, number)

	err = service.redisClient.Set(ctx, key, otp)
	if err != nil {
		return err
	}

	// send otp to user email
	emailBody := service.mailClient.GenerateOTPBody(email, otp, constants.VERIFY_TRANSFER, constants.VERIFY_TRANSFER_EXP_TIME)
	err = service.mailClient.SendEmail(ctx, email, "OTP verify transfer", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (service *TransactionService) verifyOTP(ctx *gin.Context, transferReq model.InternalTransferRequest) error {
	//regenerate key
	baseKey := constants.VERIFY_TRANSFER_KEY
	number, err := strconv.ParseInt(transferReq.TransactionId, 10, 64)
	if err != nil {
		return err
	}
	key := redis.Concat(baseKey, number)

	//get OTP and check
	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}
	if val != transferReq.Otp {
		return errors.New("invalid otp")
	}

	//delete if match OTP
	err = service.redisClient.Delete(ctx, key)
	if err != nil {
		return err
	}
	return nil
}

func (service *TransactionService) InternalTransfer(ctx *gin.Context, transferReq model.InternalTransferRequest) (*entity.Transaction, error) {
	//get customer and check exists account
	customerId := middleware.GetUserIdHelper(ctx)
	existsAccount, err := service.accountService.GetAccountByCustomerId(ctx, customerId)
	if err != nil {
		return nil, err
	}

	//check transaction by account number and transaction id
	existsTransaction, err := service.transactionRepository.GetTransactionBySourceNumberAndIdQuery(ctx, existsAccount.Number, transferReq.TransactionId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("transaction not found")
		}
		return nil, err
	}

	err = service.verifyOTP(ctx, transferReq)
	if err != nil {
		return nil, err
	}

	//update to DB
	//balance for source and target
	newSourceBalance, err := service.accountService.UpdateBalanceByAccountNumber(ctx, existsTransaction.SourceBalance, existsTransaction.SourceAccountNumber)
	if err != nil {
		return nil, err
	}
	newTargetBalance, err := service.accountService.UpdateBalanceByAccountNumber(ctx, existsTransaction.TargetBalance, existsTransaction.TargetAccountNumber)
	if err != nil {
		return nil, err
	}

	existsTransaction.Status = "success"
	existsTransaction.SourceBalance = newSourceBalance
	existsTransaction.TargetBalance = newTargetBalance

	//transaction
	err = service.transactionRepository.UpdateBalancesCommand(ctx, existsTransaction)
	if err != nil {
		return nil, err
	}
	err = service.transactionRepository.UpdateStatusCommand(ctx, existsTransaction)
	if err != nil {
		return nil, err
	}

	// notify, response history
	return existsTransaction, nil
}

func (service *TransactionService) verifyTransactionInfo(ctx *gin.Context, sourceAccountNumber string, targetAccountNumber string) (*entity.User, *entity.Account, *entity.Account, error) {
	//check input account number
	if sourceAccountNumber == targetAccountNumber {
		return nil, nil, nil, errors.New("source account number can not equal to target account number")
	}

	//get customer and check info
	customerId := middleware.GetUserIdHelper(ctx)
	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil, nil, errors.New("customer not found")
		}
		return nil, nil, nil, err
	}
	//get account by customerId and check sourceNumber
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, sourceCustomer.ID)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil, nil, errors.New("source account not found")
		}
		return nil, nil, nil, err
	}
	if sourceAccount.Number != sourceAccountNumber {
		return nil, nil, nil, errors.New("source account not match")
	}

	//check targetNumber
	targetAccount, err := service.accountService.GetAccountByNumber(ctx, targetAccountNumber)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil, nil, errors.New("target account not found")
		}
		return nil, nil, nil, err
	}
	return sourceCustomer, sourceAccount, targetAccount, nil
}

func (service *TransactionService) AddDebtReminder(ctx *gin.Context, debtReminder model.DebtReminderRequest) error {
	//verify info
	_, sourceAccount, targetAccount, err := service.verifyTransactionInfo(ctx, debtReminder.SourceAccountNumber, debtReminder.TargetAccountNumber)
	if err != nil {
		return err
	}

	//estimate fee
	fee, err := service.coreService.EstimateTransferFee(ctx, debtReminder.Amount)
	if err != nil {
		return err
	}

	*sourceAccount.Balance = debtReminder.Amount
	*targetAccount.Balance = -(debtReminder.Amount + fee)
	falseStatus := false
	//store transaction
	transaction := &entity.Transaction{
		SourceAccountNumber: sourceAccount.Number,
		TargetAccountNumber: targetAccount.Number,
		Amount:              debtReminder.Amount,
		BankId:              nil,
		Type:                debtReminder.Type,
		Description:         debtReminder.Description,
		Status:              "pending",
		IsSourceFee:         &falseStatus,
		SourceBalance:       *sourceAccount.Balance,
		TargetBalance:       *targetAccount.Balance,
	}

	//save transaction
	_, err = service.transactionRepository.CreateCommand(ctx, transaction)
	if err != nil {
		return err
	}
	return nil
}
