package pkg

import "github.com/go-playground/validator/v10"

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
