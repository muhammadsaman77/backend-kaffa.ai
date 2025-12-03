package controllers

import "github.com/gin-gonic/gin"

type AuthController interface {
	LoginUser(c *gin.Context)
	RegisterUser(c *gin.Context)
}
