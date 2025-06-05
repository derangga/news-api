package repository

import (
	"context"
	"newsapi/internal/model/entity"

	"github.com/jmoiron/sqlx"
)

type usersRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) UsersRepository {
	return usersRepository{db: db}
}

func (r usersRepository) Create(ctx context.Context, entity *entity.User) error {
	query := `INSERT INTO users (name, email) VALUES (:name, :email) RETURNING id`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, entity)

	if row.Err() != nil {
		return row.Err()
	}

	return nil
}
