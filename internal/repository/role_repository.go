package repository

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type RoleRepository interface {
	GetByUserId(ctx *gin.Context, userId int64) (*entity.Role, error)
}
