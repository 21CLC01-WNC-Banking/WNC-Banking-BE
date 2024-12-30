package serviceimplement

import (
	"errors"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

const FEE_PERCENTAGE = 2 // 2%
// const MIN_FEE = 2000     // const fee if amount <= 200k
const MIN_FEE = 2
const THRESHOLD = 200000 // limit 200k

type CoreService struct{}

func NewCoreService() service.CoreService {
	return &CoreService{}
}

func (service *CoreService) EstimateTransferFee(ctx *gin.Context, amount int64) (int64, error) {
	if amount < MIN_FEE {
		return 0, errors.New("amount must be greater than or equal to MIN_FEE")
	}
	if amount <= THRESHOLD {
		return MIN_FEE, nil
	}
	fee := (amount * FEE_PERCENTAGE) / 100

	return fee, nil
}
