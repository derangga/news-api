package repository_test

import (
	"context"
	"errors"
	"newsapi/internal/model/entity"
	"newsapi/internal/repository"
	"newsapi/internal/utils"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_CreateTopics(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewTopicRepository(sqlxDB)
	ctx := context.Background()

	query := `INSERT INTO topics \(name, description, slug\) VALUES \(\?, \?, \?\) RETURNING id`
	tests := []struct {
		testname  string
		entity    entity.Topic
		initMock  func(entity.Topic)
		assertion func(err error)
	}{
		{
			testname: "insert topic then return error",
			entity: entity.Topic{
				Name:        "test",
				Description: utils.StringPtr("test description"),
				Slug:        "slug-example",
			},
			initMock: func(t entity.Topic) {
				mockSql.ExpectPrepare(query).ExpectQuery().
					WithArgs(t.Name, t.Description, t.Slug).
					WillReturnError(errors.New("failed insert"))
			},
			assertion: func(err error) {
				assert.Error(t, err)
			},
		},
		{
			testname: "insert topic then return error nil",
			entity: entity.Topic{
				Name:        "test",
				Description: utils.StringPtr("test description"),
				Slug:        "slug-example",
			},
			initMock: func(t entity.Topic) {
				mockSql.ExpectPrepare(query).ExpectQuery().
					WithArgs(t.Name, t.Description, t.Slug).
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

func Test_GetAllTopics(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	query := `SELECT id, name, description, slug, created_at, updated_at, deleted_at 
			FROM topics 
			WHERE deleted_at IS NULL`
	now := time.Now()
	testTime := now.Truncate(time.Second)
	repos := repository.NewTopicRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testname  string
		initMock  func()
		assertion func(t []entity.Topic, err error)
	}{
		{
			testname: "get topics then return error",
			initMock: func() {
				mockSql.ExpectQuery(query).WillReturnError(errors.New("database error"))
			},
			assertion: func(tp []entity.Topic, err error) {
				assert.Error(t, err)
			},
		},
		{
			testname: "get topics then return valid topic",
			initMock: func() {
				rows := sqlmock.NewRows([]string{"id", "name", "description", "slug", "created_at", "updated_at", "deleted_at"}).
					AddRow(1, "Topic 1", "Name 1", "Description 1", testTime, testTime, nil).
					AddRow(2, "Topic 2", "Name 2", "Description 2", testTime.Add(5*time.Minute), testTime.Add(5*time.Minute), nil)

				mockSql.ExpectQuery(query).WillReturnRows(rows)
			},
			assertion: func(tp []entity.Topic, err error) {
				assert.Equal(t, 2, len(tp))
				assert.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock()
			res, err := repos.GetAll(ctx)
			tt.assertion(res, err)
		})
	}
}

func Test_UpdateTopics(t *testing.T) {
	mockDb, sqlxDB, mockSql := utils.GenerateMockDb()
	defer mockDb.Close()

	repos := repository.NewTopicRepository(sqlxDB)
	ctx := context.Background()

	tests := []struct {
		testname     string
		entity       *entity.Topic
		updateFields []string
		initMock     func(*entity.Topic)
		assertion    func(err error)
	}{
		{
			testname: "update only name field",
			entity: &entity.Topic{
				ID:   1,
				Name: "Updated Topic Name",
			},
			updateFields: []string{"name"},
			initMock: func(tp *entity.Topic) {
				expectedQuery := "UPDATE topics SET name = \\$1, updated_at = \\$2 WHERE id = \\$3"
				mockSql.ExpectExec(expectedQuery).
					WithArgs(tp.Name, sqlmock.AnyArg(), tp.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			assertion: func(err error) {
				assert.NoError(t, err)
				assert.NoError(t, mockSql.ExpectationsWereMet())
			},
		},
		{
			testname: "update name and description fields",
			entity: &entity.Topic{
				ID:          1,
				Name:        "new Name",
				Description: utils.StringPtr("new description"),
			},
			updateFields: []string{"name", "description"},
			initMock: func(tp *entity.Topic) {
				expectedQuery := "UPDATE topics SET name = \\$1, description = \\$2, updated_at = \\$3 WHERE id = \\$4"
				mockSql.ExpectExec(expectedQuery).
					WithArgs(tp.Name, tp.Description, sqlmock.AnyArg(), tp.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			assertion: func(err error) {
				assert.NoError(t, err)
				assert.NoError(t, mockSql.ExpectationsWereMet())
			},
		},
		{
			testname: "update all supported fields",
			entity: &entity.Topic{
				ID:          1,
				Name:        "all Fields Updated",
				Description: utils.StringPtr("this description is updated"),
				Slug:        "all-fields-updated",
			},
			updateFields: []string{"name", "description", "slug"},
			initMock: func(tp *entity.Topic) {
				expectedQuery := "UPDATE topics SET name = \\$1, description = \\$2, slug = \\$3, updated_at = \\$4 WHERE id = \\$5"
				mockSql.ExpectExec(expectedQuery).
					WithArgs(tp.Name, tp.Description, tp.Slug, sqlmock.AnyArg(), tp.ID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			assertion: func(err error) {
				assert.NoError(t, err)
				assert.NoError(t, mockSql.ExpectationsWereMet())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.testname, func(t *testing.T) {
			tt.initMock(tt.entity)
			err := repos.UpdateTopicFileds(ctx, tt.entity, tt.updateFields)
			tt.assertion(err)
		})
	}
}
