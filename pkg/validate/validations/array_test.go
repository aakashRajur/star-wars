package validations

import (
	"testing"

	"github.com/pkg/errors"

	"github.com/aakashRajur/star-wars/pkg/types"
)

func TestValidateArray(t *testing.T) {
	type arg struct {
		Key       string
		Value     interface{}
		Validator types.Validator
		Exists    bool
	}
	type testCase struct {
		Args     arg
		Expected error
	}

	testCases := map[string]testCase{
		`OBJECT`: {
			Args: arg{
				Key: `key1`,
				Value: map[string]interface{}{
					`field1`: 10,
					`field2`: `HELLO`,
				},
				Validator: ValidateString(),
				Exists:    true,
			},
			Expected: errors.Errorf(ArrayError, `key1`),
		},
		`HETEROGENEOUS_ARRAY`: {
			Args: arg{
				Key:       `key2`,
				Value:     []interface{}{`HELLO`, 2},
				Validator: ValidateString(),
				Exists:    true,
			},
			Expected: errors.Errorf(StringError, `key2[1]`),
		},
		`HOMOGENEOUS_ARRAY`: {
			Args: arg{
				Key:       `key2`,
				Value:     []interface{}{5, 2},
				Validator: ValidateInteger(),
				Exists:    true,
			},
			Expected: nil,
		},
	}

	for name, test := range testCases {
		args := test.Args
		validator := ValidateArray(args.Validator)

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
			t.Errorf("ValidateArray() = %+v, want %+v", e1, e2)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
