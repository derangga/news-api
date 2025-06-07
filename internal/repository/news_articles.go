package repository

import (
	"context"
	"fmt"
	"newsapi/internal/model/dto"
	"newsapi/internal/model/entity"
	"strings"
	"time"

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

func (r newsArticlesRepository) Create(ctx context.Context, entity *entity.NewsArticle) (int, error) {
	query := `INSERT INTO news_articles 
		(title, content, summary, author_id, slug, status, published_at)
		VALUES (:title, :content, :summary, :author_id, :slug, :status, :published_at) 
		RETURNING id`

	stmt, err := r.db.PrepareNamedContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	row := stmt.QueryRowxContext(ctx, entity)

	if row.Err() != nil {
		return 0, row.Err()
	}

	err = row.Scan(&entity.ID)
	if err != nil {
		return 0, err
	}

	return entity.ID, nil
}

func (r newsArticlesRepository) GetArticleBySlug(ctx context.Context, slug string) (entity.NewsArticleWithTopic, error) {
	query := `
			SELECT a.id, a.title, a.content, a.summary, a.author_id, a.slug, 
			a.status, a.published_at, a.created_at, a.updated_at,
			COALESCE(array_agg(nt.topic_id ORDER BY nt.topic_id), '{}') AS topic_ids 
			FROM news_articles a
			INNER JOIN news_topics nt ON nt.news_article_id = a.id AND nt.deleted_at IS NULL
			WHERE slug = $1 AND a.deleted_at is NULL
			GROUP BY a.id`
	var newsArticle entity.NewsArticleWithTopic

	err := r.db.GetContext(ctx, &newsArticle, query, slug)
	if err != nil {
		return entity.NewsArticleWithTopic{}, err
	}

	return newsArticle, nil
}

func (r newsArticlesRepository) GetActiveArticleBySlug(ctx context.Context, slug string) (entity.ActiveNewsWithTopic, error) {
	query := `SELECT
				a.id,
				a.title,
				a.content,
				a.slug,
				a.published_at,
				u.name,
				COALESCE(array_agg(t.name ORDER BY t.name) FILTER (WHERE t.name IS NOT NULL), '{}') AS topics
			FROM news_articles a
			INNER JOIN users u on u.id = a.author_id 
			INNER JOIN news_topics nt ON nt.news_article_id = a.id AND nt.deleted_at IS NULL
			INNER JOIN topics t ON t.id = nt.topic_id
			WHERE a.slug = $1 AND a.deleted_at IS NULL
			GROUP BY a.id, u.name`

	var newsArticle entity.ActiveNewsWithTopic

	err := r.db.GetContext(ctx, &newsArticle, query, slug)
	if err != nil {
		return entity.ActiveNewsWithTopic{}, err
	}

	return newsArticle, nil
}

func (r newsArticlesRepository) GetAll(ctx context.Context, filter dto.NewsFilter) ([]entity.NewsArticleWithTopicID, error) {
	query := `
			SELECT
				na.id,
				na.title,
				na.summary,
				na.author_id,
				na.slug,
				na.status,
				na.published_at,
				na.created_at,
				ARRAY_AGG(nt.topic_id) AS topic_ids
			FROM news_articles na
			LEFT JOIN news_topics nt ON na.id = nt.news_article_id
			WHERE
				na.deleted_at IS NULL
				AND nt.deleted_at IS NULL
		`

	var args []interface{}
	paramIdx := 1

	if filter.Status != "" {
		query += fmt.Sprintf(" AND na.status = $%d", paramIdx)
		args = append(args, filter.Status)
		paramIdx++
	}
	if filter.TopicID != "" {
		query += fmt.Sprintf(" AND nt.topic_id = $%d", paramIdx)
		args = append(args, filter.TopicID)
		paramIdx++
	}

	query += " GROUP BY na.id, na.title, na.summary, na.author_id, na.slug, na.status, na.published_at, na.created_at"

	var newsArticles []entity.NewsArticleWithTopicID
	err := r.db.SelectContext(ctx, &newsArticles, query, args...)
	if err != nil {
		return nil, err
	}

	return newsArticles, nil
}

func (r newsArticlesRepository) UpdateArticleFields(
	ctx context.Context,
	news *entity.NewsArticleWithTopic,
	updateFields []string,
) error {
	query := "UPDATE news_articles SET "
	setClauses := make([]string, 0, len(updateFields)+1)
	args := make([]interface{}, 0, len(updateFields)+1)

	for i, field := range updateFields {
		switch field {
		case "title":
			setClauses = append(setClauses, fmt.Sprintf("title = $%d", i+1))
			args = append(args, news.Title)
		case "content":
			setClauses = append(setClauses, fmt.Sprintf("content = $%d", i+1))
			args = append(args, news.Content)
		case "summary":
			setClauses = append(setClauses, fmt.Sprintf("summary = $%d", i+1))
			args = append(args, news.Summary)
		case "slug":
			setClauses = append(setClauses, fmt.Sprintf("slug = $%d", i+1))
			args = append(args, news.Slug)
		case "status":
			setClauses = append(setClauses, fmt.Sprintf("status = $%d", i+1))
			args = append(args, news.Status)
		}
	}

	setClauses = append(setClauses, fmt.Sprintf("published_at = $%d", len(updateFields)+1))
	args = append(args, time.Now())

	setClauses = append(setClauses, fmt.Sprintf("updated_at = $%d", len(updateFields)+2))
	args = append(args, time.Now())

	query += strings.Join(setClauses, ", ") + " WHERE id = $" + fmt.Sprintf("%d", len(updateFields)+3)
	args = append(args, news.ID)

	_, err := r.db.ExecContext(ctx, query, args...)

	return err
}

func (r newsArticlesRepository) DeleteBySlug(ctx context.Context, slug string) error {
	query := `UPDATE news_articles SET deleted_at = NOW() WHERE slug = $1 AND deleted_at IS NULL`

	_, err := r.db.ExecContext(ctx, query, slug)
	return err
}
