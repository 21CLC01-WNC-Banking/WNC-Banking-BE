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
	redisCLient := beanimplement.NewRedisService()
	accountRepository := repositoryimplement.NewAccountRepository(db)
	coreService := serviceimplement.NewCoreService()
	accountService := serviceimplement.NewAccountService(accountRepository, customerRepository, coreService)
	mailCLient := beanimplement.NewMailClient()
	authService := serviceimplement.NewAuthService(customerRepository, authenticationRepository, passwordEncoder, redisCLient, accountService, mailCLient)
	authHandler := v1.NewAuthHandler(authService)
	coreHandler := v1.NewCoreHandler(coreService)
	accountHandler := v1.NewAccountHandler(accountService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	server := http.NewServer(authHandler, coreHandler, accountHandler, authMiddleware)
	apiContainer := controller.NewApiContainer(server)
	return apiContainer
}

// wire.go:

var container = wire.NewSet(controller.NewApiContainer)

// may have grpc server in the future
var serverSet = wire.NewSet(http.NewServer)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(v1.NewAuthHandler, v1.NewCoreHandler, v1.NewAccountHandler)

var serviceSet = wire.NewSet(serviceimplement.NewAuthService, serviceimplement.NewAccountService, serviceimplement.NewCoreService)

var repositorySet = wire.NewSet(repositoryimplement.NewCustomerRepository, repositoryimplement.NewAuthenticationRepository, repositoryimplement.NewAccountRepository)

var middlewareSet = wire.NewSet(middleware.NewAuthMiddleware)

var beanSet = wire.NewSet(beanimplement.NewBcryptPasswordEncoder, beanimplement.NewRedisService, beanimplement.NewMailClient)
