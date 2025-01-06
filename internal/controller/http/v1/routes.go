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
	adminHandler *AdminHandler,
	partnerBankHandler *PartnerBankHandler,
	externalSearchMiddleware *middleware.ExternalSearchMiddleware,
	rsaMiddleware *middleware.RSAMiddleware,
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
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/change-password", authMiddleware.VerifyToken, authHandler.ChangePassword)
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
			customer.GET("/transaction",
				authMiddleware.VerifyToken,
				customerHandler.GetTransactions,
			)
			customer.GET("/transaction/:transactionId",
				authMiddleware.VerifyToken,
				customerHandler.GetTransactionById,
			)
			savedReceiver := customer.Group("/saved-receiver")
			{
				savedReceiver.POST("/", authMiddleware.VerifyToken, savedReceiverHandler.AddReceiver)
				savedReceiver.GET("/", authMiddleware.VerifyToken, savedReceiverHandler.GetAllReceivers)
				savedReceiver.PUT("/:id", authMiddleware.VerifyToken, savedReceiverHandler.RenameReceiver)
				savedReceiver.DELETE("/:id", authMiddleware.VerifyToken, savedReceiverHandler.DeleteReceiver)
			}
		}
		admin := v1.Group("/admin")
		{
			admin.GET("/staff",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.GetAllStaff,
			)
			admin.GET("/staff/:staffId",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.GetOneStaff,
			)
			admin.POST("/staff",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.CreateOneStaff,
			)
			admin.DELETE("/staff/:staffId",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.DeleteOneStaff,
			)
			admin.PUT("/staff",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.UpdateOneStaff,
			)
			admin.POST("/partner-bank",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.AddPartnerBank)
			admin.GET("/external-transaction",
				authMiddleware.VerifyToken,
				authMiddleware.AdminRequired,
				adminHandler.GetExternalTransactions)
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
			accounts.POST("/get-external-account-name", authMiddleware.VerifyToken, accountHandler.GetExternalAccountName)
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
			transactions.POST("/debt-reminder", authMiddleware.VerifyToken, transactionHandler.AddDebtReminder)
			transactions.PUT("/cancel-debt-reminder/:id", authMiddleware.VerifyToken, transactionHandler.CancelDebtReminder)
			transactions.GET("/received-debt-reminder", authMiddleware.VerifyToken, transactionHandler.GetReceivedDebtReminder)
			transactions.GET("/sent-debt-reminder", authMiddleware.VerifyToken, transactionHandler.GetSentDebtReminder)
			transactions.POST("/pre-debt-transfer", authMiddleware.VerifyToken, transactionHandler.PreDebtTransfer)
			transactions.POST("/pre-external-transfer", authMiddleware.VerifyToken, transactionHandler.PreExternalTransfer)
			transactions.POST("/external-transfer", authMiddleware.VerifyToken, transactionHandler.ExternalTransfer)
		}
		partnerBanks := v1.Group("/partner-bank")
		{
			partnerBanks.POST("/get-account-information", externalSearchMiddleware.VerifyAPI, partnerBankHandler.GetAccountNumberInfo)
			partnerBanks.POST("/external-transfer-rsa", rsaMiddleware.Verify, partnerBankHandler.ReceiveExternalTransfer)
			partnerBanks.GET("/", authMiddleware.VerifyToken, partnerBankHandler.GetListPartnerBank)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
