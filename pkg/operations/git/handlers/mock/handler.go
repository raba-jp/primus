// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/raba-jp/primus/pkg/operations/git/handlers (interfaces: CloneHandler)

// Package mock_handlers is a generated GoMock package.
package mock_handlers

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	handlers "github.com/raba-jp/primus/pkg/operations/git/handlers"
	reflect "reflect"
)

// MockCloneHandler is a mock of CloneHandler interface
type MockCloneHandler struct {
	ctrl     *gomock.Controller
	recorder *MockCloneHandlerMockRecorder
}

// MockCloneHandlerMockRecorder is the mock recorder for MockCloneHandler
type MockCloneHandlerMockRecorder struct {
	mock *MockCloneHandler
}

// NewMockCloneHandler creates a new mock instance
func NewMockCloneHandler(ctrl *gomock.Controller) *MockCloneHandler {
	mock := &MockCloneHandler{ctrl: ctrl}
	mock.recorder = &MockCloneHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockCloneHandler) EXPECT() *MockCloneHandlerMockRecorder {
	return m.recorder
}

// Clone mocks base method
func (m *MockCloneHandler) Clone(arg0 context.Context, arg1 bool, arg2 *handlers.CloneParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Clone", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// Clone indicates an expected call of Clone
func (mr *MockCloneHandlerMockRecorder) Clone(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Clone", reflect.TypeOf((*MockCloneHandler)(nil).Clone), arg0, arg1, arg2)
}