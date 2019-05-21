package validations

import (
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const (
	RequiredError = `%s is required`
)

func ValidateRequired() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return errors.Errorf(RequiredError, key)
		}
		return nil
	}
}
