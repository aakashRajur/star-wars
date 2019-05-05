package http

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const ProtocolName = `HTTP`
const SecureProtocolName = `HTTPS`

type Server struct {
	http.Server
	logger  types.Logger
	healthy int32
	config  ServerConfig
}

func (server *Server) Name() string {
	sslCert := server.config.SslCert
	sslKey := server.config.SslKey
	if sslCert != `` && sslKey != `` {
		return SecureProtocolName
	}
	return ProtocolName
}

func (server *Server) Start(context.Context) error {
	atomic.StoreInt32(&server.healthy, 1)

	sslCert := server.config.SslCert
	sslKey := server.config.SslKey

	go func(sslCert, sslKey string) {
		server.logger.Info("SERVER WILL HANDLE REQUESTS AT", server.Addr)
		if sslCert != `` && sslKey != `` {
			if err := server.ListenAndServeTLS(sslCert, sslKey); err != nil && err != http.ErrServerClosed {
				server.logger.Fatal(err)
			}
		} else {
			if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				server.logger.Fatal(err)
			}
		}
	}(sslCert, sslKey)
	return nil
}

func (server *Server) Stop(ctx context.Context) error {
	server.logger.Info("ATTEMPTING TO SHUTDOWN SERVER")
	atomic.StoreInt32(&server.healthy, 0)
	err := server.Shutdown(ctx)
	if err == nil {
		server.logger.Info("SERVER SHUTDOWN SUCCESSFUL")
	}
	return err
}

func (server *Server) GetStatus() int {
	if atomic.LoadInt32(&server.healthy) == 1 {
		return http.StatusOK
	}
	return http.StatusServiceUnavailable
}
