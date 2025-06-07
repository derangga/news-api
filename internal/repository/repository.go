package repository

import (
	"context"
	"newsapi/internal/model/dto"
	"newsapi/internal/model/entity"
)

type UsersRepository interface {
	Create(ctx context.Context, entity *entity.User) error
}

type TopicsRepository interface {
	Create(ctx context.Context, entity *entity.Topic) error
	GetAll(ctx context.Context) ([]entity.Topic, error)
	GetByID(ctx context.Context, id int) (entity.Topic, error)
	UpdateTopicFileds(ctx context.Context, topic *entity.Topic, updateFields []string) error
	Delete(ctx context.Context, id int) error
}

type NewsArticlesRepository interface {
	Create(ctx context.Context, entity *entity.NewsArticle) (int, error)
	GetArticleBySlug(ctx context.Context, slug string) (entity.NewsArticleWithTopic, error)
	GetActiveArticleBySlug(ctx context.Context, slug string) (entity.ActiveNewsWithTopic, error)
	GetAll(ctx context.Context, filter dto.NewsFilter) ([]entity.NewsArticleWithTopicID, error)
	UpdateArticleFields(ctx context.Context, entity *entity.NewsArticleWithTopic, updateFields []string) error
	DeleteBySlug(ctx context.Context, slug string) error
}

type NewsTopicsRepository interface {
	Create(ctx context.Context, articleID int, topicIDs []int) error
	ReplaceArticleTopics(ctx context.Context, articleID int, topicIDs []int32) error
	DeleteByArticleID(ctx context.Context, articleID int) error
}
