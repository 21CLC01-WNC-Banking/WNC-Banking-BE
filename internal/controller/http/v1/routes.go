package v1

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(router *gin.Engine, authHandler *AuthHandler, coreHandler *CoreHandler, accountHandler *AccountHandler, authMiddleware *middleware.AuthMiddleware) {
	router.Use(middleware.CorsMiddleware())
	v1 := router.Group("/api/v1")
	{
		customers := v1.Group("/auth")
		{
			customers.POST("/register", authHandler.Register)
			customers.POST("/login", authHandler.Login)
			customers.POST("/forgot-password/otp", authHandler.SendOTPToMail)
			customers.POST("/forgot-password/verify-otp", authHandler.VerifyOTP)
			customers.POST("/forgot-password", authHandler.SetPassword)
		}
		cores := v1.Group("/core")
		{
			cores.GET("/estimate-transfer-fee", coreHandler.EstimateTransferFee)
		}
		accounts := v1.Group("/account")
		{
			accounts.POST("/internal-transfer", authMiddleware.VerifyToken, accountHandler.InternalTransfer)
			accounts.GET("/customer-name", authMiddleware.VerifyToken, accountHandler.GetCustomerNameByAccountNumber)
			accounts.POST("/add-internal-receiver", authMiddleware.VerifyToken, accountHandler.AddInternalReceiver)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
