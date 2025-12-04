package services

import (
	"context"
	"errors"

	"backend-kaffa.ai/configs"
	"backend-kaffa.ai/internal/dto"
	"backend-kaffa.ai/internal/sqlc/products"
	"backend-kaffa.ai/pkg"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

type ProductServiceImpl struct {
	productsQueries *products.Queries
}

func NewProductService(productsQueries *products.Queries) ProductService {
	return &ProductServiceImpl{
		productsQueries: productsQueries,
	}
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, product *dto.CreateProductRequest) (string, error) {
	price, err := pkg.Float64ToNumeric(product.Price)
	if err != nil {
		configs.Log.Error("failed to convert price to numeric", zap.Error(err))
		return "", errors.New("INVALID_PRICE")
	}
	newProduct, err := s.productsQueries.CreateProduct(ctx, products.CreateProductParams{
		ID:          ulid.Make().String(),
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: pgtype.Text{String: product.Description, Valid: true},
		Price:       price,
		IsAvailable: pgtype.Bool{Bool: product.IsAvailable, Valid: true},
	})
	if err != nil {
		if pgErrm, ok := err.(*pgconn.PgError); ok {
			if pgErrm.Code == "23503" {
				configs.Log.Error("store not found", zap.String("store_id", product.StoreID))

				return "", errors.New("STORE_NOT_FOUND")
			}
		}
		configs.Log.Error("failed to create product", zap.Error(err))
		return "", errors.New("FAILED_TO_CREATE_PRODUCT")
	}
	configs.Log.Info("product created successfully", zap.String("product_id", newProduct.ID))
	return newProduct.ID, nil
}
func (s *ProductServiceImpl) GetProductDetails(ctx context.Context, productID string) (*dto.CreateProductRequest, error) {
	return nil, nil
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context) (*[]dto.CreateProductRequest, error) {
	return nil, nil
}
func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, productID string, product *dto.CreateProductRequest) error {
	return nil
}
func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, productID string) error {
	return nil
}
