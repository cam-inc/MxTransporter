// +build test
//

// Code generated by MockGen. DO NOT EDIT.
// Source: ./usecases/resume-token/index.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockResumeToken is a mock of ResumeToken interface.
type MockResumeToken struct {
	ctrl     *gomock.Controller
	recorder *MockResumeTokenMockRecorder
}

// MockResumeTokenMockRecorder is the mock recorder for MockResumeToken.
type MockResumeTokenMockRecorder struct {
	mock *MockResumeToken
}

// NewMockResumeToken creates a new mock instance.
func NewMockResumeToken(ctrl *gomock.Controller) *MockResumeToken {
	mock := &MockResumeToken{ctrl: ctrl}
	mock.recorder = &MockResumeTokenMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockResumeToken) EXPECT() *MockResumeTokenMockRecorder {
	return m.recorder
}

// Env mocks base method.
func (m *MockResumeToken) Env() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Env")
	ret0, _ := ret[0].(string)
	return ret0
}

// Env indicates an expected call of Env.
func (mr *MockResumeTokenMockRecorder) Env() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Env", reflect.TypeOf((*MockResumeToken)(nil).Env))
}

// ReadResumeToken mocks base method.
func (m *MockResumeToken) ReadResumeToken(ctx context.Context) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ReadResumeToken", ctx)
	ret0, _ := ret[0].(string)
	return ret0
}

// ReadResumeToken indicates an expected call of ReadResumeToken.
func (mr *MockResumeTokenMockRecorder) ReadResumeToken(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ReadResumeToken", reflect.TypeOf((*MockResumeToken)(nil).ReadResumeToken), ctx)
}

// SaveResumeToken mocks base method.
func (m *MockResumeToken) SaveResumeToken(ctx context.Context, rt string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveResumeToken", ctx, rt)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveResumeToken indicates an expected call of SaveResumeToken.
func (mr *MockResumeTokenMockRecorder) SaveResumeToken(ctx, rt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveResumeToken", reflect.TypeOf((*MockResumeToken)(nil).SaveResumeToken), ctx, rt)
}
