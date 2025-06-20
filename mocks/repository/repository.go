// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repository/repository.go

// Package mock_repository is a generated GoMock package.
package mock_repository

import (
	context "context"
	dto "newsapi/internal/model/dto"
	entity "newsapi/internal/model/entity"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUsersRepository is a mock of UsersRepository interface.
type MockUsersRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUsersRepositoryMockRecorder
}

// MockUsersRepositoryMockRecorder is the mock recorder for MockUsersRepository.
type MockUsersRepositoryMockRecorder struct {
	mock *MockUsersRepository
}

// NewMockUsersRepository creates a new mock instance.
func NewMockUsersRepository(ctrl *gomock.Controller) *MockUsersRepository {
	mock := &MockUsersRepository{ctrl: ctrl}
	mock.recorder = &MockUsersRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsersRepository) EXPECT() *MockUsersRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsersRepository) Create(ctx context.Context, entity *entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUsersRepositoryMockRecorder) Create(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsersRepository)(nil).Create), ctx, entity)
}

// MockTopicsRepository is a mock of TopicsRepository interface.
type MockTopicsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockTopicsRepositoryMockRecorder
}

// MockTopicsRepositoryMockRecorder is the mock recorder for MockTopicsRepository.
type MockTopicsRepositoryMockRecorder struct {
	mock *MockTopicsRepository
}

// NewMockTopicsRepository creates a new mock instance.
func NewMockTopicsRepository(ctrl *gomock.Controller) *MockTopicsRepository {
	mock := &MockTopicsRepository{ctrl: ctrl}
	mock.recorder = &MockTopicsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTopicsRepository) EXPECT() *MockTopicsRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockTopicsRepository) Create(ctx context.Context, entity *entity.Topic) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, entity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockTopicsRepositoryMockRecorder) Create(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockTopicsRepository)(nil).Create), ctx, entity)
}

// Delete mocks base method.
func (m *MockTopicsRepository) Delete(ctx context.Context, id int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockTopicsRepositoryMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockTopicsRepository)(nil).Delete), ctx, id)
}

// GetAll mocks base method.
func (m *MockTopicsRepository) GetAll(ctx context.Context) ([]entity.Topic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx)
	ret0, _ := ret[0].([]entity.Topic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockTopicsRepositoryMockRecorder) GetAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockTopicsRepository)(nil).GetAll), ctx)
}

// GetByID mocks base method.
func (m *MockTopicsRepository) GetByID(ctx context.Context, id int) (entity.Topic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(entity.Topic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockTopicsRepositoryMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockTopicsRepository)(nil).GetByID), ctx, id)
}

// UpdateTopicFileds mocks base method.
func (m *MockTopicsRepository) UpdateTopicFileds(ctx context.Context, topic *entity.Topic, updateFields []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateTopicFileds", ctx, topic, updateFields)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateTopicFileds indicates an expected call of UpdateTopicFileds.
func (mr *MockTopicsRepositoryMockRecorder) UpdateTopicFileds(ctx, topic, updateFields interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateTopicFileds", reflect.TypeOf((*MockTopicsRepository)(nil).UpdateTopicFileds), ctx, topic, updateFields)
}

// MockNewsArticlesRepository is a mock of NewsArticlesRepository interface.
type MockNewsArticlesRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNewsArticlesRepositoryMockRecorder
}

// MockNewsArticlesRepositoryMockRecorder is the mock recorder for MockNewsArticlesRepository.
type MockNewsArticlesRepositoryMockRecorder struct {
	mock *MockNewsArticlesRepository
}

// NewMockNewsArticlesRepository creates a new mock instance.
func NewMockNewsArticlesRepository(ctrl *gomock.Controller) *MockNewsArticlesRepository {
	mock := &MockNewsArticlesRepository{ctrl: ctrl}
	mock.recorder = &MockNewsArticlesRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNewsArticlesRepository) EXPECT() *MockNewsArticlesRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNewsArticlesRepository) Create(ctx context.Context, entity *entity.NewsArticle) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, entity)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockNewsArticlesRepositoryMockRecorder) Create(ctx, entity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNewsArticlesRepository)(nil).Create), ctx, entity)
}

// DeleteBySlug mocks base method.
func (m *MockNewsArticlesRepository) DeleteBySlug(ctx context.Context, slug string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteBySlug", ctx, slug)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteBySlug indicates an expected call of DeleteBySlug.
func (mr *MockNewsArticlesRepositoryMockRecorder) DeleteBySlug(ctx, slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteBySlug", reflect.TypeOf((*MockNewsArticlesRepository)(nil).DeleteBySlug), ctx, slug)
}

// GetActiveArticleBySlug mocks base method.
func (m *MockNewsArticlesRepository) GetActiveArticleBySlug(ctx context.Context, slug string) (entity.ActiveNewsWithTopic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetActiveArticleBySlug", ctx, slug)
	ret0, _ := ret[0].(entity.ActiveNewsWithTopic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetActiveArticleBySlug indicates an expected call of GetActiveArticleBySlug.
func (mr *MockNewsArticlesRepositoryMockRecorder) GetActiveArticleBySlug(ctx, slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetActiveArticleBySlug", reflect.TypeOf((*MockNewsArticlesRepository)(nil).GetActiveArticleBySlug), ctx, slug)
}

// GetAll mocks base method.
func (m *MockNewsArticlesRepository) GetAll(ctx context.Context, filter dto.NewsFilter) ([]entity.NewsArticleWithTopicID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, filter)
	ret0, _ := ret[0].([]entity.NewsArticleWithTopicID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockNewsArticlesRepositoryMockRecorder) GetAll(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockNewsArticlesRepository)(nil).GetAll), ctx, filter)
}

// GetArticleBySlug mocks base method.
func (m *MockNewsArticlesRepository) GetArticleBySlug(ctx context.Context, slug string) (entity.NewsArticleWithTopic, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetArticleBySlug", ctx, slug)
	ret0, _ := ret[0].(entity.NewsArticleWithTopic)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetArticleBySlug indicates an expected call of GetArticleBySlug.
func (mr *MockNewsArticlesRepositoryMockRecorder) GetArticleBySlug(ctx, slug interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetArticleBySlug", reflect.TypeOf((*MockNewsArticlesRepository)(nil).GetArticleBySlug), ctx, slug)
}

// UpdateArticleFields mocks base method.
func (m *MockNewsArticlesRepository) UpdateArticleFields(ctx context.Context, entity *entity.NewsArticleWithTopic, updateFields []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateArticleFields", ctx, entity, updateFields)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateArticleFields indicates an expected call of UpdateArticleFields.
func (mr *MockNewsArticlesRepositoryMockRecorder) UpdateArticleFields(ctx, entity, updateFields interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateArticleFields", reflect.TypeOf((*MockNewsArticlesRepository)(nil).UpdateArticleFields), ctx, entity, updateFields)
}

// MockNewsTopicsRepository is a mock of NewsTopicsRepository interface.
type MockNewsTopicsRepository struct {
	ctrl     *gomock.Controller
	recorder *MockNewsTopicsRepositoryMockRecorder
}

// MockNewsTopicsRepositoryMockRecorder is the mock recorder for MockNewsTopicsRepository.
type MockNewsTopicsRepositoryMockRecorder struct {
	mock *MockNewsTopicsRepository
}

// NewMockNewsTopicsRepository creates a new mock instance.
func NewMockNewsTopicsRepository(ctrl *gomock.Controller) *MockNewsTopicsRepository {
	mock := &MockNewsTopicsRepository{ctrl: ctrl}
	mock.recorder = &MockNewsTopicsRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockNewsTopicsRepository) EXPECT() *MockNewsTopicsRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockNewsTopicsRepository) Create(ctx context.Context, articleID int, topicIDs []int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, articleID, topicIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockNewsTopicsRepositoryMockRecorder) Create(ctx, articleID, topicIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockNewsTopicsRepository)(nil).Create), ctx, articleID, topicIDs)
}

// DeleteByArticleID mocks base method.
func (m *MockNewsTopicsRepository) DeleteByArticleID(ctx context.Context, articleID int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByArticleID", ctx, articleID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByArticleID indicates an expected call of DeleteByArticleID.
func (mr *MockNewsTopicsRepositoryMockRecorder) DeleteByArticleID(ctx, articleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByArticleID", reflect.TypeOf((*MockNewsTopicsRepository)(nil).DeleteByArticleID), ctx, articleID)
}

// ReplaceArticleTopics mocks base method.
func (m *MockNewsTopicsRepository) ReplaceArticleTopics(ctx context.Context, articleID int, topicIDs []int32) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReplaceArticleTopics", ctx, articleID, topicIDs)
	ret0, _ := ret[0].(error)
	return ret0
}

// ReplaceArticleTopics indicates an expected call of ReplaceArticleTopics.
func (mr *MockNewsTopicsRepositoryMockRecorder) ReplaceArticleTopics(ctx, articleID, topicIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReplaceArticleTopics", reflect.TypeOf((*MockNewsTopicsRepository)(nil).ReplaceArticleTopics), ctx, articleID, topicIDs)
}
