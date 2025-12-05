package pkg

import (
	"errors"

	"github.com/go-playground/validator/v10"
)

func ParseValidationErrors(err error) map[string]string {
	validationErrors := make(map[string]string)

	if errs, ok := err.(validator.ValidationErrors); ok {
		for _, fieldErr := range errs {
			fieldName := CamelToSnake(fieldErr.Field())
			var errorMsg string

			switch fieldErr.Tag() {
			case "required":
				errorMsg = fieldName + " is required."
			case "email":
				errorMsg = fieldName + " must be a valid email address."
			case "min":
				errorMsg = fieldName + " must be at least " + fieldErr.Param() + " characters long."
			case "eqfield":
				errorMsg = fieldName + " must be equal to " + fieldErr.Param() + "."
			default:
				errorMsg = "Invalid value for " + fieldName + "."
			}

			validationErrors[fieldName] = errorMsg
		}
	}

	return validationErrors
}

var (
	ErrInvalidImageMimeType      = errors.New("INVALID_IMAGE_MIME_TYPE")
	ErrImageSizeExceedsLimit     = errors.New("IMAGE_SIZE_EXCEEDS_LIMIT")
	ErrFailedToUploadImage       = errors.New("FAILED_TO_UPLOAD_IMAGE")
	ErrInvalidPrice              = errors.New("INVALID_PRICE")
	ErrFailedToCreateProduct     = errors.New("FAILED_TO_CREATE_PRODUCT")
	ErrStoreNotFound             = errors.New("STORE_NOT_FOUND")
	ErrFailedDeleteObject        = errors.New("FAILED_TO_DELETE_OBJECT")
	ErrFailedToGetProductDetails = errors.New("FAILED_TO_GET_PRODUCT_DETAILS")
	ErrFailedToDeleteProduct     = errors.New("FAILED_TO_DELETE_PRODUCT")
	ErrFailedToDeleteImage       = errors.New("FAILED_TO_DELETE_IMAGE")
	ErrProductNotFound           = errors.New("PRODUCT_NOT_FOUND")
)
