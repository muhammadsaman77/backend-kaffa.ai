package controllers

import "github.com/gin-gonic/gin"

type AuthController interface {
	LoginUser(ctx *gin.Context)
	RegisterUser(ctx *gin.Context)
}
