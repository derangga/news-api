package repository

import (
	"context"
	"newsapi/internal/model/entity"
)

type UsersRepository interface {
	Create(ctx context.Context, entity *entity.User) error
}

type TopicsRepository interface {
	Create(ctx context.Context, entity *entity.Topic) error
	GetAll(ctx context.Context) ([]entity.Topic, error)
	GetByID(ctx context.Context, id int) (entity.Topic, error)
	Update(ctx context.Context, topic *entity.Topic, updateFields []string) error
}

type NewsArticlesRepository interface {
	Create(ctx context.Context, entity entity.NewsArticle) error
}

type NewsTopicsRepository interface {
	Create(ctx context.Context, entity entity.NewsTopic) error
}
