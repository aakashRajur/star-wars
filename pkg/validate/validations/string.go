package validations

import (
	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

func ValidateString() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		if util.GetType(value) != `string` {
			return errors.Errorf(`%s should be of type string`, key)
		}
		return nil
	}
}
