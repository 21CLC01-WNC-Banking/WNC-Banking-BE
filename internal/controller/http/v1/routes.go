package v1

import (
	"github.com/VuKhoa23/advanced-web-be/internal/controller/http/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func MapRoutes(router *gin.Engine, authHandler *AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	router.Use(middleware.CorsMiddleware())
	v1 := router.Group("/api/v1")
	{
		customers := v1.Group("/auth")
		{
			customers.POST("/register", authHandler.Register)
			customers.POST("/login", authHandler.Login)
			customers.POST("/test", authMiddleware.VerifyToken, authHandler.TestJWT)
		}
	}
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
