package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Login(ctx *gin.Context, customerRequest model.LoginRequest) (*entity.User, error)
	ValidateRefreshToken(ctx *gin.Context, customerId int64) (*entity.Authentication, error)

	SendOTPToEmail(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error
	VerifyOTP(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error
	SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) error
	GetUserById(ctx *gin.Context) (*entity.User, error)
}
