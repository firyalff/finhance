package shared

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator"
)

func Validator() *validator.Validate {
	validate := validator.New()
	return validate
}

type ValidationErrBody struct {
	ErrorCode string `json:"error_code"`
	Value     string `json:"value"`
	Treshold  string `json:"threshold,omitempty"`
}

func ParseValidatorError(err error) map[string]ValidationErrBody {
	parsedErrors := make(map[string]ValidationErrBody)
	validationErrors := err.(validator.ValidationErrors)

	for _, validationError := range validationErrors {
		fieldName := strings.ToLower(validationError.StructField())
		parsedErrors[fieldName] = ValidationErrBody{
			ErrorCode: validationErrorMessage(validationError.Tag()),
			Value:     fmt.Sprintf("%v", validationError.Value()),
		}
	}

	return parsedErrors
}

func validationErrorMessage(validationRule string) string {
	const VALIDATION_ERR_PREFIX = "BAD_REQ"
	return VALIDATION_ERR_PREFIX + "_" + strings.ToUpper(validationRule)
}
