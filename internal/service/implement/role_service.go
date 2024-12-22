package serviceimplement

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/gin-gonic/gin"
)

type RoleService struct {
	roleRepository repository.RoleRepository
}

func NewRoleService(roleRepository repository.RoleRepository) service.RoleService {
	return &RoleService{
		roleRepository: roleRepository,
	}
}

func (r *RoleService) GetRoleByUserId(ctx *gin.Context, userId int64) (*entity.Role, error) {
	return r.roleRepository.GetByUserId(ctx, userId)
}
