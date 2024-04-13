package validator

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func NewValidator() (*validator.Validate, error) {
	validate := validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("dateTimeOrNil", dateTimeOrNil); err != nil {
		return nil, err
	}

	return validate, nil
}

func dateTimeOrNil(fl validator.FieldLevel) bool {
	req, ok := fl.Field().Interface().(time.Time)
	if !ok {
		return false
	}

	return xor(req == time.Time{}, req.After(time.Now()))
}

func xor(x, y bool) bool {
	return (x || y) && !(x && y)
}
