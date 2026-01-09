package server

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type httpServer struct {
	server *http.Server
}

var server *httpServer

func Init(router *gin.Engine) *httpServer {
	if server != nil {
		return server
	}

	httpserver := &http.Server{
		Handler: router.Handler(),
		BaseContext: func(net.Listener) context.Context {
			return context.Background()
		},
	}

	server := &httpServer{
		server: httpserver,
	}

	return server
}

func (server *httpServer) Run(port int) error {
	network := "tcp4"
	address := ":" + strconv.Itoa(port)

	listener, err := net.Listen(network, address)
	if err != nil {
		err = fmt.Errorf("failed to create listener for network=%s and address=%s: %w", network, address, err)
		return err
	}

	return server.server.Serve(listener)
}
