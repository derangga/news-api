package usecase

import (
	"context"
	"newsapi/internal/model/request"
)

type NewsUsecase interface {
	CreateNewsArticle(ctx context.Context, body request.CreateNewsArticleRequest)
}

type TopicsUsecase interface {
	CreateTopic(ctx context.Context, body request.CreateTopicRequest)
}
