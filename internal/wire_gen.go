// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package internal

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean/implement"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/v1"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/websocket"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository/implement"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service/implement"
	"github.com/google/wire"
)

// Injectors from wire.go:

func InitializeContainer(db database.Db) *controller.ApiContainer {
	customerRepository := repositoryimplement.NewCustomerRepository(db)
	authenticationRepository := repositoryimplement.NewAuthenticationRepository(db)
	passwordEncoder := beanimplement.NewBcryptPasswordEncoder()
	redisClient := beanimplement.NewRedisService()
	accountRepository := repositoryimplement.NewAccountRepository(db)
	partnerBankRepository := repositoryimplement.NewPartnerBankRepository(db)
	partnerBankService := serviceimplement.NewPartnerBankService(partnerBankRepository)
	accountService := serviceimplement.NewAccountService(accountRepository, customerRepository, partnerBankService)
	mailClient := beanimplement.NewMailClient()
	roleRepository := repositoryimplement.NewRoleRepository(db)
	authService := serviceimplement.NewAuthService(customerRepository, authenticationRepository, passwordEncoder, redisClient, accountService, mailClient, roleRepository)
	authHandler := v1.NewAuthHandler(authService)
	coreService := serviceimplement.NewCoreService()
	server := websocket.NewServer()
	notificationRepository := repositoryimplement.NewNotificationRepository(db)
	notificationClient := beanimplement.NewNotificationClient(server, notificationRepository)
	coreHandler := v1.NewCoreHandler(coreService, notificationClient)
	savedReceiverRepository := repositoryimplement.NewSavedReceiverRepository(db)
	savedReceiverService := serviceimplement.NewSavedReceiverService(savedReceiverRepository, accountService, partnerBankService)
	accountHandler := v1.NewAccountHandler(accountService, savedReceiverService, authService)
	transactionRepository := repositoryimplement.NewTransactionRepository(db)
	staffService := serviceimplement.NewStaffService(customerRepository, passwordEncoder, accountService, accountRepository, transactionRepository, mailClient, notificationRepository, notificationClient)
	staffHandler := v1.NewStaffHandler(staffService)
	roleService := serviceimplement.NewRoleService(roleRepository)
	authMiddleware := middleware.NewAuthMiddleware(authService, roleService)
	debtReplyRepository := repositoryimplement.NewDebtReplyRepository(db)
	externalSearchMiddleware := middleware.NewExternalSearchMiddleware(partnerBankService)
	rsaMiddleware := middleware.NewRSAMiddleware(externalSearchMiddleware)
	transactionService := serviceimplement.NewTransactionService(transactionRepository, customerRepository, accountService, coreService, redisClient, mailClient, debtReplyRepository, notificationRepository, notificationClient, partnerBankService, rsaMiddleware)
	pgpMiddleware := middleware.NewPGPMiddleware(externalSearchMiddleware)
	transactionHandler := v1.NewTransactionHandler(transactionService, pgpMiddleware)
	savedReceiverHandler := v1.NewSavedReceiverHandler(savedReceiverService)
	notificationService := serviceimplement.NewNotificationService(notificationRepository)
	customerHandler := v1.NewCustomerHandler(notificationService, transactionService)
	staffRepository := repositoryimplement.NewStaffRepository(db)
	adminService := serviceimplement.NewAdminService(staffRepository, passwordEncoder, transactionRepository)
	adminHandler := v1.NewAdminHandler(adminService, partnerBankService)
	partnerBankHandler := v1.NewPartnerBankHandler(accountService, transactionService, partnerBankService, rsaMiddleware, pgpMiddleware)
	debtReplyService := serviceimplement.NewDebtReplyService(debtReplyRepository)
	debtReplyHandler := v1.NewDebtReplyHandler(debtReplyService)
	httpServer := http.NewServer(authHandler, coreHandler, accountHandler, staffHandler, authMiddleware, transactionHandler, savedReceiverHandler, customerHandler, adminHandler, partnerBankHandler, externalSearchMiddleware, rsaMiddleware, pgpMiddleware, debtReplyHandler)
	apiContainer := controller.NewApiContainer(httpServer, server)
	return apiContainer
}

// wire.go:

var container = wire.NewSet(controller.NewApiContainer)

// may have grpc server in the future
var serverSet = wire.NewSet(http.NewServer, websocket.NewServer)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(v1.NewAuthHandler, v1.NewCoreHandler, v1.NewAccountHandler, v1.NewStaffHandler, v1.NewTransactionHandler, v1.NewSavedReceiverHandler, v1.NewCustomerHandler, v1.NewAdminHandler, v1.NewPartnerBankHandler, v1.NewDebtReplyHandler)

var serviceSet = wire.NewSet(serviceimplement.NewAuthService, serviceimplement.NewAccountService, serviceimplement.NewCoreService, serviceimplement.NewRoleService, serviceimplement.NewTransactionService, serviceimplement.NewSavedReceiverService, serviceimplement.NewStaffService, serviceimplement.NewNotificationService, serviceimplement.NewAdminService, serviceimplement.NewPartnerBankService, serviceimplement.NewDebtReplyService)

var repositorySet = wire.NewSet(repositoryimplement.NewCustomerRepository, repositoryimplement.NewAuthenticationRepository, repositoryimplement.NewAccountRepository, repositoryimplement.NewRoleRepository, repositoryimplement.NewTransactionRepository, repositoryimplement.NewSavedReceiverRepository, repositoryimplement.NewNotificationRepository, repositoryimplement.NewDebtReplyRepository, repositoryimplement.NewStaffRepository, repositoryimplement.NewPartnerBankRepository)

var middlewareSet = wire.NewSet(middleware.NewAuthMiddleware, middleware.NewExternalSearchMiddleware, middleware.NewRSAMiddleware, middleware.NewPGPMiddleware)

var beanSet = wire.NewSet(beanimplement.NewBcryptPasswordEncoder, beanimplement.NewRedisService, beanimplement.NewMailClient, beanimplement.NewNotificationClient)
