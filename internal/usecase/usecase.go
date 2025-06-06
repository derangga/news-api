package usecase

import (
	"context"
	"newsapi/internal/model/dto"
	"newsapi/internal/model/request"
	"newsapi/internal/model/response"
)

type UsersUsecase interface {
	CreateUser(ctx context.Context, body request.CreateUserRequest) error
}

type NewsUsecase interface {
	CreateNewsArticle(ctx context.Context, body request.CreateNewsArticleRequest) error
	GetNewsArticles(ctx context.Context, filter dto.NewsFilter) ([]response.NewsArticle, error)
	GetNewsArticleBySlug(ctx context.Context, slug string) (response.NewsArticleWithTopic, error)
	UpdateNewsArticleBySlug(ctx context.Context, slug string, body request.UpdateNewsArticleRequest) error
	DeleteNewsArticleBySlug(ctx context.Context, slug string) error
}

type TopicsUsecase interface {
	CreateTopic(ctx context.Context, body request.CreateTopicRequest) error
	GetTopics(ctx context.Context) ([]response.Topic, error)
	UpdateTopic(ctx context.Context, id int, body request.UpdateTopicRequest) error
}
