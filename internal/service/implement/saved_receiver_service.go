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
}

func NewSavedReceiverService(savedReceiverRepository repository.SavedReceiverRepository) service.SavedReceiverService {
	return &SavedReceiverService{savedReceiverRepository: savedReceiverRepository}
}

func (service *SavedReceiverService) AddInternalReceiver(ctx *gin.Context, receiver model.InternalReceiver) error {
	customerId, exists := ctx.Get("customerId")
	if !exists {
		return errors.New("customer not exists")
	}

	err := service.savedReceiverRepository.CreateCommand(ctx, &entity.SavedReceiver{
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
