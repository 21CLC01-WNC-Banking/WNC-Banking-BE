package serviceimplement

import (
	"errors"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type PartnerBankService struct {
	partnerBankRepo repository.PartnerBankRepository
}

func NewPartnerBankService(partnerBankRepo repository.PartnerBankRepository) service.PartnerBankService {
	return &PartnerBankService{partnerBankRepo: partnerBankRepo}
}

func (service *PartnerBankService) AddPartnerBank(c *gin.Context, request model.PartnerBankRequest) error {
	partnerBank := &entity.PartnerBank{
		BankCode:    request.BankCode,
		BankName:    request.BankName,
		ShortName:   request.ShortName,
		LogoUrl:     request.LogoUrl,
		ResearchApi: request.ResearchApi,
		TransferApi: request.TransferApi,
		PublicKey:   request.PublicKey,
	}
	err := service.partnerBankRepo.CreateCommand(c, partnerBank)
	if err != nil {
		return err
	}
	return nil
}

func (service *PartnerBankService) GetPartnerBankByBankCode(c *gin.Context, bankCode string) (*entity.PartnerBank, error) {
	partnerBank, err := service.partnerBankRepo.GetOneByBankCode(c, bankCode)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("partner bank not found")
		}
		return nil, err
	}
	return partnerBank, nil
}

func (service *PartnerBankService) GetListPartnerBank(c *gin.Context) ([]entity.PartnerBank, error) {
	return service.partnerBankRepo.GetListBank(c)
}

func (service *PartnerBankService) GetBankById(c *gin.Context, bankId int64) (*entity.PartnerBank, error) {
	partnerBank, err := service.partnerBankRepo.GetBankById(c, bankId)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("partner bank not found")
		}
		return nil, err
	}
	return partnerBank, nil
}
