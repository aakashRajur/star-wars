package normalizations

import (
	"fmt"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func NormalizeString() types.Normalizor {
	return func(key string, value interface{}) (interface{}, error) {
		if value == nil {
			return nil, nil
		}
		sValue, ok := value.(string)
		if !ok {
			return nil, errors.Errorf(`%s is not a valid string`, key)
		}
		return fmt.Sprintf(`"%s"`, sValue), nil
	}
}
