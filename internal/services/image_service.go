package services

import (
	"context"
	"mime/multipart"

	"github.com/jackc/pgx/v5"
)

type ImageService interface {
	UploadImage(ctx context.Context, tx pgx.Tx, imageHeader *multipart.FileHeader) (string, string, error)
	GetImageURL(ctx context.Context, imageID string) (string, error)
	DeleteImage(ctx context.Context, keyName *string) error
	DeleteImageWithMetadata(ctx context.Context, tx pgx.Tx, imageID string, keyName *string) error
}
