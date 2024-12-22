package serviceimplement

import (
	"errors"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type SavedReceiverService struct {
	savedReceiverRepository repository.SavedReceiverRepository
	accountService          service.AccountService
}

func NewSavedReceiverService(savedReceiverRepository repository.SavedReceiverRepository, accountService service.AccountService) service.SavedReceiverService {
	return &SavedReceiverService{
		savedReceiverRepository: savedReceiverRepository,
		accountService:          accountService,
	}
}

func (service *SavedReceiverService) AddInternalReceiver(ctx *gin.Context, receiver model.InternalReceiver) error {
	customerId, exists := ctx.Get("userId")
	if !exists {
		return errors.New("customer not exists")
	}

	account, err := service.accountService.GetAccountByCustomerId(ctx, customerId.(int64))
	if err != nil {
		return err
	}
	if account.Number == receiver.ReceiverAccountNumber {
		return errors.New("cannot add yourself as receiver")
	}

	exists, err = service.existsByAccountNumberAndBankID(ctx, receiver.ReceiverAccountNumber, nil)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("receiver already exists")
	}

	err = service.savedReceiverRepository.CreateCommand(ctx, &entity.SavedReceiver{
		CustomerId:            customerId.(int64),
		ReceiverAccountNumber: receiver.ReceiverAccountNumber,
		ReceiverNickname:      receiver.ReceiverNickname,
		BankId:                nil,
	})
	if err != nil {
		return err
	}
	return nil
}

func (service *SavedReceiverService) AddExternalReceiver(ctx *gin.Context, receiver model.ExternalReceiver) error {
	panic("unimplemented")
}

func (service *SavedReceiverService) existsByAccountNumberAndBankID(ctx *gin.Context, accountNumber string, bankID *int64) (bool, error) {
	return service.savedReceiverRepository.ExistsByAccountNumberAndBankID(ctx, accountNumber, bankID)
}
