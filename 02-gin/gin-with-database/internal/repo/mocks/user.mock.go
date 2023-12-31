// Code generated by MockGen. DO NOT EDIT.
// Source: ./repo/user.go

// Package repomocks is a generated GoMock package.
package repomocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"

	"gin-with-database/internal/domain/entity"
)

// MockUserRepository is a mock of UserRepository interface.
type MockUserRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryMockRecorder
}

// MockUserRepositoryMockRecorder is the mock recorder for MockUserRepository.
type MockUserRepositoryMockRecorder struct {
	mock *MockUserRepository
}

// NewMockUserRepository creates a new mock instance.
func NewMockUserRepository(ctrl *gomock.Controller) *MockUserRepository {
	mock := &MockUserRepository{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepository) EXPECT() *MockUserRepositoryMockRecorder {
	return m.recorder
}

// FindByEmail mocks base method.
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByEmail", ctx, email)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByEmail indicates an expected call of FindByEmail.
func (mr *MockUserRepositoryMockRecorder) FindByEmail(ctx, email interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByEmail", reflect.TypeOf((*MockUserRepository)(nil).FindByEmail), ctx, email)
}

// FindByID mocks base method.
func (m *MockUserRepository) FindByID(ctx context.Context, id uint64) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByID", ctx, id)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByID indicates an expected call of FindByID.
func (mr *MockUserRepositoryMockRecorder) FindByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByID", reflect.TypeOf((*MockUserRepository)(nil).FindByID), ctx, id)
}

// FindByPhone mocks base method.
func (m *MockUserRepository) FindByPhone(ctx context.Context, phone string) (*entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindByPhone", ctx, phone)
	ret0, _ := ret[0].(*entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindByPhone indicates an expected call of FindByPhone.
func (mr *MockUserRepositoryMockRecorder) FindByPhone(ctx, phone interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindByPhone", reflect.TypeOf((*MockUserRepository)(nil).FindByPhone), ctx, phone)
}

// FindUserPage mocks base method.
func (m *MockUserRepository) FindUserPage(ctx context.Context, name string, page, size int64) ([]*entity.User, uint, bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserPage", ctx, name, page, size)
	ret0, _ := ret[0].([]*entity.User)
	ret1, _ := ret[1].(uint)
	ret2, _ := ret[2].(bool)
	ret3, _ := ret[3].(error)
	return ret0, ret1, ret2, ret3
}

// FindUserPage indicates an expected call of FindUserPage.
func (mr *MockUserRepositoryMockRecorder) FindUserPage(ctx, name, page, size interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserPage", reflect.TypeOf((*MockUserRepository)(nil).FindUserPage), ctx, name, page, size)
}

// Remove mocks base method.
func (m *MockUserRepository) Remove(ctx context.Context, userEntity entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", ctx, userEntity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockUserRepositoryMockRecorder) Remove(ctx, userEntity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockUserRepository)(nil).Remove), ctx, userEntity)
}

// Save mocks base method.
func (m *MockUserRepository) Save(ctx context.Context, userEntity entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, userEntity)
	ret0, _ := ret[0].(error)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockUserRepositoryMockRecorder) Save(ctx, userEntity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockUserRepository)(nil).Save), ctx, userEntity)
}

// SaveAndCache mocks base method.
func (m *MockUserRepository) SaveAndCache(ctx context.Context, userEntity entity.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveAndCache", ctx, userEntity)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveAndCache indicates an expected call of SaveAndCache.
func (mr *MockUserRepositoryMockRecorder) SaveAndCache(ctx, userEntity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveAndCache", reflect.TypeOf((*MockUserRepository)(nil).SaveAndCache), ctx, userEntity)
}
