package normalizations

import (
	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NormalizeInteger() types.Normalizor {
	return func(key string, value interface{}) (interface{}, error) {
		if value == nil {
			return nil, nil
		}
		fValue, ok := value.(float64)
		if !ok {
			return nil, errors.Errorf(`%s is not a valid int`, key)
		}
		return int(fValue), nil
	}
}
