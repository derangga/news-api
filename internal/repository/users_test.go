package repository_test

import (
	"context"
	"errors"
	"newsapi/internal/model/entity"
	"newsapi/internal/repository"
	"newsapi/internal/utils"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateUser(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewUsersRepository(sqlxDB)
	ctx := context.Background()

	query := `INSERT INTO users \(name, email\) VALUES \(\?, \?\) RETURNING id`
	tests := []struct {
		testname  string
		entity    entity.User
		initMock  func(entity.User)
		assertion func(err error)
	}{
		{
			testname: "insert user then return error",
			entity: entity.User{
				Name:  "test",
				Email: "john@example.com",
			},
			initMock: func(t entity.User) {
				mockSql.ExpectPrepare(query).ExpectQuery().
					WithArgs(t.Name, t.Email).
					WillReturnError(errors.New("failed insert"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testname: "insert user then return error nil",
			entity: entity.User{
				Name:  "test",
				Email: "john@example.com",
			},
			initMock: func(t entity.User) {
				mockSql.ExpectPrepare(query).ExpectQuery().
					WithArgs(t.Name, t.Email).
					WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
			},
			assertion: func(err error) {
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.entity)
			err := repos.Create(ctx, &tt.entity)
			tt.assertion(err)
		})
	}
}
