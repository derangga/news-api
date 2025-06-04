package repository

import (
	"context"
	"newsapi/internal/model/entity"

	"github.com/jmoiron/sqlx"
)

type newsArticlesRepository struct {
	db *sqlx.DB
}

func NewNewsArticlesRepository(db *sqlx.DB) NewsArticlesRepository {
	return newsArticlesRepository{
		db: db,
	}
}

func (r newsArticlesRepository) Create(ctx context.Context, entity entity.NewsArticle) error {
	return nil
}
