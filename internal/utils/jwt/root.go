package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	Payload interface{}
	jwt.RegisteredClaims
}

func GenerateToken(duration time.Duration, secretKey string, payload interface{}) (string, error) {
	// Create a new JWT token
	claims := jwt.MapClaims{
		"exp":     time.Now().Add(duration).Unix(), // Set the expiration time
		"iat":     time.Now().Unix(),               // Set the issued at time
		"payload": payload,                         // Include the custom payload
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateTokenByClaims(claim TokenClaims, secretKey string) (string, error) {
	claims := jwt.MapClaims{
		"exp":     claim.ExpiresAt, // Set the expiration time
		"iat":     claim.IssuedAt,  // Set the issued at time
		"payload": claim.Payload,   // Include the custom payload
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func VerifyToken(tokenString, secretKey string) (*TokenClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
