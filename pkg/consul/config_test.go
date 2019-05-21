package consul

import (
	"net/url"
	"reflect"
	"testing"
	"time"
)

const (
	address      = `consul:8500`
	scheme       = `http`
	watchTimeout = 5 * time.Minute
	username     = `hello`
	password     = `example1245`
)

type testCase struct {
	Args     Config
	Expected url.URL
}

var testCases = map[string]testCase{
	`DEFAULT`: {
		Args: Config{
			Address:      address,
			Scheme:       scheme,
			WatchTimeout: watchTimeout,
		},
		Expected: url.URL{
			Scheme: scheme,
			Host:   address,
		},
	},
	`ONLY_USER`: {
		Args: Config{
			Address:      address,
			Scheme:       scheme,
			WatchTimeout: watchTimeout,
			Username:     username,
		},
		Expected: url.URL{
			Scheme: scheme,
			Host:   address,
			User:   url.User(username),
		},
	},
	`USER_WITH_PASSWORD`: {
		Args: Config{
			Address:      address,
			Scheme:       scheme,
			WatchTimeout: watchTimeout,
			Username:     username,
			Password:     password,
		},
		Expected: url.URL{
			Scheme: scheme,
			Host:   address,
			User:   url.UserPassword(username, password),
		},
	},
}

func TestConfig_Url(t *testing.T) {
	for name, test := range testCases {
		got := test.Args.Url()
		if !reflect.DeepEqual(got, test.Expected) {
			t.Errorf(`Config.Url() = %+v, want %+v`, got, test.Expected)
		} else {
			t.Logf(`✔ %s`, name)
		}
	}
}
