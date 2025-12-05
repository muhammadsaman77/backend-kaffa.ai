package services

import (
	"context"
	"mime/multipart"
	"strings"

	"backend-kaffa.ai/configs"
	"backend-kaffa.ai/internal/sqlc/images"
	"backend-kaffa.ai/pkg"
	"go.uber.org/zap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/jackc/pgx/v5"
	"github.com/oklog/ulid/v2"
)

type ImageServiceImpl struct {
	S3Client     *s3.Client
	Bucket       *string
	ImageQueries *images.Queries
}

func NewImageService(s3Client *s3.Client, bucket *string, imageQueries *images.Queries) ImageService {
	return &ImageServiceImpl{
		S3Client:     s3Client,
		Bucket:       bucket,
		ImageQueries: imageQueries,
	}
}

func (s *ImageServiceImpl) UploadImage(ctx context.Context, tx pgx.Tx, imageHeader *multipart.FileHeader) (string, string, error) {
	mimeType := imageHeader.Header.Get("Content-Type")
	if !strings.HasPrefix(mimeType, "image/") {
		configs.Log.Error("invalid mime type for image", zap.String("mimeType", mimeType))
		return "", "", pkg.ErrInvalidImageMimeType
	}
	size := imageHeader.Size
	if size > 5*1024*1024 { // 5 MB limit
		configs.Log.Error("image size exceeds limit", zap.Int64("size", size))
		return "", "", pkg.ErrImageSizeExceedsLimit
	}
	fileName := imageHeader.Filename
	imageId := ulid.Make().String()
	keyName := "products/" + imageId + "_" + fileName
	file, err := imageHeader.Open()

	if err != nil {
		configs.Log.Error("failed to open image file", zap.Error(err))
		return "", "", pkg.ErrFailedToUploadImage
	}
	defer file.Close()

	_, err = s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      s.Bucket,
		Key:         aws.String(keyName),
		Body:        file,
		ContentType: aws.String(mimeType),
	})
	if err != nil {
		configs.Log.Error("failed to upload image to S3", zap.Error(err))
		return "", "", pkg.ErrFailedToUploadImage
	}
	_, err = s.ImageQueries.WithTx(tx).CreateImage(ctx, images.CreateImageParams{
		ID:           imageId,
		OriginalName: fileName,
		Size:         int32(size),
		MimeType:     mimeType,
		Path:         keyName,
	})
	if err != nil {
		configs.Log.Error("failed to save image metadata", zap.Error(err))
		err = s.DeleteImage(ctx, aws.String(keyName))
		if err != nil {
			return "", "", pkg.ErrFailedToUploadImage
		}
		return "", "", pkg.ErrFailedToUploadImage
	}
	configs.Log.Info("image uploaded successfully", zap.String("image_id", imageId), zap.String("key_name", keyName))
	return imageId, keyName, nil
}

func (s *ImageServiceImpl) DeleteImage(ctx context.Context, keyName *string) error {
	_, err := s.S3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: s.Bucket,
		Key:    keyName,
	})
	if err != nil {
		configs.Log.Error("failed to delete image from S3", zap.String("key_name", *keyName), zap.Error(err))
		return pkg.ErrFailedDeleteObject
	}
	configs.Log.Info("image deleted successfully from S3", zap.String("key_name", *keyName))
	return nil

}

func (s *ImageServiceImpl) GetImageURL(ctx context.Context, imageID string) (string, error) {
	// Implement logic to get image URL from S3 here
	return "", nil
}

func (s *ImageServiceImpl) DeleteImageWithMetadata(ctx context.Context, tx pgx.Tx, imageID string, keyName *string) error {
	err := s.DeleteImage(ctx, keyName)
	if err != nil {
		return err
	}
	err = s.ImageQueries.WithTx(tx).DeleteImage(ctx, imageID)
	if err != nil {
		configs.Log.Error("failed to delete image metadata", zap.String("image_id", imageID), zap.Error(err))
		return pkg.ErrFailedToDeleteImage
	}
	configs.Log.Info("image metadata deleted successfully", zap.String("image_id", imageID))
	return nil
}
