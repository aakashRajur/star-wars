package http

import (
	"net/http"
)

func NewInstance(config ServerConfig, router *Router) *Server {
	server := Server{
		config: config,
		Server: http.Server{
			Addr:         config.Port,
			ReadTimeout:  config.Timeout.Read,
			WriteTimeout: config.Timeout.Write,
			IdleTimeout:  config.Timeout.Idle,
			Handler:      router,
		},
		logger:  config.Logger,
		healthy: 0,
	}

	return &server
}
