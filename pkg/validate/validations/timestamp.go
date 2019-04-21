package validations

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func ValidateTimestamp() types.Validator {
	return func(key string, value interface{}, exists bool) error {
		if !exists || value == nil {
			return nil
		}

		str := value.(string)

		_, err := time.Parse(`2006-01-02T15:04:05Z0700`, str)
		if err != nil {
			return err
		}
		return nil
	}
}
