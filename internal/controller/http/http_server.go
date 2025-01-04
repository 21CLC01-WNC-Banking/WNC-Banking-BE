package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/middleware"
	v1 "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/controller/http/v1"
)

type Server struct {
	authHandler              *v1.AuthHandler
	coreHandler              *v1.CoreHandler
	accountHandler           *v1.AccountHandler
	staffHandler             *v1.StaffHandler
	authMiddleware           *middleware.AuthMiddleware
	transactionHandler       *v1.TransactionHandler
	savedReceiverHandler     *v1.SavedReceiverHandler
	customerHandler          *v1.CustomerHandler
	adminHandler             *v1.AdminHandler
	partnerBankHandler       *v1.PartnerBankHandler
	externalSearchMiddleware *middleware.ExternalSearchMiddleware
	rsaMiddleware            *middleware.RSAMiddleware
	pgpMiddleware            *middleware.PGPMiddleware
}

func NewServer(authHandler *v1.AuthHandler,
	coreHandler *v1.CoreHandler,
	accountHandler *v1.AccountHandler,
	staffHandler *v1.StaffHandler,
	authMiddleware *middleware.AuthMiddleware,
	transactionHandler *v1.TransactionHandler,
	savedReceiverHandler *v1.SavedReceiverHandler,
	customerHandler *v1.CustomerHandler,
	adminHandler *v1.AdminHandler,
	partnerBankHandler *v1.PartnerBankHandler,
	externalSearchMiddleware *middleware.ExternalSearchMiddleware,
	rsaMiddleware *middleware.RSAMiddleware,
	pgpMiddleware *middleware.PGPMiddleware,
) *Server {
	return &Server{
		authHandler:              authHandler,
		authMiddleware:           authMiddleware,
		coreHandler:              coreHandler,
		accountHandler:           accountHandler,
		staffHandler:             staffHandler,
		transactionHandler:       transactionHandler,
		savedReceiverHandler:     savedReceiverHandler,
		customerHandler:          customerHandler,
		adminHandler:             adminHandler,
		partnerBankHandler:       partnerBankHandler,
		externalSearchMiddleware: externalSearchMiddleware,
		rsaMiddleware:            rsaMiddleware,
		pgpMiddleware:            pgpMiddleware,
	}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	v1.MapRoutes(router, s.authHandler, s.coreHandler, s.accountHandler, s.authMiddleware, s.staffHandler, s.transactionHandler, s.savedReceiverHandler, s.customerHandler, s.adminHandler, s.partnerBankHandler, s.externalSearchMiddleware, s.rsaMiddleware, s.pgpMiddleware)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		return
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)
}
