package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type AdminService interface {
	GetAllStaff(ctx *gin.Context) ([]entity.User, error)
	GetOneStaff(ctx *gin.Context, staffId int64) (*entity.User, error)
	CreateOneStaff(ctx *gin.Context, request *model.CreateStaffRequest) (int64, error) // return created staff id
	DeleteOneStaff(ctx *gin.Context, staffId int64) error
	UpdateOneStaff(ctx *gin.Context, request *model.UpdateStaffRequest) error
}
