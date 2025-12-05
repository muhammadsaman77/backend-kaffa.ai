package services

import (
	"context"
	"mime/multipart"

	"backend-kaffa.ai/internal/dto"
)

type ProductService interface {
	GetProductDetails(ctx context.Context, productID string) (*dto.GetProductDetailsResponse, error)
	ListProducts(ctx context.Context, storeId string) ([]dto.GetAllProductResponse, error)
	CreateProduct(ctx context.Context, product *dto.CreateProductRequest, imageHeader *multipart.FileHeader) (string, error)
	UpdateProduct(ctx context.Context, productID string, product *dto.CreateProductRequest) error
	DeleteProduct(ctx context.Context, productID string) error
}
