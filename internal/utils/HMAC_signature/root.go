package HMAC_signature

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
)

func calculateHMAC(data, secretKey string) []byte {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(data))
	return h.Sum(nil)
}

func GenerateHMAC(data, secretKey string) string {
	signature := calculateHMAC(data, secretKey)
	fmt.Println("data: " + data)
	fmt.Println("expect hashed: " + base64.StdEncoding.EncodeToString(signature))
	return base64.StdEncoding.EncodeToString(signature)
}

func VerifyHMAC(data, signature, secretKey string) bool {
	expectedSignature := calculateHMAC(data, secretKey)
	fmt.Println("data: " + data)
	fmt.Println("input hashed: " + signature)
	fmt.Println("expect hashed: " + base64.StdEncoding.EncodeToString(expectedSignature))
	actualSignature, err := base64.StdEncoding.DecodeString(signature)
	if err != nil {
		fmt.Println("Failed to decode signature:", err)
		return false
	}

	return hmac.Equal(expectedSignature, actualSignature)
}
