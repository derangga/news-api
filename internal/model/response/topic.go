package response

import "time"

type Topic struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"` // nullable
	Slug        string    `json:"slug" db:"slug"`
	UpdatedAt   time.Time `json:"updated_at"`
}
