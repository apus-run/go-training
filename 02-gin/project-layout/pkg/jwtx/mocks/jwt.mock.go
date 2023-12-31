// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/jwtx/types.go

// Package jwtxmocks is a generated GoMock package.
package jwtxmocks

import (
	jwtx "project-layout/pkg/jwtx"
	reflect "reflect"

	jwt "github.com/golang-jwt/jwt/v5"
	gomock "go.uber.org/mock/gomock"
)

// MockJwtToken is a mock of JwtToken interface.
type MockJwtToken struct {
	ctrl     *gomock.Controller
	recorder *MockJwtTokenMockRecorder
}

// MockJwtTokenMockRecorder is the mock recorder for MockJwtToken.
type MockJwtTokenMockRecorder struct {
	mock *MockJwtToken
}

// NewMockJwtToken creates a new mock instance.
func NewMockJwtToken(ctrl *gomock.Controller) *MockJwtToken {
	mock := &MockJwtToken{ctrl: ctrl}
	mock.recorder = &MockJwtTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJwtToken) EXPECT() *MockJwtTokenMockRecorder {
	return m.recorder
}

// GenerateToken mocks base method.
func (m *MockJwtToken) GenerateToken(options ...jwtx.Option) (string, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{}
	for _, a := range options {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GenerateToken", varargs...)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateToken indicates an expected call of GenerateToken.
func (mr *MockJwtTokenMockRecorder) GenerateToken(options ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateToken", reflect.TypeOf((*MockJwtToken)(nil).GenerateToken), options...)
}

// ParseToken mocks base method.
func (m *MockJwtToken) ParseToken(tokenString, secretKey string) (*jwtx.CustomClaims, *jwt.Token, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseToken", tokenString, secretKey)
	ret0, _ := ret[0].(*jwtx.CustomClaims)
	ret1, _ := ret[1].(*jwt.Token)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// ParseToken indicates an expected call of ParseToken.
func (mr *MockJwtTokenMockRecorder) ParseToken(tokenString, secretKey interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseToken", reflect.TypeOf((*MockJwtToken)(nil).ParseToken), tokenString, secretKey)
}
