package routing

import (
	"net/http"
	"newsapi/internal/handler"

	"github.com/labstack/echo/v4"
)

type AppRoutes struct {
	echo    *echo.Echo
	handler handler.HandlerRegistry
}

func NewAppRoutes(
	echo *echo.Echo,
	handler handler.HandlerRegistry,
) AppRoutes {
	return AppRoutes{
		echo:    echo,
		handler: handler,
	}
}

func (r AppRoutes) RegisterRoute() {
	h := r.handler
	r.registerRoute(r.echo, http.MethodGet, "/", h.TopicsHandler.GetTopics)
}

func (r AppRoutes) registerRoute(e *echo.Echo, method string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	e.Add(method, path, h, m...)
}
