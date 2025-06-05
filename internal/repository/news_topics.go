package repository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
)

type newsTopicsRepository struct {
	db *sqlx.DB
}

func NewNewsTopicsRepository(db *sqlx.DB) NewsTopicsRepository {
	return newsTopicsRepository{db: db}
}

func (r newsTopicsRepository) Create(ctx context.Context, articleID int, topicIDs []int) error {
	if len(topicIDs) == 0 {
		return nil
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	query := "INSERT INTO news_topics (news_article_id, topic_id) VALUES "
	values := make([]interface{}, 0, len(topicIDs)*2)
	valueArgs := make([]string, 0, len(topicIDs))

	for i, topicID := range topicIDs {
		valueArgs = append(valueArgs, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		values = append(values, articleID, topicID)
	}

	query += strings.Join(valueArgs, ",")

	_, err = tx.ExecContext(ctx, query, values...)
	if err != nil {
		if strings.Contains(err.Error(), "foreign key constraint") {
			return errors.New("one or more topic IDs are invalid")
		}
		return err
	}

	return tx.Commit()
}

func (r newsTopicsRepository) ReplaceArticleTopics(ctx context.Context, articleID int, topicIDs []int32) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get current (non-deleted) topic_ids
	var currentTopicIDs []int32
	err = tx.SelectContext(ctx, &currentTopicIDs,
		`SELECT topic_id FROM news_topics WHERE news_article_id = $1 AND deleted_at IS NULL`,
		articleID,
	)
	if err != nil {
		return err
	}

	// Convert slices to maps for quick lookup
	newTopicMap := make(map[int32]struct{}, len(topicIDs))
	for _, id := range topicIDs {
		newTopicMap[id] = struct{}{}
	}

	currentMap := make(map[int32]struct{}, len(currentTopicIDs))
	for _, id := range currentTopicIDs {
		currentMap[id] = struct{}{}
	}

	// Soft-delete topics that are no longer present
	for _, id := range currentTopicIDs {
		if _, keep := newTopicMap[id]; !keep {
			_, err := tx.ExecContext(ctx,
				`UPDATE news_topics SET deleted_at = NOW() WHERE news_article_id = $1 AND topic_id = $2 AND deleted_at IS NULL`,
				articleID, id)
			if err != nil {
				return err
			}
		}
	}

	// Insert new topics or undelete them if previously deleted
	for _, id := range topicIDs {
		_, err := tx.ExecContext(ctx,
			`INSERT INTO news_topics (news_article_id, topic_id, created_at, deleted_at)
			 VALUES ($1, $2, NOW(), NULL)
			 ON CONFLICT (news_article_id, topic_id)
			 DO UPDATE SET deleted_at = NULL`,
			articleID, id)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r newsTopicsRepository) DeleteByArticleID(ctx context.Context, articleID int) error {
	query := `UPDATE news_topics SET deleted_at = NOW() WHERE news_article_id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, articleID)
	return err
}
