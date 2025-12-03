package controllers

import (
	"backend-kaffa.ai/internal/dto"
	"backend-kaffa.ai/internal/services"
	"backend-kaffa.ai/pkg"
	"github.com/gin-gonic/gin"
)

type AuthControllerImpl struct {
	AuthService services.AuthService
}

func NewAuthController(authService services.AuthService) AuthController {
	return &AuthControllerImpl{
		AuthService: authService,
	}
}

func (controller *AuthControllerImpl) LoginUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		var loginRequest dto.LoginRequest
		if err := c.ShouldBindJSON(&loginRequest); err != nil {
			c.JSON(400, gin.H{
				"mesage": "Invalid request payload",
				"error":  pkg.ParseValidationErrors(err),
			})
			return
		}
		loginResponse, err := controller.AuthService.LoginUser(c.Request.Context(), &loginRequest)
		if err != nil {
			if err.Error() == "USER_NOT_FOUND" || err.Error() == "INVALID_PASSWORD" {
				c.JSON(401, gin.H{
					"message": "Unauthorized",
					"error":   "Invalid username or password",
				})
				return
			}
			if err.Error() == "TOKEN_GENERATION_FAILED" {
				c.JSON(500, gin.H{
					"message": "Internal Server Error",
					"error":   "Failed to generate access token",
				})
				return
			}
			return
		}
		c.JSON(200, gin.H{
			"message": "Login successful",
			"payload": loginResponse,
		})
	}
}

func (a *AuthControllerImpl) RegisterUser() gin.HandlerFunc {
	return func(c *gin.Context) {}
}
