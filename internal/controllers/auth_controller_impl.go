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

func (controller *AuthControllerImpl) LoginUser(c *gin.Context) {
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

	}
	c.JSON(200, gin.H{
		"message": "Login successful",
		"payload": loginResponse,
	})
}

func (controller *AuthControllerImpl) RegisterUser(c *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(400, gin.H{
			"mesage": "Invalid request payload",
			"error":  pkg.ParseValidationErrors(err),
		})
		return
	}
	newUser, err := controller.AuthService.RegisterUser(c.Request.Context(), &registerRequest)
	if err != nil {
		if err.Error() == "PASSWORD_HASHING_FAILED" {
			c.JSON(500, gin.H{
				"message": "Internal Server Error",
				"error":   "Failed to hash password",
			})
			return
		}
		if err.Error() == "USER_ALREADY_EXISTS" {
			c.JSON(409, gin.H{
				"message": "Conflict",
				"error":   "User with given email or username already exists",
			})
			return
		}
		c.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   err,
		})
		return
	}
	c.JSON(201, gin.H{
		"message": "User registered successfully",
		"payload": newUser,
	})

}
