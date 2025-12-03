package controllers

import "github.com/gin-gonic/gin"

type AuthController interface {
	LoginUser() gin.HandlerFunc
	RegisterUser() gin.HandlerFunc
}
