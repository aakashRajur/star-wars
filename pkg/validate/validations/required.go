package validations

import (
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func Required() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return errors.Errorf(`%s is required`, key)
		}
		return nil
	}
}
