package entity

import (
	"database/sql"
	"time"
)

type Topic struct {
	ID          int          `db:"id"`
	Name        string       `db:"name"`
	Description *string      `db:"description"`
	Slug        string       `db:"slug"`
	CreatedAt   time.Time    `db:"created_at"`
	UpdatedAt   time.Time    `db:"updated_at"`
	DeletedAt   sql.NullTime `db:"deleted_at"`
}
