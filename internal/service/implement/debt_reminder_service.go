package serviceimplement

import (
	"errors"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type DebtReminderService struct {
	transactionRepo    repository.TransactionRepository
	customerRepository repository.CustomerRepository
	accountService     service.AccountService
}

func NewDebtReminderService(
	transactionRepo repository.TransactionRepository,
	customerRepository repository.CustomerRepository,
	accountService service.AccountService,
) service.DebtReminderService {
	return &DebtReminderService{
		transactionRepo:    transactionRepo,
		customerRepository: customerRepository,
		accountService:     accountService,
	}
}

func (service *DebtReminderService) AddNewDebtReminder(ctx *gin.Context, debtReq *model.DebtReminderRequest) error {
	//check input account number
	if debtReq.SourceAccountNumber == debtReq.TargetAccountNumber {
		return errors.New("source account number can not equal to target account number")
	}

	//get customer and check info
	customerId, exists := ctx.Get("userId")
	if !exists {
		return errors.New("customer not exists")
	}

	//check customerId
	sourceCustomer, err := service.customerRepository.GetOneByIdQuery(ctx, customerId.(int64))
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("customer not found")
		}
		return err
	}
	//get account by id and check sourceNumber
	sourceAccount, err := service.accountService.GetAccountByCustomerId(ctx, sourceCustomer.ID)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("source account not found")
		}
		return err
	}
	if sourceAccount.Number != debtReq.SourceAccountNumber {
		return errors.New("source account not match")
	}

	//check targetNumber
	_, err = service.accountService.GetAccountByNumber(ctx, debtReq.TargetAccountNumber)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return errors.New("target account not found")
		}
		return err
	}
	return nil
}
