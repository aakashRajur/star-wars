package validations

import (
	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
	"github.com/aakashRajur/star-wars/pkg/util"
)

const (
	StringError = `%s SHOULD BE OF TYPE STRING`
)

func ValidateString() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}
		if util.GetType(value) != `string` {
			return errors.Errorf(StringError, key)
		}
		return nil
	}
}
