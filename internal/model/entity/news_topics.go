package entity

import "time"

type NewsTopic struct {
	ID            int       `db:"id"`
	NewsArticleID int       `db:"news_article_id"`
	TopicID       int       ` db:"topic_id"`
	CreatedAt     time.Time `db:"created_at"`
}
