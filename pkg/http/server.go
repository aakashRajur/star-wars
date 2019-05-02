package http

import (
	"context"
	"net/http"
	"sync/atomic"

	"github.com/aakashRajur/star-wars/pkg/types"
)

const ProtocolName = `HTTP`

type Server struct {
	http.Server
	logger  types.Logger
	healthy int32
	config  ServerConfig
}

func (server *Server) Name() string {
	return ProtocolName
}

func (server *Server) Start(context.Context) error {
	atomic.StoreInt32(&server.healthy, 1)

	sslCert := server.config.SslCert
	sslKey := server.config.SslKey

	go func() {
		server.logger.Info("SERVER WILL HANDLE REQUESTS AT", server.Addr)
		if err := server.ListenAndServeTLS(sslCert, sslKey); err != nil && err != http.ErrServerClosed {
			server.logger.Fatal(err)
		}
	}()
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
