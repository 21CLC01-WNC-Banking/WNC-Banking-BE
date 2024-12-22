package serviceimplement

import (
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/constants"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/generate_number_code"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/mail"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/redis"
	"github.com/gin-gonic/gin"
	"strconv"
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
	//check input account number
	if transferReq.SourceAccountNumber == transferReq.TargetAccountNumber {
		return "", errors.New("source account number can not equal to target account number")
	}

	//get customer and check info
	customerId, exists := ctx.Get("customerId")
	if !exists {
		return "", errors.New("customer not exists")
	}

	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId.(int64))
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return "", errors.New("customer not found")
		}
		return "", err
	}
	//get account by customerId and check sourceNumber
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, customerId.(int64))
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return "", errors.New("source account not found")
		}
		return "", err
	}
	if sourceAccount.Number != transferReq.SourceAccountNumber {
		return "", errors.New("source account not match")
	}

	//check targetNumber
	targetAccount, err := service.accountService.GetAccountByNumber(ctx, transferReq.TargetAccountNumber)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return "", errors.New("target account not found")
		}
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
		if sourceAccount.Balance < totalDeduction {
			return "", errors.New("insufficient balance in source account")
		}
		sourceAccount.Balance -= totalDeduction
		targetAccount.Balance += transferReq.Amount
	} else {
		if sourceAccount.Balance < transferReq.Amount {
			return "", errors.New("insufficient balance in source account")
		}
		sourceAccount.Balance -= transferReq.Amount
		targetAccount.Balance += transferReq.Amount - fee
	}
	//generate id
	transactionId := generate_number_code.GenerateRandomNumber(10)

	//store transaction
	transaction := &entity.Transaction{
		Id:                  transactionId,
		SourceAccountNumber: sourceAccount.Number,
		TargetAccountNumber: targetAccount.Number,
		Amount:              transferReq.Amount,
		BankId:              nil,
		Type:                transferReq.Type,
		Description:         transferReq.Description,
		Status:              "pending",
		IsSourceFee:         transferReq.IsSourceFee,
		SourceBalance:       sourceAccount.Balance,
		TargetBalance:       targetAccount.Balance,
	}

	//send OTP
	err = service.SendOTPToEmail(ctx, sourceCustomer.Email, transactionId)
	if err != nil {
		return "", err
	}

	//save transaction
	err = service.transactionRepository.CreateCommand(ctx, transaction)
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
	err = service.mailClient.SendEmail(ctx, email, "OTP verify transfer", otp, constants.VERIFY_TRANSFER, constants.VERIFY_TRANSFER_EXP_TIME)
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
	customerId, exists := ctx.Get("customerId")
	if !exists {
		return nil, errors.New("customer not exists")
	}
	existsAccount, err := service.accountService.GetAccountByCustomerId(ctx, customerId.(int64))
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
	existsTransaction.Status = "success"

	//update to DB
	//transaction
	err = service.transactionRepository.UpdateStatusCommand(ctx, existsTransaction)
	if err != nil {
		return nil, err
	}

	//balance for source and target
	err = service.accountService.UpdateBalanceByAccountNumber(ctx, existsTransaction.SourceBalance, existsTransaction.SourceAccountNumber)
	if err != nil {
		return nil, err
	}
	err = service.accountService.UpdateBalanceByAccountNumber(ctx, existsTransaction.TargetBalance, existsTransaction.TargetAccountNumber)
	if err != nil {
		return nil, err
	}

	// notify, response history
	return existsTransaction, nil
}
