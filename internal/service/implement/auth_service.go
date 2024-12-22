package serviceimplement

import (
	"database/sql"
	"errors"
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
	redisCLient              bean.RedisClient
	mailCLient               bean.MailClient
}

func NewAuthService(customerRepository repository.CustomerRepository,
	authenticationRepository repository.AuthenticationRepository,
	encoder bean.PasswordEncoder,
	redisCLient bean.RedisClient,
	accountSer service.AccountService,
	mailCLient bean.MailClient,
) service.AuthService {
	return &AuthService{
		customerRepository:       customerRepository,
		authenticationRepository: authenticationRepository,
		passwordEncoder:          encoder,
		redisCLient:              redisCLient,
		accountService:           accountSer,
		mailCLient:               mailCLient,
	}
}

func (service *AuthService) Register(ctx *gin.Context, registerRequest model.RegisterRequest) error {
	existsCustomer, err := service.customerRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil {
		return err
	}
	if existsCustomer != nil {
		return errors.New("Email have already registered")
	}
	hashPW, err := service.passwordEncoder.Encrypt(registerRequest.Password)
	if err != nil {
		return err
	}
	newCustomer := &entity.Customer{
		Email:       registerRequest.Email,
		Name:        registerRequest.Name,
		PhoneNumber: registerRequest.PhoneNumber,
		Password:    string(hashPW),
		RoleId:      1,
	}
	err = service.customerRepository.CreateCommand(ctx, newCustomer)
	if err != nil {
		return err
	}

	// auto create an account
	currentCustomer, err := service.customerRepository.GetOneByEmailQuery(ctx, registerRequest.Email)
	if err != nil {
		return err
	}
	err = service.accountService.AddNewAccount(ctx, currentCustomer.ID)
	if err != nil {
		return err
	}
	return nil
}

func (service *AuthService) Login(ctx *gin.Context, loginRequest model.LoginRequest) (*entity.Customer, error) {
	// validate captcha
	isValid, err := google_recaptcha.ValidateRecaptcha(ctx, loginRequest.RecaptchaToken)
	if err != nil || !isValid {
		return nil, err
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
		"id": existsCustomer.ID,
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
		"id": existsCustomer.ID,
	})
	if err != nil {
		return nil, err
	}

	// Check if a refresh token already exists
	existingRefreshToken, err := service.authenticationRepository.GetOneByCustomerIdQuery(ctx, existsCustomer.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	if existingRefreshToken == nil {
		// Create a new refresh token
		err = service.authenticationRepository.CreateCommand(ctx, entity.Authentication{
			UserId:       existsCustomer.ID,
			RefreshToken: refreshToken,
		})
		if err != nil {
			return nil, err
		}
	} else {
		// Update the existing refresh token
		err = service.authenticationRepository.UpdateCommand(ctx, entity.Authentication{
			UserId:       existsCustomer.ID,
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

	return existsCustomer, nil
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

	err = service.redisCLient.Set(ctx, key, otp)
	if err != nil {
		return err
	}

	// send otp to user email
	err = service.mailCLient.SendEmail(ctx, sendOTPRequest.Email, "OTP reset password", otp, constants.FORGOT_PASSWORD, constants.RESET_PASSWORD_EXP_TIME)
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

	val, err := service.redisCLient.Get(ctx, key)
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

	val, err := service.redisCLient.Get(ctx, key)
	if err != nil {
		return err
	}

	if val == setPasswordRequest.OTP {
		service.redisCLient.Delete(ctx, key)

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
