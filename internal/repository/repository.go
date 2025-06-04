package repository

import (
	"context"
	"newsapi/internal/model/entity"
)

type TopicsRepository interface {
	Create(ctx context.Context, entity entity.Topic) error
}

type NewsArticlesRepository interface {
	Create(ctx context.Context, entity entity.NewsArticle) error
}

type NewsTopicsRepository interface {
	Create(ctx context.Context, entity entity.NewsTopic) error
}
