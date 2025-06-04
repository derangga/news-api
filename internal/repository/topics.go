package repository

import (
	"context"
	"newsapi/internal/model/entity"

	"github.com/jmoiron/sqlx"
)

type topicRepository struct {
	db *sqlx.DB
}

func NewTopicRepository(db *sqlx.DB) TopicsRepository {
	return topicRepository{
		db: db,
	}
}

func (r topicRepository) Create(ctx context.Context, entity entity.Topic) error {
	return nil
}
