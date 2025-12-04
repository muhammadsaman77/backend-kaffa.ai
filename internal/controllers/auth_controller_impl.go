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

func (c *AuthControllerImpl) LoginUser(ctx *gin.Context) {
	var loginRequest dto.LoginRequest
	if err := ctx.ShouldBindJSON(&loginRequest); err != nil {
		ctx.JSON(400, gin.H{
			"mesage": "Invalid request payload",
			"error":  pkg.ParseValidationErrors(err),
		})

		return
	}
	loginResponse, err := c.AuthService.LoginUser(ctx.Request.Context(), &loginRequest)
	if err != nil {
		if err.Error() == "USER_NOT_FOUND" || err.Error() == "INVALID_PASSWORD" {
			ctx.JSON(401, gin.H{
				"message": "Unauthorized",
				"error":   "Invalid username or password",
			})
			return

		}
		if err.Error() == "TOKEN_GENERATION_FAILED" {
			ctx.JSON(500, gin.H{
				"message": "Internal Server Error",
				"error":   "Failed to generate access token",
			})
			return
		}

	}
	ctx.JSON(200, gin.H{
		"message": "Login successful",
		"payload": loginResponse,
	})
}

func (c *AuthControllerImpl) RegisterUser(ctx *gin.Context) {
	var registerRequest dto.RegisterRequest
	if err := ctx.ShouldBindJSON(&registerRequest); err != nil {
		ctx.JSON(400, gin.H{
			"mesage": "Invalid request payload",
			"error":  pkg.ParseValidationErrors(err),
		})
		return
	}
	newUser, err := c.AuthService.RegisterUser(ctx.Request.Context(), &registerRequest)
	if err != nil {
		if err.Error() == "PASSWORD_HASHING_FAILED" {
			ctx.JSON(500, gin.H{
				"message": "Internal Server Error",
				"error":   "Failed to hash password",
			})
			return
		}
		if err.Error() == "USER_ALREADY_EXISTS" {
			ctx.JSON(409, gin.H{
				"message": "Conflict",
				"error":   "User with given email or username already exists",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   err,
		})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "User registered successfully",
		"payload": newUser,
	})

}
