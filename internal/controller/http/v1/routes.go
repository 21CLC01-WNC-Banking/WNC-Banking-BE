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
	customerHandler *CustomerHandler,
) {
	router.Use(middleware.CorsMiddleware())
	v1 := router.Group("/api/v1")
	{
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/forgot-password/otp", authHandler.SendOTPToMail)
			auth.POST("/forgot-password/verify-otp", authHandler.VerifyOTP)
			auth.POST("/forgot-password", authHandler.SetPassword)
		}
		customer := v1.Group("/customer")
		{
			customer.GET("/notification",
				authMiddleware.VerifyToken,
				customerHandler.GetNotifications,
			)
			customer.PATCH("/notification/:notificationId",
				authMiddleware.VerifyToken,
				customerHandler.SeenNotification,
			)
		}
		cores := v1.Group("/core")
		{
			cores.GET("/estimate-transfer-fee", coreHandler.EstimateTransferFee)
			cores.POST("/test-notification", coreHandler.Notification)
		}
		accounts := v1.Group("/account")
		{
			accounts.GET("/customer-name", authMiddleware.VerifyToken, accountHandler.GetCustomerNameByAccountNumber)
			accounts.GET("/", authMiddleware.VerifyToken, accountHandler.GetAccountByCustomerId)
		}
		staff := v1.Group("/staff")
		{
			staff.POST("/register-customer",
				authMiddleware.VerifyToken,
				authMiddleware.StaffRequired,
				staffHandler.RegisterCustomer,
			)
			staff.POST("/add-amount",
				authMiddleware.VerifyToken,
				authMiddleware.StaffRequired,
				staffHandler.AddAmountToAccount,
			)
			staff.GET("/transactions-by-account",
				authMiddleware.VerifyToken,
				authMiddleware.StaffRequired,
				staffHandler.GetTransactionsByAccountNumber,
			)
		}
		transactions := v1.Group("/transaction")
		{
			transactions.POST("/pre-internal-transfer", authMiddleware.VerifyToken, transactionHandler.PreInternalTransfer)
			transactions.POST("/internal-transfer", authMiddleware.VerifyToken, transactionHandler.InternalTransfer)
		}
		savedReceiver := v1.Group("/saved-receiver")
		{
			savedReceiver.POST("/", authMiddleware.VerifyToken, savedReceiverHandler.AddReceiver)
			savedReceiver.GET("/", authMiddleware.VerifyToken, savedReceiverHandler.GetAllReceivers)
			savedReceiver.PUT("/:id", authMiddleware.VerifyToken, savedReceiverHandler.RenameReceiver)
			savedReceiver.DELETE("/:id", authMiddleware.VerifyToken, savedReceiverHandler.DeleteReceiver)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
