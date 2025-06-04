package repository

import (
	"context"
	"newsapi/internal/model/entity"

	"github.com/jmoiron/sqlx"
)

type newsTopicsRepository struct {
	db *sqlx.DB
}

func NewNewsTopicsRepository(db *sqlx.DB) NewsTopicsRepository {
	return newsTopicsRepository{db: db}
}

func (r newsTopicsRepository) Create(ctx context.Context, entity entity.NewsTopic) error {
	return nil
}
