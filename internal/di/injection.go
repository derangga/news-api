package di

import (
	"newsapi/internal/config/db"
	"newsapi/internal/config/env"
	"newsapi/internal/config/server"
	"newsapi/internal/handler"
	"newsapi/internal/repository"
	"newsapi/internal/usecase"

	"github.com/go-playground/validator/v10"
	"github.com/jmoiron/sqlx"
)

func provideValidator() *validator.Validate {
	return validator.New()
}

func provideConfig() *env.Config {
	return env.BuildConfig()
}

func provideDB(config env.DatabaseConfig) *sqlx.DB {
	return db.NewPostgresDatabase(config)
}

func provideTopicsRepository(db *sqlx.DB) repository.TopicsRepository {
	return repository.NewTopicRepository(db)
}

func provideNewsArticlesRepository(db *sqlx.DB) repository.NewsArticlesRepository {
	return repository.NewNewsArticlesRepository(db)
}

func provideNewsTopicsRepository(db *sqlx.DB) repository.NewsTopicsRepository {
	return repository.NewNewsTopicsRepository(db)
}

func provideTopicsUsecase(repos repository.TopicsRepository) usecase.TopicsUsecase {
	return usecase.NewTopicsUsecase(repos)
}

func provideNewsUsecase(
	newsArticles repository.NewsArticlesRepository,
	newsTopics repository.NewsTopicsRepository,
) usecase.NewsUsecase {
	return usecase.NewNewsArticlesUsecase(newsArticles, newsTopics)
}

func provideTopicsHandler(
	validator *validator.Validate,
	uc usecase.TopicsUsecase,
) handler.TopicsHandler {
	return handler.NewTopicsHandler(validator, uc)
}

func provideNewsHandler(validator *validator.Validate,
	uc usecase.NewsUsecase,
) handler.NewsHandler {
	return handler.NewNewsArticlesHandler(validator, uc)
}

func provideHandlerRegistry(
	topicsHandler handler.TopicsHandler,
	newsHandler handler.NewsHandler,
) handler.HandlerRegistry {
	return handler.NewHandlerRegistry(topicsHandler, newsHandler)
}

func provideHttpServer(env *env.Config, handler handler.HandlerRegistry) *server.HttpServer {
	return server.NewHttpServer(env, handler)
}

func InitHTTPServer() *server.HttpServer {
	validator := provideValidator()
	config := provideConfig()
	sqlClient := provideDB(config.DatabaseConfig)
	topicsRepo := provideTopicsRepository(sqlClient)
	newsArticlesRepo := provideNewsArticlesRepository(sqlClient)
	newsTopicsRepo := provideNewsTopicsRepository(sqlClient)
	topicsUC := provideTopicsUsecase(topicsRepo)
	newsUC := provideNewsUsecase(newsArticlesRepo, newsTopicsRepo)
	topicsHandler := provideTopicsHandler(validator, topicsUC)
	newsHandler := provideNewsHandler(validator, newsUC)
	handlerRegistry := provideHandlerRegistry(topicsHandler, newsHandler)
	httpServer := provideHttpServer(config, handlerRegistry)

	return httpServer
}
