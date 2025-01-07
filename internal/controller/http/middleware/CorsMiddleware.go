package middleware

import (
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"
	"github.com/gin-gonic/gin"
)

func CorsMiddleware() gin.HandlerFunc {
	url, err := env.GetEnv("ALLOWED_ORIGIN")
	if err != nil {
		panic(err)
	}
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", url)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
			return
		}
		c.Next()
	}
}
