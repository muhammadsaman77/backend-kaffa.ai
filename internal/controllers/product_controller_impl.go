package controllers

import (
	"errors"

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
	if err := ctx.ShouldBind(&createProductRequest); err != nil {
		ctx.JSON(400, gin.H{
			"message": "Invalid request payload",
			"error":   pkg.ParseValidationErrors(err),
		})
		return
	}
	imageHeader, err := ctx.FormFile("image")
	if err != nil {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
			"error":   "Image file is required",
		})
		return
	}

	createdProduct, err := c.ProductService.CreateProduct(ctx.Request.Context(), &createProductRequest, imageHeader)
	if errors.Is(err, pkg.ErrInvalidPrice) {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
			"error":   "Price must be a valid number",
		})
		return
	}
	if errors.Is(err, pkg.ErrStoreNotFound) {
		ctx.JSON(404, gin.H{
			"message": "Store not found",
			"error":   "The specified store does not exist",
		})
		return
	}
	if errors.Is(err, pkg.ErrInvalidImageMimeType) {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
			"error":   "Unsupported image format",
		})
		return
	}
	if errors.Is(err, pkg.ErrImageSizeExceedsLimit) {
		ctx.JSON(400, gin.H{
			"message": "Bad Request",
			"error":   "Image size exceeds the allowed limit",
		})
		return
	}
	if errors.Is(err, pkg.ErrFailedToUploadImage) {
		ctx.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   "Failed to upload image",
		})
		return
	}
	if errors.Is(err, pkg.ErrFailedToCreateProduct) {
		ctx.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   "Failed to create product",
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

	productId := ctx.Param("id")

	err := c.ProductService.DeleteProduct(ctx.Request.Context(), productId)
	if errors.Is(err, pkg.ErrProductNotFound) {
		ctx.JSON(404, gin.H{
			"message": "Product not found",
			"error":   "The specified product does not exist",
		})
		return
	}
	if errors.Is(err, pkg.ErrFailedToDeleteProduct) {
		ctx.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   "Failed to delete product",
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": "Product deleted successfully",
	})

}

func (c *ProductControllerImpl) ListProducts(ctx *gin.Context) {
	storeId := ctx.Query("store_id")

	productsList, err := c.ProductService.ListProducts(ctx.Request.Context(), storeId)
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": "Internal Server Error",
			"error":   "Failed to list products",
		})
		return
	}
	ctx.JSON(200, gin.H{
		"message": "Products retrieved successfully",
		"payload": productsList,
	})
}
