package controllers

import (
	"backend-kaffa.ai/internal/dto"
	"backend-kaffa.ai/internal/services"
	"backend-kaffa.ai/pkg"
	"github.com/gin-gonic/gin"
)

type ProductControllerImpl struct {
	// Add necessary services here, e.g., ProductService
	ProductService services.ProductService
}

func NewProductController(productService services.ProductService) ProductController {
	return &ProductControllerImpl{
		ProductService: productService,
	}
}

func (c *ProductControllerImpl) CreateProduct(ctx *gin.Context) {
	var createProductRequest dto.CreateProductRequest
	if err := ctx.ShouldBindJSON(&createProductRequest); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request payload",
			"error":   pkg.ParseValidationErrors(err),
		})
		return
	}

	createdProduct, err := c.ProductService.CreateProduct(ctx.Request.Context(), &createProductRequest)
	if err != nil {
		if err.Error() == "STORE_NOT_FOUND" {
			ctx.JSON(404, gin.H{
				"message": "Store not found",
				"error":   "The specified store does not exist",
			})
			return
		}
		ctx.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   "Failed to create product",
			"details": err,
		})
		return
	}
	ctx.JSON(201, gin.H{
		"message": "Product created successfully",
		"payload": createdProduct,
	})
}

func (c *ProductControllerImpl) GetProduct(ctx *gin.Context) {
	// Implementation for getting a product
}

func (c *ProductControllerImpl) UpdateProduct(ctx *gin.Context) {
	// Implementation for updating a product
}

func (c *ProductControllerImpl) DeleteProduct(ctx *gin.Context) {
	// Implementation for deleting a product
}

func (c *ProductControllerImpl) ListProducts(ctx *gin.Context) {
	// Implementation for listing products
}
