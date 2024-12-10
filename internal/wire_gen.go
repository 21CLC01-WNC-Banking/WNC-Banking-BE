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
	authService := serviceimplement.NewAuthService(customerRepository, authenticationRepository, passwordEncoder)
	authHandler := v1.NewAuthHandler(authService)
	authMiddleware := middleware.NewAuthMiddleware(authService)
	server := http.NewServer(authHandler, authMiddleware)
	apiContainer := controller.NewApiContainer(server)
	return apiContainer
}

// wire.go:

var container = wire.NewSet(controller.NewApiContainer)

// may have grpc server in the future
var serverSet = wire.NewSet(http.NewServer)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(v1.NewAuthHandler)

var serviceSet = wire.NewSet(serviceimplement.NewAuthService)

var repositorySet = wire.NewSet(repositoryimplement.NewCustomerRepository, repositoryimplement.NewAuthenticationRepository)

var middlewareSet = wire.NewSet(middleware.NewAuthMiddleware)

var beanSet = wire.NewSet(beanimplement.NewBcryptPasswordEncoder)
