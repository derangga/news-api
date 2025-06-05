package request

type CreateNewsArticleRequest struct {
	Title    string  `json:"title" validate:"required,min=5,max=255"`
	Content  string  `json:"content" validate:"required,min=10"`
	Summary  *string `json:"summary,omitempty" validate:"omitempty,max=500"`
	AuthorID int     `json:"author_id" validate:"required,min=1"`
	Slug     string  `json:"slug" validate:"required,min=5,max=255"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=draft published deleted"`
	TopicIDs []int   `json:"topic_ids,omitempty" validate:"omitempty,dive,min=1"`
}

// UpdateNewsArticleRequest represents the request payload for updating a news article
type UpdateNewsArticleRequest struct {
	Title    *string `json:"title,omitempty" validate:"omitempty,min=5,max=255"`
	Content  *string `json:"content,omitempty" validate:"omitempty,min=10"`
	Summary  *string `json:"summary,omitempty" validate:"omitempty,max=500"`
	Slug     *string `json:"slug,omitempty" validate:"omitempty,min=5,max=255"`
	Status   *string `json:"status,omitempty" validate:"omitempty,oneof=draft published deleted"`
	TopicIDs []int32 `json:"topic_ids,omitempty" validate:"omitempty,dive,min=1"`
}
