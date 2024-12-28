package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/gin-gonic/gin"
)

type AdminService interface {
	GetAllStaff(ctx *gin.Context) ([]entity.User, error)
	GetOneStaff(ctx *gin.Context, staffId int64) (*entity.User, error)
}
