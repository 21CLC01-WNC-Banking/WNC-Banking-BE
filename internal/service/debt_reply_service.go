package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type DebtReplyService interface {
	GetReplyByDebtReminderId(ctx *gin.Context, debtReminderId string) (*entity.DebtReply, error)
}
