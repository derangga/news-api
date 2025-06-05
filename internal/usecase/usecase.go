package usecase

import (
	"context"
	"newsapi/internal/model/request"
	"newsapi/internal/model/response"
)

type UsersUsecase interface {
	CreateUser(ctx context.Context, body request.CreateUserRequest) error
}

type NewsUsecase interface {
	CreateNewsArticle(ctx context.Context, body request.CreateNewsArticleRequest)
}

type TopicsUsecase interface {
	CreateTopic(ctx context.Context, body request.CreateTopicRequest) error
	GetTopics(ctx context.Context) ([]response.Topic, error)
	UpdateTopic(ctx context.Context, body request.UpdateTopicRequest) error
}
