package entity

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
)

type ArticleStatus string

const (
	StatusDraft     ArticleStatus = "draft"
	StatusPublished ArticleStatus = "published"
	StatusDeleted   ArticleStatus = "deleted"
)

type NewsArticle struct {
	ID          int           `db:"id"`
	Title       string        `db:"title"`
	Content     string        `db:"content"`
	Summary     *string       `db:"summary"`
	AuthorID    int           `db:"author_id"`
	Slug        string        `db:"slug"`
	Status      ArticleStatus `db:"status"`
	PublishedAt sql.NullTime  `db:"published_at"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	DeletedAt   sql.NullTime  `db:"deleted_at"`
}

type ActiveNewsWithTopic struct {
	ID          int            `db:"id"`
	Title       string         `db:"title"`
	Content     string         `db:"content"`
	Slug        string         `db:"slug"`
	AuthorName  string         `db:"name"`
	PublishedAt sql.NullTime   `db:"published_at"`
	Topics      pq.StringArray `db:"topics"`
}

type NewsArticleWithTopic struct {
	ID          int           `db:"id"`
	Title       string        `db:"title"`
	Content     string        `db:"content"`
	Summary     *string       `db:"summary"`
	AuthorID    int           `db:"author_id"`
	Slug        string        `db:"slug"`
	Status      ArticleStatus `db:"status"`
	PublishedAt sql.NullTime  `db:"published_at"`
	CreatedAt   time.Time     `db:"created_at"`
	UpdatedAt   time.Time     `db:"updated_at"`
	DeletedAt   sql.NullTime  `db:"deleted_at"`
	Topics      pq.Int32Array `db:"topic_ids"`
}
