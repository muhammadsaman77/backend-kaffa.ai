package controllers

import "github.com/gin-gonic/gin"

type ProductController interface {
	CreateProduct(ctx *gin.Context)
	GetProduct(ctx *gin.Context)
	UpdateProduct(ctx *gin.Context)
	DeleteProduct(ctx *gin.Context)
	ListProducts(ctx *gin.Context)
}
