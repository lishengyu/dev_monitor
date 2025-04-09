package validate

import (
	validator "github.com/go-playground/validator/v10"
)

func IsStructValidate(s any) bool {
	v := validator.New()
	if err := v.Struct(s); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return false
		}
		for _, err := range err.(validator.ValidationErrors) {
			if err.Tag() == "required" {
				return false
			}
		}
		return true
	}
	return true
}
