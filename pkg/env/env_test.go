package env

import (
	"os"
	"strconv"
	"testing"
)

const (
	ENV_TYPE_STRING = `STRING`
	ENV_TYPE_INT    = `INT`
	ENV_TYPE_BOOL   = `BOOL`
)

type Env struct {
	EnvType  string
	EnvKey   string
	EnvValue interface{}
}

var fakeStringEnvs = []Env{
	{
		EnvKey:   `PG_URI`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `postgresql://wiki-db@localhost:5433/star_wars`,
	},
	{
		EnvKey:   `REDIS_URI`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `redis://star_wars@localhost:6379`,
	},
	{
		EnvKey:   `HTTP_SSL_CERT`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `assets/localhost.crt`,
	},
	{
		EnvKey:   `HTTP_SSL_KEY`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `assets/localhost.key`,
	},
	{
		EnvKey:   `LOG_LEVEL`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `trace`,
	},
	{
		EnvKey:   `LOG_FORMATTER`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `text`,
	},
	{
		EnvKey:   `LOG_PREFIX`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `api-debug`,
	},
	{
		EnvKey:   `PROJECT_NAME`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `github.com/aakashRajur/star-wars`,
	},
	{
		EnvKey:   `PROJECT_PATH`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `/star-wars`,
	},
	{
		EnvKey:   `MAIN`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `cmd/main.go`,
	},
	{
		EnvKey:   `EXECUTABLE`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `/main`,
	},
	{
		EnvKey:   `PG_USER`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `api_server`,
	},
	{
		EnvKey:   `PG_DB`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `star_wars`,
	},
	{
		EnvKey:   `PG_BACKUP`,
		EnvType:  ENV_TYPE_STRING,
		EnvValue: `/backup.sql`,
	},
}

var fakeIntEnvs = []Env{
	{
		EnvKey:   `PG_POOL_LIMIT`,
		EnvType:  ENV_TYPE_INT,
		EnvValue: 2,
	},
	{
		EnvKey:   `HTTP_PORT`,
		EnvType:  ENV_TYPE_INT,
		EnvValue: 8000,
	},
	{
		EnvKey:   `HTTP_READ_TIMEOUT`,
		EnvType:  ENV_TYPE_INT,
		EnvValue: 5,
	},
	{
		EnvKey:   `HTTP_WRITE_TIMEOUT`,
		EnvType:  ENV_TYPE_INT,
		EnvValue: 10,
	},
	{
		EnvKey:   `HTTP_IDLE_TIMEOUT`,
		EnvType:  ENV_TYPE_INT,
		EnvValue: 15,
	},
}

var fakeBoolEnvs = []Env{
	{
		EnvKey:   `BOOL_TRUE`,
		EnvType:  ENV_TYPE_BOOL,
		EnvValue: true,
	},
	{
		EnvKey:   `BOOL_FALSE`,
		EnvType:  ENV_TYPE_BOOL,
		EnvValue: false,
	},
}

func init() {
	compiled := append(fakeStringEnvs, fakeIntEnvs...)
	compiled = append(compiled, fakeBoolEnvs...)
	for _, each := range compiled {
		switch each.EnvType {
		case ENV_TYPE_INT:
			{
				err := os.Setenv(each.EnvKey, strconv.Itoa(each.EnvValue.(int)))
				if err != nil {
					panic(err)
				}
				break
			}
		case ENV_TYPE_STRING:
			err := os.Setenv(each.EnvKey, each.EnvValue.(string))
			if err != nil {
				panic(err)
			}
			break
		case ENV_TYPE_BOOL:
			{
				err := os.Setenv(each.EnvKey, strconv.FormatBool(each.EnvValue.(bool)))
				if err != nil {
					panic(err)
				}
				break
			}
		}
	}
}

func TestGetInt(t *testing.T) {
	t.Logf(`Testing setting of Integer Envs`)
	for _, each := range fakeIntEnvs {
		value := GetInt(each.EnvKey)
		expected := each.EnvValue.(int)
		if value != expected {
			t.Errorf(`%s was set incorrectly, got: %d, expected: %d`, each.EnvKey, value, expected)
			continue
		} else {
			t.Logf(`✔ %s`, each.EnvKey)
		}
	}
}

func TestGetString(t *testing.T) {
	t.Logf(`Testing setting of String Envs`)
	for _, each := range fakeStringEnvs {
		value := GetString(each.EnvKey)
		expected := each.EnvValue.(string)
		if value != expected {
			t.Errorf(`%s was set incorrectly, got: %s, expected: %s`, each.EnvKey, value, expected)
			continue
		} else {
			t.Logf(`✔ %s`, each.EnvKey)
		}
	}
}

func TestGetBool(t *testing.T) {
	t.Logf(`Testing setting of Bool Envs`)
	for _, each := range fakeBoolEnvs {
		value := GetBool(each.EnvKey)
		expected := each.EnvValue.(bool)
		if value != expected {
			t.Errorf(`%s was set incorrectly, got: %v, expected: %v`, each.EnvKey, value, expected)
			continue
		} else {
			t.Logf(`✔ %s`, each.EnvKey)
		}
	}
}
