package v1

import (
	"errors"
	"net/http"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/validation"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary Login
// @Description Login to account
// @Tags Auths
// @Accept json
// @Param request body model.LoginRequest true "Auth payload"
// @Produce  json
// @Router /auth/login [post]
// @Success 200 {object} httpcommon.HttpResponse[entity.User]
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) Login(ctx *gin.Context) {
	var loginRequest model.LoginRequest

	if err := validation.BindJsonAndValidate(ctx, &loginRequest); err != nil {
		return
	}

	customer, err := handler.authService.Login(ctx, loginRequest)
	if err != nil || customer == nil {
		if customer == nil && err == nil {
			err = errors.New(httpcommon.ErrorMessage.SqlxNoRow)
		}
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	ctx.JSON(200, httpcommon.NewSuccessResponse(customer))
}

// @Summary Send OTP to Mail
// @Description Send OTP to user email
// @Tags Auths
// @Accept json
// @Param request body model.SendOTPRequest true "Send OTP payload"
// @Produce json
// @Router /auth/forgot-password/otp [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SendOTPToMail(ctx *gin.Context) {
	var sendOTPRequest model.SendOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &sendOTPRequest); err != nil {
		return
	}

	err := handler.authService.SendOTPToEmail(ctx, sendOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Verify OTP
// @Description Verify OTP with email and otp
// @Tags Auths
// @Accept json
// @Param request body model.VerifyOTPRequest true "Verify OTP payload"
// @Produce json
// @Router /auth/forgot-password/verify-otp [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) VerifyOTP(ctx *gin.Context) {
	var verifyOTPRequest model.VerifyOTPRequest

	if err := validation.BindJsonAndValidate(ctx, &verifyOTPRequest); err != nil {
		return
	}

	err := handler.authService.VerifyOTP(ctx, verifyOTPRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Set Password
// @Description Set a new password after OTP verification
// @Tags Auths
// @Accept json
// @Param request body model.SetPasswordRequest true "Set Password payload"
// @Produce json
// @Router /auth/forgot-password [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) SetPassword(ctx *gin.Context) {
	var setPasswordRequest model.SetPasswordRequest

	if err := validation.BindJsonAndValidate(ctx, &setPasswordRequest); err != nil {
		return
	}

	err := handler.authService.SetPassword(ctx, setPasswordRequest)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InvalidRequest,
			},
		))
		return
	}

	ctx.AbortWithStatus(204)
}

// @Summary Logout
// @Description Logout
// @Tags Auths
// @Router /auth/logout [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) Logout(ctx *gin.Context) {
	refreshToken, err := ctx.Cookie("refresh_token")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: "Refresh token required"},
		))
	}

	handler.authService.Logout(ctx, refreshToken)
	ctx.AbortWithStatus(204)
}

// @Summary Change Password
// @Description Change password for logged in users
// @Tags Auths
// @Accept json
// @Param request body model.ChangePasswordRequest true "Set Password payload"
// @Produce json
// @Router /auth/change-password [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) ChangePassword(ctx *gin.Context) {
	var request model.ChangePasswordRequest

	if err := validation.BindJsonAndValidate(ctx, &request); err != nil {
		return
	}

	err := handler.authService.ChangePassword(ctx, request)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error(), Code: httpcommon.ErrorResponseCode.InvalidRequest},
		))
		return
	}
	ctx.AbortWithStatus(204)
}

// @Summary Close account
// @Description Close account
// @Tags Auths
// @Accept json
// @Produce json
// @Router /auth/close [post]
// @Success 204 "No Content"
// @Failure 400 {object} httpcommon.HttpResponse[any]
// @Failure 500 {object} httpcommon.HttpResponse[any]
func (handler *AuthHandler) CloseAccount(ctx *gin.Context) {
	err := handler.authService.Close(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
			httpcommon.Error{Message: err.Error(), Code: httpcommon.ErrorResponseCode.InvalidRequest},
		))
		return
	}
	ctx.AbortWithStatus(204)
}
