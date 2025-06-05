package repository

import (
	"context"
	"fmt"
	"newsapi/internal/model/entity"
	"strings"
	"time"

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

func (r topicRepository) Create(ctx context.Context, entity *entity.Topic) error {
	query := `INSERT INTO topics (name, description, slug) VALUES (:name, :description, :slug)`
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

func (r topicRepository) GetAll(ctx context.Context) ([]entity.Topic, error) {
	query := `
	SELECT id, name, description, slug, created_at, updated_at, deleted_at 
	FROM topics 
	WHERE deleted_at IS NULL
`
	var topics []entity.Topic
	err := r.db.SelectContext(ctx, &topics, query)
	if err != nil {
		return nil, err
	}

	return topics, nil
}

func (r topicRepository) GetByID(ctx context.Context, id int) (entity.Topic, error) {
	query := `
		SELECT id, name, description, slug, created_at, updated_at 
		FROM topics 
		WHERE id = $1 AND deleted_at IS NULL
	`
	var topic entity.Topic

	err := r.db.GetContext(ctx, &topic, query, id)
	if err != nil {
		return entity.Topic{}, err
	}

	return topic, nil
}

func (r topicRepository) UpdateTopicFileds(ctx context.Context, topic *entity.Topic, updateFields []string) error {
	query := "UPDATE topics SET "
	setClauses := make([]string, 0, len(updateFields)+1)
	args := make([]interface{}, 0, len(updateFields)+1)

	for i, field := range updateFields {
		switch field {
		case "name":
			setClauses = append(setClauses, fmt.Sprintf("name = $%d", i+1))
			args = append(args, topic.Name)
		case "description":
			setClauses = append(setClauses, fmt.Sprintf("description = $%d", i+1))
			args = append(args, topic.Description)
		case "slug":
			setClauses = append(setClauses, fmt.Sprintf("slug = $%d", i+1))
			args = append(args, topic.Slug)
		}
	}

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", len(updateFields)+1))
	args = append(args, time.Now())

	query += strings.Join(setClauses, ", ") + " WHERE id = $" + fmt.Sprintf("%d", len(updateFields)+2)
	args = append(args, topic.ID)

	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}
