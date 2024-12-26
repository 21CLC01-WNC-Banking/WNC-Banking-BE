package model

type RegisterRequest struct {
	Email       string `json:"email" binding:"required,email,min=10,max=255"`
	Name        string `json:"name" binding:"required,min=5,max=255"`
	PhoneNumber string `json:"phoneNumber" binding:"required,min=10,max=255"`
	Password    string `json:"password" binding:"required,min=8,max=255"`
}

type LoginRequest struct {
	Email          string `json:"email" binding:"required,email,min=10,max=255"`
	Password       string `json:"password" binding:"required,min=8,max=255"`
	RecaptchaToken string `json:"recaptchaToken" binding:"required"`
}

type AuthenticationRequest struct {
	CustomerID   int64  `json:"customerId" binding:"required,min=1,max=255"`
	RefreshToken string `json:"refreshToken" binding:"required"`
}

type SendOTPRequest struct {
	Email string `json:"email" binding:"required,email,min=10,max=255"`
}

type VerifyOTPRequest struct {
	Email string `json:"email" binding:"required,email,min=10,max=255"`
	OTP   string `json:"otp" binding:"required,min=6,max=6"`
}

type SetPasswordRequest struct {
	Email    string `json:"email" binding:"required,email,min=10,max=255"`
	OTP      string `json:"otp" binding:"required,min=6,max=6"`
	Password string `json:"password" binding:"required,min=8,max=255"`
}
