package validations

import (
	"testing"

	"github.com/juju/errors"
)

func TestValidateInteger(t *testing.T) {
	type arg struct {
		Key    string
		Value  interface{}
		Exists bool
	}
	type testCase struct {
		Args     arg
		Expected error
	}

	testCases := map[string]testCase{
		`INDEX_1`: {
			Args: arg{
				Key:    `key`,
				Value:  4.0,
				Exists: true,
			},
			Expected: nil,
		},
		`INDEX_2`: {
			Args: arg{
				Key:    `key`,
				Value:  5,
				Exists: true,
			},
			Expected: nil,
		},
		`FLOAT`: {
			Args: arg{
				Key:    `key1`,
				Value:  10.5,
				Exists: true,
			},
			Expected: errors.Errorf(IntegerError, `10.5`, `key1`),
		},
		`INT`: {
			Args: arg{
				Key:    `key1`,
				Value:  16,
				Exists: true,
			},
			Expected: nil,
		},
		`STRING`: {
			Args: arg{
				Key:    `key2`,
				Value:  `445.98`,
				Exists: true,
			},
			Expected: errors.Errorf(IntegerError, `445.98`, `key2`),
		},
		`BOOL`: {
			Args: arg{
				Key:    `key3`,
				Value:  true,
				Exists: true,
			},
			Expected: errors.Errorf(IntegerError, true, `key3`),
		},
		`EXISTS`: {
			Args: arg{
				Key:    `key4`,
				Value:  nil,
				Exists: false,
			},
			Expected: nil,
		},
	}

	validator := ValidateInteger()

	for name, test := range testCases {
		args := test.Args
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
			t.Errorf("ValidateInteger() = %+v, want %+v", e1, e2)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
