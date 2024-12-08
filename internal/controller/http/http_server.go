package http

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"

	"github.com/VuKhoa23/advanced-web-be/internal/controller/http/middleware"
	v1 "github.com/VuKhoa23/advanced-web-be/internal/controller/http/v1"
)

type Server struct {
	authHandler *v1.AuthHandler
	authMiddleware *middleware.AuthMiddleware
}

func NewServer(authHandler *v1.AuthHandler, authMiddleware *middleware.AuthMiddleware) *Server {
	return &Server{
		authHandler: authHandler,
		authMiddleware: authMiddleware,
	}
}

func (s *Server) Run() {
	router := gin.New()
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	httpServerInstance := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	v1.MapRoutes(router, s.authHandler, s.authMiddleware)
	err := httpServerInstance.ListenAndServe()
	if err != nil {
		return
	}
	fmt.Println("Server running at " + httpServerInstance.Addr)
}
