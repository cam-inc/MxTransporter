// +build test

// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/client/storage/storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockStorageClient is a mock of StorageClient interface.
type MockStorageClient struct {
	ctrl     *gomock.Controller
	recorder *MockStorageClientMockRecorder
}

// MockStorageClientMockRecorder is the mock recorder for MockStorageClient.
type MockStorageClientMockRecorder struct {
	mock *MockStorageClient
}

// NewMockStorageClient creates a new mock instance.
func NewMockStorageClient(ctrl *gomock.Controller) *MockStorageClient {
	mock := &MockStorageClient{ctrl: ctrl}
	mock.recorder = &MockStorageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorageClient) EXPECT() *MockStorageClientMockRecorder {
	return m.recorder
}

// DeleteObject mocks base method.
func (m *MockStorageClient) DeleteObject(ctx context.Context, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteObject", ctx, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteObject indicates an expected call of DeleteObject.
func (mr *MockStorageClientMockRecorder) DeleteObject(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteObject", reflect.TypeOf((*MockStorageClient)(nil).DeleteObject), ctx, key)
}

// GetObject mocks base method.
func (m *MockStorageClient) GetObject(ctx context.Context, key string) ([]byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetObject", ctx, key)
	ret0, _ := ret[0].([]byte)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetObject indicates an expected call of GetObject.
func (mr *MockStorageClientMockRecorder) GetObject(ctx, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetObject", reflect.TypeOf((*MockStorageClient)(nil).GetObject), ctx, key)
}

// PutObject mocks base method.
func (m *MockStorageClient) PutObject(ctx context.Context, key, value string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PutObject", ctx, key, value)
	ret0, _ := ret[0].(error)
	return ret0
}

// PutObject indicates an expected call of PutObject.
func (mr *MockStorageClientMockRecorder) PutObject(ctx, key, value interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PutObject", reflect.TypeOf((*MockStorageClient)(nil).PutObject), ctx, key, value)
}
