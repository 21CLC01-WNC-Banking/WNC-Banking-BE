package v1

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(router *gin.Engine,
	authHandler *AuthHandler,
	coreHandler *CoreHandler,
	accountHandler *AccountHandler,
	authMiddleware *middleware.AuthMiddleware,
	staffHandler *StaffHandler,
	transactionHandler *TransactionHandler,
	savedReceiverHandler *SavedReceiverHandler,
) {
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
			accounts.GET("/customer-name", authMiddleware.VerifyToken, accountHandler.GetCustomerNameByAccountNumber)
		}
		staff := v1.Group("/staff")
		{
			staff.POST("/register-customer",
				authMiddleware.VerifyToken,
				authMiddleware.StaffRequired,
				staffHandler.RegisterCustomer,
			)
		}
		transactions := v1.Group("/transaction")
		{
			transactions.POST("/pre-internal-transfer", authMiddleware.VerifyToken, transactionHandler.PreInternalTransfer)
			transactions.POST("/internal-transfer", authMiddleware.VerifyToken, transactionHandler.InternalTransfer)
		}
		savedReceiver := v1.Group("/saved-receiver")
		{
			savedReceiver.POST("/add-internal", authMiddleware.VerifyToken, savedReceiverHandler.AddReceiver)
			savedReceiver.GET("/", authMiddleware.VerifyToken, savedReceiverHandler.GetAllReceivers)
			savedReceiver.PUT("/:id", authMiddleware.VerifyToken, savedReceiverHandler.RenameReceiver)
			savedReceiver.DELETE("/:id", authMiddleware.VerifyToken, savedReceiverHandler.DeleteReceiver)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
