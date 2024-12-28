package serviceimplement

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type AdminService struct {
	staffRepository repository.StaffRepository
}

func NewAdminService(staffRepository repository.StaffRepository) service.AdminService {
	return &AdminService{staffRepository: staffRepository}
}

func (a AdminService) GetAllStaff(ctx *gin.Context) ([]entity.User, error) {
	return a.staffRepository.GetAll(ctx)
}

func (a AdminService) GetOneStaff(ctx *gin.Context, staffId int64) (*entity.User, error) {
	return a.staffRepository.GetOneById(ctx, staffId)
}
