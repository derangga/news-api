package middleware

import (
	"newsapi/internal/config/env"

	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func SetupGlobalMiddleware(e *echo.Echo, config env.ApplicationConfig) {
	e.Use(middleware.ContextTimeoutWithConfig(middleware.ContextTimeoutConfig{Skipper: middleware.DefaultSkipper, Timeout: config.Timeout}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(echoprometheus.NewMiddleware("news_api_app"))
}
