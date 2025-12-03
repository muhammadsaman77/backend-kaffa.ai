package middlewares

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
			"error":   "Missing access token",
		})
		c.Abort()
		return
	}
	accessToken := authHeader[len("Bearer "):]
	if accessToken == "" {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
			"error":   "Invalid access token",
		})
		c.Abort()
		return
	}
	keyData, _ := os.ReadFile("keys/public.pem")
	pubKey, _ := jwt.ParseRSAPublicKeyFromPEM(keyData)
	_, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		return pubKey, nil
	})
	if err != nil {
		c.JSON(401, gin.H{
			"message": "Unauthorized",
			"error":   "Invalid access token",
		})
		c.Abort()
		return
	}
	c.Next()
}
