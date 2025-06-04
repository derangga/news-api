package server

import (
	"fmt"
	"newsapi/internal/config/env"
	"newsapi/internal/handler"
	"newsapi/internal/middleware"
	"newsapi/internal/routing"

	"github.com/labstack/echo/v4"
)

type HttpServer struct {
	echo    *echo.Echo
	config  *env.Config
	handler handler.HandlerRegistry
}

func NewHttpServer(
	config *env.Config,
	handler handler.HandlerRegistry,
) *HttpServer {
	e := echo.New()
	middleware.SetupGlobalMiddleware(e, config.ApplicationConfig)

	server := &HttpServer{
		echo:    e,
		config:  config,
		handler: handler,
	}

	return server
}

func (s *HttpServer) ListenAndServe() error {
	serviceUrl := fmt.Sprintf("%s:%s", s.config.ApplicationConfig.Host, s.config.ApplicationConfig.Port)
	return s.echo.Start(serviceUrl)
}

func (s *HttpServer) ConnectCoreWithEcho() {
	appRoutes := routing.NewAppRoutes(s.echo, s.handler)
	appRoutes.RegisterRoute()
}
