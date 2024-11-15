package helper

import (
	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func ValidateRequest(data interface{}) (map[string]string, error) {
	err := validate.Struct(data)
	if err == nil {
		return nil, nil
	}

	errors := make(map[string]string)
	for _, err := range err.(validator.ValidationErrors) {
		errors[err.Field()] = err.Tag()
	}
	return errors, err
}
