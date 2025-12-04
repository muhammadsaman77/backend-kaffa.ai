package services

import (
	"context"

	"backend-kaffa.ai/internal/dto"
)

type ProductService interface {
	GetProductDetails(ctx context.Context, productID string) (*dto.CreateProductRequest, error)
	ListProducts(ctx context.Context) (*[]dto.CreateProductRequest, error)
	CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (string, error)
	UpdateProduct(ctx context.Context, productID string, product *dto.CreateProductRequest) error
	DeleteProduct(ctx context.Context, productID string) error
}
