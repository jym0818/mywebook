// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/user.go
//
// Generated by this command:
//
//	mockgen -source=./internal/service/user.go -package=svcmocks -destination=./internal/service/mocks/user.mock.go
//

// Package svcmocks is a generated GoMock package.
package svcmocks

import (
	context "context"
	reflect "reflect"

	domain "github.com/jym/mywebook/internal/domain"
	gomock "go.uber.org/mock/gomock"
)

// MockUserService is a mock of UserService interface.
type MockUserService struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceMockRecorder
	isgomock struct{}
}

// MockUserServiceMockRecorder is the mock recorder for MockUserService.
type MockUserServiceMockRecorder struct {
	mock *MockUserService
}

// NewMockUserService creates a new mock instance.
func NewMockUserService(ctrl *gomock.Controller) *MockUserService {
	mock := &MockUserService{ctrl: ctrl}
	mock.recorder = &MockUserServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserService) EXPECT() *MockUserServiceMockRecorder {
	return m.recorder
}

// FindOrCreate mocks base method.
func (m *MockUserService) FindOrCreate(ctx context.Context, phone string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindOrCreate", ctx, phone)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindOrCreate indicates an expected call of FindOrCreate.
func (mr *MockUserServiceMockRecorder) FindOrCreate(ctx, phone any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindOrCreate", reflect.TypeOf((*MockUserService)(nil).FindOrCreate), ctx, phone)
}

// Login mocks base method.
func (m *MockUserService) Login(ctx context.Context, email, password string) (domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Login", ctx, email, password)
	ret0, _ := ret[0].(domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Login indicates an expected call of Login.
func (mr *MockUserServiceMockRecorder) Login(ctx, email, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Login", reflect.TypeOf((*MockUserService)(nil).Login), ctx, email, password)
}

// Signup mocks base method.
func (m *MockUserService) Signup(ctx context.Context, user domain.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Signup", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// Signup indicates an expected call of Signup.
func (mr *MockUserServiceMockRecorder) Signup(ctx, user any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Signup", reflect.TypeOf((*MockUserService)(nil).Signup), ctx, user)
}
