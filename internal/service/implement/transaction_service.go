package serviceimplement

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/HMAC_signature"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"

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

type TransferRSATeamResponse struct {
	Success *bool    `json:"success"`
	Message []string `json:"message"`
	Data    string   `json:"data"`
}

type TransactionService struct {
	transactionRepository  repository.TransactionRepository
	customerRepository     repository.CustomerRepository
	accountService         service.AccountService
	coreService            service.CoreService
	redisClient            bean.RedisClient
	mailClient             bean.MailClient
	debtReply              repository.DebtReplyRepository
	notificationRepository repository.NotificationRepository
	notificationClient     bean.NotificationClient
	partnerBankService     service.PartnerBankService
	rsaMiddleware          *middleware.RSAMiddleware
}

func NewTransactionService(transactionRepository repository.TransactionRepository,
	customerRepository repository.CustomerRepository,
	accountService service.AccountService,
	coreService service.CoreService,
	redisClient bean.RedisClient,
	mailClient bean.MailClient,
	debtReply repository.DebtReplyRepository,
	notificationRepository repository.NotificationRepository,
	notificationClient bean.NotificationClient,
	partnerBankService service.PartnerBankService,
	rsaMiddleware *middleware.RSAMiddleware,
) service.TransactionService {
	return &TransactionService{
		transactionRepository:  transactionRepository,
		customerRepository:     customerRepository,
		accountService:         accountService,
		coreService:            coreService,
		redisClient:            redisClient,
		mailClient:             mailClient,
		debtReply:              debtReply,
		notificationRepository: notificationRepository,
		notificationClient:     notificationClient,
		partnerBankService:     partnerBankService,
		rsaMiddleware:          rsaMiddleware,
	}
}

func (service *TransactionService) PreInternalTransfer(ctx *gin.Context, transferReq model.PreInternalTransferRequest) (string, error) {
	//verify info
	sourceCustomer, sourceAccount, targetAccount, err := service.verifyTransactionInfo(ctx, transferReq.SourceAccountNumber, transferReq.TargetAccountNumber)
	if err != nil {
		return "", err
	}

	//check type
	if transferReq.Type != "internal" {
		return "", errors.New("invalid transaction type")
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

	fmt.Println("OTP: ", otp)

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

func (service *TransactionService) verifyOTP(ctx *gin.Context, transferReq model.TransferRequest) error {
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

func (service *TransactionService) InternalTransfer(ctx *gin.Context, transferReq model.TransferRequest) (*entity.Transaction, error) {
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

	// create notification response
	//get source customer and target customer's name
	sourceCustomer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, existsTransaction.SourceAccountNumber)
	if err != nil {
		fmt.Println(err)
	}
	targetCustomer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, existsTransaction.TargetAccountNumber)
	if err != nil {
		fmt.Println(err)
	}

	notificationForSourceCustomerResp := &model.TransactionNotificationContent{
		DeviceId:      int(sourceCustomer.Id),
		Name:          targetCustomer.Name,
		Amount:        int(existsTransaction.Amount),
		TransactionId: existsTransaction.Id,
		Type:          "outgoing_transfer",
		CreatedAt:     existsTransaction.UpdatedAt,
	}

	notificationForTargetCustomerResp := &model.TransactionNotificationContent{
		DeviceId:      int(targetCustomer.Id),
		Name:          sourceCustomer.Name,
		Amount:        int(existsTransaction.Amount),
		TransactionId: existsTransaction.Id,
		Type:          "incoming_transfer",
		CreatedAt:     existsTransaction.UpdatedAt,
	}

	// notify, response history
	service.notificationClient.SaveAndSend(ctx, *notificationForSourceCustomerResp)
	service.notificationClient.SaveAndSend(ctx, *notificationForTargetCustomerResp)

	return existsTransaction, nil
}

func (service *TransactionService) verifyTransactionInfo(ctx *gin.Context, sourceAccountNumber string, targetAccountNumber string) (*entity.User, *entity.Account, *entity.Account, error) {
	//check input account number
	if sourceAccountNumber == targetAccountNumber {
		return nil, nil, nil, errors.New("source account number can not equal to target account number")
	}

	//get customer
	customerId := middleware.GetUserIdHelper(ctx)
	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil, nil, errors.New("customer not found")
		}
		return nil, nil, nil, err
	}
	//get account by customerId and check sourceNumber is internal
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, sourceCustomer.Id)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil, nil, errors.New("source account not found")
		}
		return nil, nil, nil, err
	}
	if sourceAccount.Number != sourceAccountNumber {
		return nil, nil, nil, errors.New("source account not match")
	}

	//check targetNumber is internal
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
	trueStatus := true
	//store transaction
	transaction := &entity.Transaction{
		SourceAccountNumber: targetAccount.Number,
		TargetAccountNumber: sourceAccount.Number,
		Amount:              debtReminder.Amount,
		BankId:              nil,
		Type:                debtReminder.Type,
		Description:         debtReminder.Description,
		Status:              "pending",
		IsSourceFee:         &trueStatus,
		SourceBalance:       *targetAccount.Balance,
		TargetBalance:       *sourceAccount.Balance,
	}

	//save transaction
	targetTransactionId, err := service.transactionRepository.CreateCommand(ctx, transaction)
	if err != nil {
		return err
	}

	targetTransaction, err := service.transactionRepository.GetTransactionByIdQuery(ctx, targetTransactionId)

	//get target and source customer's name
	sourceCustomer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, transaction.SourceAccountNumber)
	if err != nil {
		fmt.Println(err)
	}
	targetCustomer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, transaction.TargetAccountNumber)
	if err != nil {
		fmt.Println(err)
	}

	// create notification response
	notificationForTargetCustomerResp := &model.TransactionNotificationContent{
		DeviceId:      int(sourceCustomer.Id),
		Name:          targetCustomer.Name,
		Amount:        int(transaction.Amount),
		TransactionId: transaction.Id,
		Type:          "debt_reminder",
		CreatedAt:     targetTransaction.CreatedAt,
	}
	fmt.Println(notificationForTargetCustomerResp)

	// notify, response history
	service.notificationClient.SaveAndSend(ctx, *notificationForTargetCustomerResp)

	return nil
}

func (service *TransactionService) CancelDebtReminder(ctx *gin.Context, debtReminderId string, debtReply model.DebtReminderReplyRequest) error {
	//get customer and check info
	customerId := middleware.GetUserIdHelper(ctx)
	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("customer not found")
		}
		return err
	}

	//get account by customerId and check sourceNumber is internal
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, sourceCustomer.Id)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("source account not found")
		}
		return err
	}

	//get transaction by debtReminderId
	debtReminder, err := service.transactionRepository.GetTransactionByIdQuery(ctx, debtReminderId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("debt reminder not found")
		}
		return err
	}
	//check accountNumber in debtReminder
	if sourceAccount.Number != debtReminder.SourceAccountNumber && sourceAccount.Number != debtReminder.TargetAccountNumber {
		return errors.New("account number not match")
	}

	//create reply
	reply := &entity.DebtReply{
		Content:        debtReply.Content,
		DebtReminderId: debtReminderId,
		UserReplyName:  sourceCustomer.Name,
	}

	//save reply
	err = service.debtReply.CreateCommand(ctx, reply)
	if err != nil {
		return err
	}
	//update status and save
	debtReminder.Status = "failed"
	err = service.transactionRepository.UpdateStatusCommand(ctx, debtReminder)
	if err != nil {
		return err
	}

	//notify...
	//if current user is source, then find the other user to notify
	if sourceAccount.Number == debtReminder.SourceAccountNumber {
		creditor, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, debtReminder.TargetAccountNumber)
		if err != nil {
			fmt.Println(err)
		}
		debtor, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, debtReminder.SourceAccountNumber)
		if err != nil {
			fmt.Println(err)
		}

		// create notification response
		notificationForTargetCustomerResp := &model.TransactionNotificationContent{
			DeviceId:      int(creditor.Id),
			Name:          debtor.Name,
			Amount:        int(debtReminder.Amount),
			TransactionId: debtReminder.Id,
			Type:          "debt_cancel",
			CreatedAt:     debtReminder.CreatedAt,
		}

		// notify, response history
		service.notificationClient.SaveAndSend(ctx, *notificationForTargetCustomerResp)
	} else {
		//if current user is target, then find the other user to notify
		creditor, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, debtReminder.TargetAccountNumber)
		if err != nil {
			fmt.Println(err)
		}
		debtor, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, debtReminder.SourceAccountNumber)
		if err != nil {
			fmt.Println(err)
		}

		// create notification response
		notificationForSourceCustomerResp := &model.TransactionNotificationContent{
			DeviceId:      int(debtor.Id),
			Name:          creditor.Name,
			Amount:        int(debtReminder.Amount),
			TransactionId: debtReminder.Id,
			Type:          "debt_cancel",
			CreatedAt:     debtReminder.CreatedAt,
		}

		// notify, response history
		service.notificationClient.SaveAndSend(ctx, *notificationForSourceCustomerResp)
	}
	return nil
}

func (service *TransactionService) GetReceivedDebtReminder(ctx *gin.Context) ([]model.DebtReminderResponse, error) {
	//get customer and check info
	customerId := middleware.GetUserIdHelper(ctx)
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}
	//get debt receive
	debtList, err := service.transactionRepository.GetReceivedDebtReminderByCustomerIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil
		}
	}
	var resList []model.DebtReminderResponse
	for _, debt := range *debtList {
		sender, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, debt.TargetAccountNumber)
		if err != nil {
			if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
				return nil, errors.New("customer not found")
			}
		}
		if debt.Status == "failed" {
			reply, err := service.debtReply.GetReplyByDebtIdQuery(ctx, debt.Id)
			if err != nil {
				return nil, err
			}
			resList = append(resList, model.DebtReminderResponse{
				Sender:       sender.Name,
				Receiver:     sourceCustomer.Name,
				DebtReminder: &debt,
				Reply:        reply,
			})
		} else {
			resList = append(resList, model.DebtReminderResponse{
				Sender:       sender.Name,
				Receiver:     sourceCustomer.Name,
				DebtReminder: &debt,
				Reply:        nil,
			})
		}
	}
	return resList, nil
}

func (service *TransactionService) GetSentDebtReminder(ctx *gin.Context) ([]model.DebtReminderResponse, error) {
	customerId := middleware.GetUserIdHelper(ctx)
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}
	//get debt sent
	debtList, err := service.transactionRepository.GetSentDebtReminderByCustomerIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, nil
		}
	}
	var resList []model.DebtReminderResponse
	for _, debt := range *debtList {
		receiver, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, debt.SourceAccountNumber)
		if err != nil {
			if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
				return nil, errors.New("customer not found")
			}
		}
		if debt.Status == "failed" {
			reply, err := service.debtReply.GetReplyByDebtIdQuery(ctx, debt.Id)
			if err != nil {
				return nil, err
			}
			resList = append(resList, model.DebtReminderResponse{
				Sender:       sourceCustomer.Name,
				Receiver:     receiver.Name,
				DebtReminder: &debt,
				Reply:        reply,
			})
		} else {
			resList = append(resList, model.DebtReminderResponse{
				Sender:       sourceCustomer.Name,
				Receiver:     receiver.Name,
				DebtReminder: &debt,
				Reply:        nil,
			})
		}
	}
	return resList, nil
}

func (s *TransactionService) GetTransactionsByCustomerId(ctx *gin.Context, customerId int64) ([]model.GetTransactionsResponseSum, error) {
	//get account by customerId
	sourceAccount, err := s.accountService.GetAccountByCustomerId(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("source account not found")
		}
		return nil, err
	}

	transactions, err := s.transactionRepository.GetTransactionByAccountNumber(ctx, sourceAccount.Number)
	if err != nil {
		return nil, err
	}

	transactionResp := make([]model.GetTransactionsResponseSum, 0)

	for _, transaction := range transactions {
		trans := service.TransactionUtilsEntityToResponse(transaction, sourceAccount.Number)
		var bankInfo *entity.PartnerBank
		if transaction.BankId != nil {
			bankInfo, _ = s.partnerBankService.GetBankById(ctx, *transaction.BankId)
			transactionResp = append(transactionResp, model.GetTransactionsResponseSum{
				Transaction: trans,
				BankCode:    &bankInfo.BankCode,
				BankName:    &bankInfo.BankName,
			})
		} else {
			transactionResp = append(transactionResp, model.GetTransactionsResponseSum{
				Transaction: trans,
				BankCode:    nil,
				BankName:    nil,
			})
		}
	}

	return transactionResp, nil
}

func (s *TransactionService) GetTransactionByIdAndCustomerId(ctx *gin.Context, customerId int64, id string) (*model.GetTransactionsResponseSum, error) {
	//get account by customerId
	sourceAccount, err := s.accountService.GetAccountByCustomerId(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("source account not found")
		}
		return nil, err
	}

	transaction, err := s.transactionRepository.GetTransactionByAccountNumberAndIdQuery(ctx, sourceAccount.Number, id)
	if err != nil {
		return nil, err
	}
	transactionResp := service.TransactionUtilsEntityToResponse(*transaction, sourceAccount.Number)
	var bankInfo *entity.PartnerBank
	var transactionRes *model.GetTransactionsResponseSum
	if transaction.BankId != nil {
		bankInfo, _ = s.partnerBankService.GetBankById(ctx, *transaction.BankId)
		transactionRes = &model.GetTransactionsResponseSum{
			Transaction: transactionResp,
			BankCode:    &bankInfo.BankCode,
			BankName:    &bankInfo.BankName,
		}
	} else {
		transactionRes = &model.GetTransactionsResponseSum{
			Transaction: transactionResp,
			BankCode:    nil,
			BankName:    nil,
		}
	}

	return transactionRes, nil
}

func (service *TransactionService) PreDebtTransfer(ctx *gin.Context, transferReq model.PreDebtTransferRequest) error {
	//get customer
	customerId := middleware.GetUserIdHelper(ctx)
	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("customer not found")
		}
		return err
	}
	//get account by customerId
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, sourceCustomer.Id)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("source account not found")
		}
		return err
	}
	//get transaction by id and check valid
	transaction, err := service.transactionRepository.GetTransactionByIdQuery(ctx, transferReq.TransactionId)
	if err != nil {
		return err
	}
	if transaction.SourceAccountNumber != sourceAccount.Number {
		return errors.New("source account not match")
	}
	//check balance
	if *sourceAccount.Balance < -(transaction.SourceBalance) {
		return errors.New("insufficient balance in source account")
	}

	//send OTP
	err = service.SendOTPToEmail(ctx, sourceCustomer.Email, transaction.Id)
	if err != nil {
		return err
	}
	return nil
}

func (service *TransactionService) ReceiveExternalTransfer(ctx *gin.Context, transferReq model.ExternalPayload, partnerBankId int64) error {
	//update to DB
	//balance for target account
	newTargetBalance, err := service.accountService.UpdateBalanceByAccountNumber(ctx, transferReq.Amount, transferReq.DesAccountNumber)
	if err != nil {
		return err
	}
	transaction := &entity.Transaction{
		SourceAccountNumber: transferReq.SrcAccountNumber,
		TargetAccountNumber: transferReq.DesAccountNumber,
		Amount:              transferReq.Amount,
		BankId:              &partnerBankId,
		Type:                "external",
		Description:         transferReq.Description,
		Status:              "success",
		IsSourceFee:         transferReq.IsSourceFee,
		SourceBalance:       -1,
		TargetBalance:       newTargetBalance,
	}
	//save transaction
	transactionId, err := service.transactionRepository.CreateCommand(ctx, transaction)
	if err != nil {
		return err
	}
	//check notify
	//find user by desAccountNumber -> notify
	existsTransaction, err := service.transactionRepository.GetTransactionByIdQuery(ctx, transactionId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			fmt.Println(err)
		}
		fmt.Println(err)
	}

	targetCustomer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, existsTransaction.TargetAccountNumber)
	if err != nil {
		fmt.Println(err)
	}

	notificationForTargetCustomerResp := &model.TransactionNotificationContent{
		DeviceId:      int(targetCustomer.Id),
		Name:          "ABC", // will call api get source customer name later
		Amount:        int(existsTransaction.Amount),
		TransactionId: existsTransaction.Id,
		Type:          "incoming_transfer",
		CreatedAt:     existsTransaction.UpdatedAt,
	}

	// notify, response history
	service.notificationClient.SaveAndSend(ctx, *notificationForTargetCustomerResp)
	return nil
}

func (service *TransactionService) PreExternalTransfer(ctx *gin.Context, transferReq model.PreExternalTransferRequest) (string, error) {
	//get customer
	customerId := middleware.GetUserIdHelper(ctx)
	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return "", errors.New("customer not found")
		}
		return "", err
	}
	//get account by customerId and check sourceNumber is internal
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, sourceCustomer.Id)
	if err != nil {
		return "", err
	}
	if sourceAccount.Number != transferReq.SourceAccountNumber {
		return "", errors.New("source account not match")
	}
	//check targetNumber is in partner bank
	//_, err = service.accountService.GetExternalAccountName(ctx, model.GetExternalAccountNameRequest{
	//	BankId:        transferReq.PartnerBankId,
	//	AccountNumber: transferReq.TargetAccountNumber,
	//})
	//
	//if err != nil {
	//	return "", err
	//}

	//check type
	if transferReq.Type != "external" {
		return "", errors.New("invalid transaction type")
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
	} else {
		if *sourceAccount.Balance < transferReq.Amount {
			return "", errors.New("insufficient balance in source account")
		}
		*sourceAccount.Balance = -(transferReq.Amount)
	}

	//store transaction
	transaction := &entity.Transaction{
		SourceAccountNumber: sourceAccount.Number,
		TargetAccountNumber: transferReq.TargetAccountNumber,
		Amount:              transferReq.Amount,
		BankId:              &transferReq.PartnerBankId,
		Type:                transferReq.Type,
		Description:         transferReq.Description,
		Status:              "pending",
		IsSourceFee:         transferReq.IsSourceFee,
		SourceBalance:       *sourceAccount.Balance,
		TargetBalance:       -1,
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

func (service *TransactionService) ExternalTransfer(ctx *gin.Context, transferReq model.TransferRequest) (*entity.Transaction, error) {
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
	//verify OTP
	err = service.verifyOTP(ctx, transferReq)
	if err != nil {
		return nil, err
	}
	/*
		call api to partner bank to confirm transaction
	*/
	//err = service.callExternalTransfer(ctx, existsTransaction)
	//if err != nil {
	//	return nil, err
	//}
	//update to DB
	//balance for source
	newSourceBalance, err := service.accountService.UpdateBalanceByAccountNumber(ctx, existsTransaction.SourceBalance, existsTransaction.SourceAccountNumber)
	if err != nil {
		return nil, err
	}

	existsTransaction.Status = "success"
	existsTransaction.SourceBalance = newSourceBalance

	//transaction
	err = service.transactionRepository.UpdateBalancesCommand(ctx, existsTransaction)
	if err != nil {
		return nil, err
	}
	err = service.transactionRepository.UpdateStatusCommand(ctx, existsTransaction)
	if err != nil {
		return nil, err
	}
	/*
		notify transfer success
	*/
	sourceCustomer, err := service.customerRepository.GetCustomerByAccountNumberQuery(ctx, existsTransaction.SourceAccountNumber)
	if err != nil {
		fmt.Println(err)
	}
	//get target name
	//targetAccountName, err := service.accountService.GetExternalAccountName(ctx, model.GetExternalAccountNameRequest{
	//	BankId:        *existsTransaction.BankId,
	//	AccountNumber: existsTransaction.TargetAccountNumber,
	//})
	if err != nil {
		return nil, err
	}
	//noti
	notificationForSourceCustomerResp := &model.TransactionNotificationContent{
		DeviceId:      int(sourceCustomer.Id),
		Name:          "HUYNH THIEN HUU",
		Amount:        int(existsTransaction.Amount),
		TransactionId: existsTransaction.Id,
		Type:          "outgoing_transfer",
		CreatedAt:     existsTransaction.UpdatedAt,
	}

	// notify, response history
	service.notificationClient.SaveAndSend(ctx, *notificationForSourceCustomerResp)

	// response history
	return existsTransaction, nil
}

func (service *TransactionService) callExternalTransfer(ctx *gin.Context, transaction *entity.Transaction) error {
	//setup payload
	bankIdInRsaTeamInt, err := strconv.ParseInt(bankIdTeam3, 10, 64)
	if err != nil {
		return err
	}
	//get partner bank info
	partnerBank, err := service.partnerBankService.GetBankById(ctx, *transaction.BankId)
	if err != nil {
		return err
	}
	//setup payload and call api
	utcTime := time.Now().UTC()
	formattedTime := fmt.Sprintf("%04d-%02d-%02dT%02d:%02d:%02d.%09dZ",
		utcTime.Year(), utcTime.Month(), utcTime.Day(),
		utcTime.Hour(), utcTime.Minute(), utcTime.Second(),
		utcTime.Nanosecond())
	payload := &model.ExternalTransferToRSATeam{
		BankId:               bankIdInRsaTeamInt,
		AccountNumber:        transaction.TargetAccountNumber,
		ForeignAccountNumber: transaction.SourceAccountNumber,
		Amount:               transaction.Amount,
		Description:          transaction.Description,
		Timestamp:            formattedTime,
	}
	reqBytes, err := json.Marshal(payload)
	//hash data
	secretString := os.Getenv("SECRET_KEY_FOR_EXTERNAL_BANK")
	hashedData := HMAC_signature.GenerateHMAC(string(reqBytes), secretString)
	//sign data
	signedData, err := service.rsaMiddleware.SignDataRSA(string(reqBytes))
	if err != nil {
		return err
	}
	//setup and call to partner bank server
	request, err := http.NewRequest("POST", partnerBank.TransferApi, bytes.NewBuffer(reqBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("HMAC", hashedData)
	request.Header.Set("RSA-Signature", signedData)
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)
	//handler response
	var res TransferRSATeamResponse
	err = json.Unmarshal(body, &res)
	if err != nil {
		return err
	}
	if *res.Success != true {
		return errors.New(res.Message[0])
	}
	//Verify signature
	partnerSignature := response.Header.Get("RSA-Signature")
	dataBytes, _ := json.Marshal(res)
	err = service.rsaMiddleware.VerifyRSASignature(*partnerBank, dataBytes, partnerSignature)
	if err != nil {
		return err
	}
	return nil
}
