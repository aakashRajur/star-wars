package consul

import (
	"net/url"
	"time"
)

type Config struct {
	Address      string
	Scheme       string
	WatchTimeout time.Duration
	Datacenter   string
	Username     string
	Password     string
}

func (config Config) Url() url.URL {
	username := config.Username
	password := config.Password

	compiled := url.URL{
		Scheme: config.Scheme,
		Host:   config.Address,
	}

	if username != `` {
		if password != `` {
			compiled.User = url.UserPassword(username, password)
		} else {
			compiled.User = url.User(username)
		}
	}

	return compiled
}
