package entity

import (
	"database/sql"
	"time"
)

type NewsTopic struct {
	ID            int          `db:"id"`
	NewsArticleID int          `db:"news_article_id"`
	TopicID       int          ` db:"topic_id"`
	CreatedAt     time.Time    `db:"created_at"`
	DeletedAt     sql.NullTime `db:"deleted_at"`
}
