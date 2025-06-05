package response

import (
	"newsapi/internal/model/entity"
	"time"
)

type NewsArticle struct {
	ID          int                  `json:"id"`
	Title       string               `json:"title"`
	Content     string               `json:"content"`
	Summary     *string              `json:"summary"`
	AuthorID    int                  `json:"author_id"`
	Slug        string               `json:"slug" db:"slug"`
	Status      entity.ArticleStatus `json:"status"`
	PublishedAt *time.Time           `json:"published_at"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   *time.Time           `json:"deleted_at"`
}

func NewsArticleSeriliazer(entity entity.NewsArticle) NewsArticle {
	pub := &entity.PublishedAt.Time
	del := &entity.DeletedAt.Time

	if !entity.PublishedAt.Valid {
		pub = nil
	}
	if !entity.DeletedAt.Valid {
		del = nil
	}
	return NewsArticle{
		ID:          entity.ID,
		Title:       entity.Title,
		Content:     entity.Content,
		Summary:     entity.Summary,
		AuthorID:    entity.AuthorID,
		Slug:        entity.Slug,
		Status:      entity.Status,
		PublishedAt: pub,
		CreatedAt:   entity.CreatedAt,
		UpdatedAt:   entity.UpdatedAt,
		DeletedAt:   del,
	}
}

type NewsArticleWithTopic struct {
	ID          int        `json:"id"`
	Title       string     `json:"title"`
	Content     string     `json:"content"`
	Slug        string     `json:"slug"`
	AuthorName  string     `json:"author_name"`
	PublishedAt *time.Time `json:"published_at"`
	Topics      []string   `json:"topics"`
}

func NewsArticleWithTopicSerializer(entity entity.ActiveNewsWithTopic) NewsArticleWithTopic {
	pub := &entity.PublishedAt.Time
	if !entity.PublishedAt.Valid {
		pub = nil
	}

	return NewsArticleWithTopic{
		ID:          entity.ID,
		Title:       entity.Title,
		Content:     entity.Content,
		Slug:        entity.Slug,
		AuthorName:  entity.AuthorName,
		PublishedAt: pub,
		Topics:      append([]string(nil), entity.Topics...),
	}
}
