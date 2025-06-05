package response

import (
	"newsapi/internal/model/entity"
	"time"
)

type Topic struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` // nullable
	Slug        string    `json:"slug" db:"slug"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func TopicSeriliazer(entity entity.Topic) Topic {
	return Topic{
		ID:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description,
		Slug:        entity.Slug,
		UpdatedAt:   entity.UpdatedAt,
	}
}
