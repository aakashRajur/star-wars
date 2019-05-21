package validations

import (
	"testing"

	"github.com/juju/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func TestValidateShape(t *testing.T) {
	type arg struct {
		Key    string
		Shape  map[string][]types.Validator
		Value  interface{}
		Exists bool
	}
	type testCase struct {
		Args     arg
		Expected error
	}

	testCases := map[string]testCase{
		`INCORRECT_DATA_TYPE`: {
			Args: arg{
				Key: `key1`,
				Shape: map[string][]types.Validator{
					`field1`: {ValidateRequired(), ValidateString()},
					`field2`: {ValidateInteger()},
				},
				Value:  `HELLO_WORLD`,
				Exists: true,
			},
			Expected: errors.Errorf(ShapeError, `key1`),
		},
		`INCORRECT_FIELD_VALUE_TYPE`: {
			Args: arg{
				Key: `key1`,
				Shape: map[string][]types.Validator{
					`field1`: {ValidateRequired(), ValidateString()},
					`field2`: {ValidateInteger()},
				},
				Value: map[string]interface{}{
					`field1`: 10,
					`field2`: 56,
				},
				Exists: true,
			},
			Expected: errors.Errorf(StringError, `field1`),
		},
		`CORRECT_FIELD_VALUE_TYPE`: {
			Args: arg{
				Key: `key1`,
				Shape: map[string][]types.Validator{
					`field1`: {ValidateRequired(), ValidateString()},
					`field2`: {ValidateInteger()},
					`field3`: {ValidateFloat()},
				},
				Value: map[string]interface{}{
					`field1`: `HELLO_WORLD`,
					`field2`: 56,
				},
				Exists: true,
			},
			Expected: nil,
		},
	}

	for name, test := range testCases {
		args := test.Args
		validator := ValidateShape(args.Shape)

		got := validator(args.Key, args.Value, args.Exists)
		expected := test.Expected

		e1 := `nil`
		e2 := `nil`

		if got != nil {
			e1 = got.Error()
		}
		if expected != nil {
			e2 = expected.Error()
		}

		if e1 != e2 {
			t.Errorf("ValidateShape() = %+v, want %+v", e1, e2)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
