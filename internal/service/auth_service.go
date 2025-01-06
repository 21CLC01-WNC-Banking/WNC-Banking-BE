package service

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Login(ctx *gin.Context, customerRequest model.LoginRequest) (*model.LoginResponse, error)
	ValidateRefreshToken(ctx *gin.Context, customerId int64) (*entity.Authentication, error)

	SendOTPToEmail(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error
	VerifyOTP(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error
	SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) error
	GetUserById(ctx *gin.Context, userId int64) (*entity.User, error)
	Logout(ctx *gin.Context, refreshToken string)
	ChangePassword(ctx *gin.Context, request model.ChangePasswordRequest) error
}
