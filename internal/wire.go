//go:build wireinject
// +build wireinject

package internal

import (
	beanimplement "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/bean/implement"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	v1 "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/v1"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/database"
	repositoryimplement "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/repository/implement"
	serviceimplement "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service/implement"
	"github.com/google/wire"
)

var container = wire.NewSet(
	controller.NewApiContainer,
)

// may have grpc server in the future
var serverSet = wire.NewSet(
	http.NewServer,
)

// handler === controller | with service and repository layers to form 3 layers architecture
var handlerSet = wire.NewSet(
	v1.NewAuthHandler,
	v1.NewCoreHandler,
	v1.NewAccountHandler,
	v1.NewTransactionHandler,
)

var serviceSet = wire.NewSet(
	serviceimplement.NewAuthService,
	serviceimplement.NewAccountService,
	serviceimplement.NewCoreService,
	serviceimplement.NewTransactionService,
)

var repositorySet = wire.NewSet(
	repositoryimplement.NewCustomerRepository,
	repositoryimplement.NewAuthenticationRepository,
	repositoryimplement.NewAccountRepository,
	repositoryimplement.NewTransactionRepository,
)

var middlewareSet = wire.NewSet(
	middleware.NewAuthMiddleware,
)
var beanSet = wire.NewSet(
	beanimplement.NewBcryptPasswordEncoder,
	beanimplement.NewRedisService,
	beanimplement.NewMailClient,
)

func InitializeContainer(
	db database.Db,
) *controller.ApiContainer {
	wire.Build(serverSet, handlerSet, serviceSet, repositorySet, middlewareSet, beanSet, container)
	return &controller.ApiContainer{}
}
