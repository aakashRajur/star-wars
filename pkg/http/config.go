package http

import (
	"time"

	"github.com/aakashRajur/star-wars/pkg/types"
)

type Timeout struct {
	Read  time.Duration
	Write time.Duration
	Idle  time.Duration
}

type ServerConfig struct {
	Port    string
	Timeout Timeout
	SslCert string
	SslKey  string
	Logger  types.Logger
}