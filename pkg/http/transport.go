package http

import (
	"net/http"
	"time"
)

func NewTransport(timeout time.Duration) *http.Client {
	return &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 10,
			MaxIdleConns:        10,
			IdleConnTimeout:     1 * time.Second,
		},
	}
}
