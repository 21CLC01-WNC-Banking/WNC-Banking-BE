package serviceimplement

import (
	"database/sql"
	"errors"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/entity"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/model"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/constants"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/google_recaptcha"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/jwt"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/mail"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/redis"
	"github.com/gin-gonic/gin"
)

type AuthService struct {
	customerRepository       repository.CustomerRepository
	authenticationRepository repository.AuthenticationRepository
	passwordEncoder          bean.PasswordEncoder
	accountService           service.AccountService
	redisClient              bean.RedisClient
	mailClient               bean.MailClient
	roleRepository           repository.RoleRepository
}

func NewAuthService(customerRepository repository.CustomerRepository,
	authenticationRepository repository.AuthenticationRepository,
	encoder bean.PasswordEncoder,
	redisClient bean.RedisClient,
	accountSer service.AccountService,
	mailClient bean.MailClient,
	roleRepository repository.RoleRepository,
) service.AuthService {
	return &AuthService{
		customerRepository:       customerRepository,
		authenticationRepository: authenticationRepository,
		passwordEncoder:          encoder,
		redisClient:              redisClient,
		accountService:           accountSer,
		mailClient:               mailClient,
		roleRepository:           roleRepository,
	}
}

func (service *AuthService) Login(ctx *gin.Context, loginRequest model.LoginRequest) (*model.LoginResponse, error) {
	// validate captcha
	isValid, err := google_recaptcha.ValidateRecaptcha(ctx, loginRequest.RecaptchaToken)
	if err != nil || !isValid {
		return nil, errors.New("invalid recaptcha token")
	}

	existsCustomer, err := service.customerRepository.GetOneByEmailQuery(ctx, loginRequest.Email)
	if err != nil {
		if err.Error() == httpcommon.ErrorMessage.SqlxNoRow {
			return nil, errors.New("Email not found")
		}
		return nil, err
	}
	checkPw := service.passwordEncoder.Compare(existsCustomer.Password, loginRequest.Password)
	if checkPw == false {
		return nil, errors.New("invalid password")
	}

	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		return nil, err
	}
	accessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": existsCustomer.Id,
	})

	if err == nil {
		ctx.SetCookie(
			"access_token",
			accessToken,
			constants.COOKIE_DURATION,
			"/",
			"",
			false,
			true,
		)
	}

	refreshToken, err := jwt.GenerateToken(constants.REFRESH_TOKEN_DURATION, jwtSecret, map[string]interface{}{
		"id": existsCustomer.Id,
	})
	if err != nil {
		return nil, err
	}

	// Check if a refresh token already exists
	existingRefreshToken, err := service.authenticationRepository.GetOneByCustomerIdQuery(ctx, existsCustomer.Id)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingRefreshToken == nil {
		// Create a new refresh token
		err = service.authenticationRepository.CreateCommand(ctx, entity.Authentication{
			UserId:       existsCustomer.Id,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Update the existing refresh token
		err = service.authenticationRepository.UpdateCommand(ctx, entity.Authentication{
			UserId:       existsCustomer.Id,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
	}

	// Set refresh token as a cookie
	ctx.SetCookie(
		"refresh_token",
		refreshToken,
		constants.COOKIE_DURATION,
		"/",
		"",
		false,
		true,
	)

	role, err := service.roleRepository.GetByUserId(ctx, existsCustomer.Id)
	if err != nil {
		return nil, err
	}

	return &model.LoginResponse{
		Name:         existsCustomer.Name,
		Email:        existsCustomer.Email,
		UserId:       existsCustomer.Id,
		Role:         role.Name,
		RefreshToken: refreshToken,
	}, nil
}

func (service *AuthService) ValidateRefreshToken(ctx *gin.Context, customerId int64) (*entity.Authentication, error) {
	refreshToken, err := service.authenticationRepository.GetOneByCustomerIdQuery(ctx, customerId)
	if err != nil {
		return nil, err
	}
	return refreshToken, nil
}

func (service *AuthService) SendOTPToEmail(ctx *gin.Context, sendOTPRequest model.SendOTPRequest) error {
	// generate otp
	otp := mail.GenerateOTP(6)

	// store otp in redis
	customerId, err := service.customerRepository.GetIdByEmailQuery(ctx, sendOTPRequest.Email)
	if err != nil {
		return err
	}
	baseKey := constants.RESET_PASSWORD_KEY
	key := redis.Concat(baseKey, customerId)

	err = service.redisClient.Set(ctx, key, otp)
	if err != nil {
		return err
	}

	// send otp to user email
	emailBody := service.mailClient.GenerateOTPBody(sendOTPRequest.Email, otp, constants.FORGOT_PASSWORD, constants.RESET_PASSWORD_EXP_TIME)
	err = service.mailClient.SendEmail(ctx, sendOTPRequest.Email, "OTP reset password", emailBody)
	if err != nil {
		return err
	}

	return nil
}

func (service *AuthService) VerifyOTP(ctx *gin.Context, verifyOTPRequest model.VerifyOTPRequest) error {
	customerId, err := service.customerRepository.GetIdByEmailQuery(ctx, verifyOTPRequest.Email)
	if err != nil {
		return err
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redis.Concat(baseKey, customerId)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	if val != verifyOTPRequest.OTP {
		return errors.New("Invalid OTP")
	}

	return nil
}

func (service *AuthService) SetPassword(ctx *gin.Context, setPasswordRequest model.SetPasswordRequest) error {
	customerId, err := service.customerRepository.GetIdByEmailQuery(ctx, setPasswordRequest.Email)
	if err != nil {
		return err
	}

	baseKey := constants.RESET_PASSWORD_KEY
	key := redis.Concat(baseKey, customerId)

	val, err := service.redisClient.Get(ctx, key)
	if err != nil {
		return err
	}

	if val == setPasswordRequest.OTP {
		service.redisClient.Delete(ctx, key)

		hashedPW, err := service.passwordEncoder.Encrypt(setPasswordRequest.Password)
		if err != nil {
			return err
		}

		err = service.customerRepository.UpdatePasswordByIdQuery(ctx, customerId, hashedPW)
		if err != nil {
			return err
		}
	} else {
		return errors.New("Invalid OTP")
	}

	return nil
}

func (service *AuthService) GetUserById(ctx *gin.Context, userId int64) (*entity.User, error) {
	return service.customerRepository.GetOneByIdQuery(ctx, userId)
}

func (service *AuthService) Logout(ctx *gin.Context, refreshToken string) {
	_ = service.authenticationRepository.DeleteByRefreshToken(ctx, refreshToken)
	ctx.SetCookie(
		"access_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
	ctx.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		false,
		true,
	)
}

func (service *AuthService) ChangePassword(ctx *gin.Context, request model.ChangePasswordRequest) error {
	userId := middleware.GetUserIdHelper(ctx)

	user, err := service.GetUserById(ctx, userId)
	if err != nil {
		return err
	}

	if !service.passwordEncoder.Compare(user.Password, request.Password) {
		return errors.New("invalid password")
	}

	newPasswordHashed, err := service.passwordEncoder.Encrypt(request.NewPassword)
	if err != nil {
		return err
	}
	err = service.customerRepository.UpdatePasswordByIdQuery(ctx, userId, newPasswordHashed)
	return err
}

func (service *AuthService) Close(ctx *gin.Context) error {
	userId := middleware.GetUserIdHelper(ctx)

	return service.customerRepository.DeleteById(ctx, userId)
}
