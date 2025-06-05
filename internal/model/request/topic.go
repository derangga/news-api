package request

type CreateTopicRequest struct {
	Name        string  `json:"name" validate:"required,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Slug        string  `json:"slug" validate:"required,min=2,max=100"`
}

type UpdateTopicRequest struct {
	ID          int
	Name        *string `json:"name,omitempty" validate:"omitempty,min=2,max=100"`
	Description *string `json:"description,omitempty" validate:"omitempty,max=1000"`
	Slug        *string `json:"slug,omitempty" validate:"omitempty,min=2,max=100"`
}
