package validation

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

var (
	mobileRegex   = regexp.MustCompile(`^1[3-9]\d{9}$`)
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_-]{4,16}$`)
)

// RegisterCustomValidators registers custom validation rules
func RegisterCustomValidators(v *validator.Validate) error {
	if err := v.RegisterValidation("mobile", validateMobile); err != nil {
		return err
	}
	if err := v.RegisterValidation("username", validateUsername); err != nil {
		return err
	}
	return nil
}

func validateMobile(fl validator.FieldLevel) bool {
	mobile := strings.TrimSpace(fl.Field().String())
	return mobileRegex.MatchString(mobile)
}

func validateUsername(fl validator.FieldLevel) bool {
	username := strings.TrimSpace(fl.Field().String())
	return usernameRegex.MatchString(username)
}
