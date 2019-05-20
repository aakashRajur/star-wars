package args

import (
	"flag"
	"os"
	"testing"
)

type testCase struct {
	Args  []string
	Value string
}

var conditions = map[string]testCase{
	`DEFAULT`: {
		Value: `dev`,
		Args:  []string{"cmd"},
	},
	`DEV`: {
		Value: `dev`,
		Args:  []string{"cmd", "-env=dev"},
	},
	`PROD`: {
		Value: `prod`,
		Args:  []string{"cmd", "-env=prod"},
	},
}

func TestLoadArgs(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	for name, tc := range conditions {
		flag.CommandLine = flag.NewFlagSet(tc.Args[0], flag.ContinueOnError)
		flag.CommandLine.Usage = flag.Usage
		os.Args = tc.Args

		LoadArgs()
		value := os.Getenv(`env`)
		if value != tc.Value {
			t.Errorf(`%s failed, got: %s, expected: %s`, name, value, tc.Value)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
