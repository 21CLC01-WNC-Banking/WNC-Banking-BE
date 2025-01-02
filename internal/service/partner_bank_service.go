package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type PartnerBankService interface {
	AddPartnerBank(c *gin.Context, request model.PartnerBankRequest) error
	GetPartnerBankByBankCode(c *gin.Context, bankCode string) (*entity.PartnerBank, error)
	GetListPartnerBank(c *gin.Context) ([]entity.PartnerBank, error)
}
