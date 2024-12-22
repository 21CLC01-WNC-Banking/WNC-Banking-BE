package middleware

import (
	"net/http"

	httpcommon "github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/domain/http_common"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/service"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/constants"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/env"
	"github.com/21CLC01-WNC-Banking/WNC-Banking-BE/internal/utils/jwt"
	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	authService service.AuthService
}

func NewAuthMiddleware(authService service.AuthService) *AuthMiddleware {
	return &AuthMiddleware{
		authService: authService,
	}
}

func getAccessToken(c *gin.Context) (token string) {
	token, err := c.Cookie("access_token")
	if err != nil {
		return ""
	}
	return token
}

func getRefreshToken(c *gin.Context) (token string) {
	token, err := c.Cookie("refresh_token")
	if err != nil {
		return ""
	}
	return token
}

func (a *AuthMiddleware) VerifyToken(c *gin.Context) {
	// Get the JWT secret from the environment
	jwtSecret, err := env.GetEnv("JWT_SECRET")
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
			httpcommon.Error{
				Message: err.Error(),
				Code:    httpcommon.ErrorResponseCode.InternalServerError,
			},
		))
		return
	}

	// Retrieve the access token from the header or cookies
	accessToken := getAccessToken(c)
	claims, err := jwt.VerifyToken(accessToken, jwtSecret)
	if err == nil {
		// If the access token is valid, extract customer ID and proceed
		if payload, ok := claims.Payload.(map[string]interface{}); ok {
			customerId := int64(payload["id"].(float64))
			c.Set("customerId", customerId)
			c.Next()
			return
		}
	}

	// If the access token is expired, check the refresh token
	if err.Error() == httpcommon.ErrorMessage.TokenExpired {
		refreshToken := getRefreshToken(c)
		refreshClaims, errRf := jwt.VerifyToken(refreshToken, jwtSecret)
		if errRf != nil {
			// If the refresh token is invalid or expired, abort with unauthorized
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: httpcommon.ErrorMessage.BadCredential,
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				},
			))
			return
		}

		// Extract customer ID from refresh token claims
		payload, ok := refreshClaims.Payload.(map[string]interface{})
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: httpcommon.ErrorMessage.BadCredential,
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				},
			))
			return
		}
		customerId := int64(payload["id"].(float64))

		// Check if the refresh token exists and is still valid in the database
		refreshTokenEntity, err := a.authService.ValidateRefreshToken(c, customerId)
		if err != nil || refreshTokenEntity == nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: httpcommon.ErrorMessage.BadCredential,
					Code:    httpcommon.ErrorResponseCode.Unauthorized,
				},
			))
			return
		}

		// Generate a new access token
		newAccessToken, err := jwt.GenerateToken(constants.ACCESS_TOKEN_DURATION, jwtSecret, map[string]interface{}{
			"id": customerId,
		})
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, httpcommon.NewErrorResponse(
				httpcommon.Error{
					Message: err.Error(),
					Code:    httpcommon.ErrorResponseCode.InternalServerError,
				},
			))
			return
		}

		// Set the new access token in cookies
		c.SetCookie(
			"access_token",
			newAccessToken,
			constants.COOKIE_DURATION,
			"/",
			"",
			false,
			true,
		)

		// Proceed with the customer ID set in the context
		c.Set("customerId", customerId)
		c.Next()
		return
	}

	// For all other errors, abort with unauthorized
	c.AbortWithStatusJSON(http.StatusUnauthorized, httpcommon.NewErrorResponse(
		httpcommon.Error{
			Message: httpcommon.ErrorMessage.BadCredential,
			Code:    httpcommon.ErrorResponseCode.Unauthorized,
		},
	))
}
