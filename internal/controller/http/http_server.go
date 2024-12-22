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
	authHandler        *v1.AuthHandler
	coreHandler        *v1.CoreHandler
	accountHandler     *v1.AccountHandler
	transactionHandler *v1.TransactionHandler
	authMiddleware     *middleware.AuthMiddleware
}

func NewServer(authHandler *v1.AuthHandler, coreHandler *v1.CoreHandler, accountHandler *v1.AccountHandler, transactionHandler *v1.TransactionHandler, authMiddleware *middleware.AuthMiddleware) *Server {
	return &Server{
		authHandler:        authHandler,
		authMiddleware:     authMiddleware,
		coreHandler:        coreHandler,
		accountHandler:     accountHandler,
		transactionHandler: transactionHandler,
	}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	v1.MapRoutes(router, s.authHandler, s.coreHandler, s.accountHandler, s.transactionHandler, s.authMiddleware)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		return
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)
}
