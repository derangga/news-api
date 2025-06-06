package utils

import (
	"database/sql"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func GenerateMockDb() (*sql.DB, *sqlx.DB, sqlmock.Sqlmock) {
	mockDb, mockSql, _ := sqlmock.New()
	sqlxDB := sqlx.NewDb(mockDb, "sqlmock")
	return mockDb, sqlxDB, mockSql
}

func StringPtr(s string) *string {
	return &s
}
