package serviceimplement

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type DebtReplyService struct {
	debtReplyRepo repository.DebtReplyRepository
}

func NewDebtReplyService(repo repository.DebtReplyRepository) service.DebtReplyService {
	return &DebtReplyService{
		debtReplyRepo: repo,
	}
}

func (service *DebtReplyService) GetReplyByDebtReminderId(ctx *gin.Context, debtReminderId string) (*entity.DebtReply, error) {
	reply, err := service.debtReplyRepo.GetReplyByDebtIdQuery(ctx, debtReminderId)
	if err != nil {
		return nil, err
	}
	return reply, nil
}
