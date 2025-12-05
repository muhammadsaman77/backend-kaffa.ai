package services

import (
	"context"
	"errors"
	"mime/multipart"

	"backend-kaffa.ai/configs"
	"backend-kaffa.ai/internal/dto"
	"backend-kaffa.ai/internal/sqlc/products"
	"backend-kaffa.ai/pkg"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/oklog/ulid/v2"
	"go.uber.org/zap"
)

type ProductServiceImpl struct {
	ProductsQueries *products.Queries
	ImageService    ImageService
	DB              *pgxpool.Pool
}

func NewProductService(productsQueries *products.Queries, imageService ImageService, db *pgxpool.Pool) ProductService {
	return &ProductServiceImpl{
		ProductsQueries: productsQueries,
		ImageService:    imageService,
		DB:              db,
	}
}

func (s *ProductServiceImpl) CreateProduct(ctx context.Context, product *dto.CreateProductRequest, imageHeader *multipart.FileHeader) (string, error) {

	price, err := pkg.Float64ToNumeric(product.Price)
	if err != nil {
		configs.Log.Error("failed to convert price to numeric", zap.Error(err))
		return "", pkg.ErrInvalidPrice
	}
	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		configs.Log.Error("failed to begin transaction", zap.Error(err))
		return "", pkg.ErrFailedToCreateProduct
	}
	var (
		imageKey string
		finalErr error
	)
	defer func() {
		if finalErr != nil {
			tx.Rollback(ctx)
			if imageKey != "" {
				s.ImageService.DeleteImage(ctx, &imageKey)
			}
		}
	}()
	imageId, key, err := s.ImageService.UploadImage(ctx, tx, imageHeader)
	if err != nil {
		finalErr = err
		return "", finalErr
	}
	imageKey = key
	productId := ulid.Make().String()
	newProduct, err := s.ProductsQueries.WithTx(tx).CreateProduct(ctx, products.CreateProductParams{
		ID:      productId,
		StoreID: product.StoreID,
		Name:    product.Name,
		ImageID: pgtype.Text{
			String: imageId,
			Valid:  imageId != "",
		},
		Description: pgtype.Text{String: product.Description, Valid: true},
		Price:       price,
		IsAvailable: pgtype.Bool{Bool: product.IsAvailable, Valid: true},
	})
	if err != nil {
		if pgErrm, ok := err.(*pgconn.PgError); ok {
			if pgErrm.Code == "23503" {
				tx.Rollback(ctx)
				configs.Log.Error("Store not found", zap.String("store_id", product.StoreID))
				finalErr = pkg.ErrStoreNotFound
				return "", finalErr
			}

		}
		tx.Rollback(ctx)
		configs.Log.Error("failed to create product", zap.Error(err))
		finalErr = pkg.ErrFailedToCreateProduct
		return "", finalErr
	}
	if err := tx.Commit(ctx); err != nil {
		configs.Log.Error("failed to commit transaction", zap.Error(err))
		finalErr = pkg.ErrFailedToCreateProduct
		return "", finalErr
	}
	configs.Log.Info("product created successfully", zap.String("product_id", newProduct.ID))
	return newProduct.ID, nil
}

func (s *ProductServiceImpl) ListProducts(ctx context.Context, storeId string) ([]dto.GetAllProductResponse, error) {
	productsList, err := s.ProductsQueries.GetListProductsByStoreId(ctx, storeId)
	if err != nil {
		configs.Log.Error("failed to list products", zap.Error(err))
		return nil, errors.New("FAILED_TO_LIST_PRODUCTS")
	}
	response := make([]dto.GetAllProductResponse, 0)
	for _, p := range productsList {
		priceFloat, err := p.Price.Float64Value()
		if err != nil {
			configs.Log.Error("failed to convert price to float64", zap.Error(err))
			return nil, errors.New("INVALID_PRICE_FORMAT")
		}
		product := dto.GetAllProductResponse{
			Id:          p.ID,
			Name:        p.Name,
			Description: p.Description.String,
			Price:       priceFloat.Float64,
			IsAvailable: p.IsAvailable.Bool,
			ImagePath:   p.Path.String,
			CreatedAt:   p.CreatedAt.Time.String(),
			UpdatedAt:   p.UpdatedAt.Time.String(),
		}
		response = append(response, product)
	}
	configs.Log.Info("products listed successfully", zap.Int("count", len(response)))
	return response, nil
}

func (s *ProductServiceImpl) DeleteProduct(ctx context.Context, productID string) error {
	product, err := s.GetProductDetails(ctx, productID)
	if err != nil {
		return pkg.ErrFailedToCreateProduct
	}
	tx, err := s.DB.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		configs.Log.Error("failed to begin transaction", zap.Error(err))
		return pkg.ErrFailedToDeleteProduct
	}
	var finalErr error
	defer func() {
		if finalErr != nil {
			tx.Rollback(ctx)
		}
	}()

	err = s.ProductsQueries.WithTx(tx).DeleteProduct(ctx, productID)
	if err != nil {
		configs.Log.Error("failed to delete product", zap.Error(err))
		finalErr = pkg.ErrFailedToDeleteProduct
		return finalErr
	}
	err = s.ImageService.DeleteImageWithMetadata(ctx, tx, product.ImageID, aws.String(product.ImagePath))
	if err != nil {
		finalErr = err
		return finalErr
	}
	if err := tx.Commit(ctx); err != nil {
		configs.Log.Error("failed to commit transaction", zap.Error(err))
		finalErr = pkg.ErrFailedToDeleteProduct
		return finalErr
	}
	configs.Log.Info("product deleted successfully", zap.String("product_id", productID))
	return nil
}

func (s *ProductServiceImpl) GetProductDetails(ctx context.Context, productID string) (*dto.GetProductDetailsResponse, error) {
	product, err := s.ProductsQueries.GetProductById(ctx, productID)
	if err != nil {
		configs.Log.Error("failed to get product details", zap.Error(err))
		return nil, pkg.ErrProductNotFound
	}
	priceFloat, err := product.Price.Float64Value()
	if err != nil {
		configs.Log.Error("failed to convert price to float64", zap.Error(err))
		return nil, pkg.ErrFailedToGetProductDetails
	}
	return &dto.GetProductDetailsResponse{
		Id:          product.ID,
		ImageID:     product.ImageID.String,
		StoreID:     product.StoreID,
		Name:        product.Name,
		Description: product.Description.String,
		Price:       priceFloat.Float64,
		IsAvailable: product.IsAvailable.Bool,
		ImagePath:   product.Path.String,
		CreatedAt:   product.CreatedAt.Time.String(),
		UpdatedAt:   product.UpdatedAt.Time.String(),
	}, nil

}

func (s *ProductServiceImpl) UpdateProduct(ctx context.Context, productID string, product *dto.CreateProductRequest) error {
	return nil
}
