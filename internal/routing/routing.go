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

	users := r.echo.Group("/api/v1/users")
	r.registerGroupRoute(users, http.MethodPost, "", h.UsersHandler.CreateUser)

	topics := r.echo.Group("/api/v1/topics")
	r.registerGroupRoute(topics, http.MethodGet, "", h.TopicsHandler.GetTopics)
	r.registerGroupRoute(topics, http.MethodPost, "", h.TopicsHandler.CreateTopic)
	r.registerGroupRoute(topics, http.MethodPatch, "/:id", h.TopicsHandler.UpdateTopic)

	newsArticle := r.echo.Group("/api/v1/news")
	r.registerGroupRoute(newsArticle, http.MethodGet, "", h.NewsArticlesHandler.GetNewsArticles)
	r.registerGroupRoute(newsArticle, http.MethodPost, "", h.NewsArticlesHandler.CreateNews)
	r.registerGroupRoute(newsArticle, http.MethodGet, "/:slug", h.NewsArticlesHandler.GetNewsBySlug)
	r.registerGroupRoute(newsArticle, http.MethodPatch, "/:slug", h.NewsArticlesHandler.UpdateNewsArticle)
	r.registerGroupRoute(newsArticle, http.MethodDelete, "/:slug", h.NewsArticlesHandler.DeleteNewsArticle)
}

func (r AppRoutes) registerGroupRoute(g *echo.Group, method string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	g.Add(method, path, h, m...)
}

func (r AppRoutes) registerRoute(method string, path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) {
	r.echo.Add(method, path, h, m...)
}
