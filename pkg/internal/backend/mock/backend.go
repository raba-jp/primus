// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/raba-jp/primus/pkg/internal/backend (interfaces: Backend)

// Package mock_backend is a generated GoMock package.
package mock_backend

import (
	context "context"
	gomock "github.com/golang/mock/gomock"
	backend "github.com/raba-jp/primus/pkg/internal/backend"
	reflect "reflect"
)

// MockBackend is a mock of Backend interface
type MockBackend struct {
	ctrl     *gomock.Controller
	recorder *MockBackendMockRecorder
}

// MockBackendMockRecorder is the mock recorder for MockBackend
type MockBackendMockRecorder struct {
	mock *MockBackend
}

// NewMockBackend creates a new mock instance
func NewMockBackend(ctrl *gomock.Controller) *MockBackend {
	mock := &MockBackend{ctrl: ctrl}
	mock.recorder = &MockBackendMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockBackend) EXPECT() *MockBackendMockRecorder {
	return m.recorder
}

// CheckInstall mocks base method
func (m *MockBackend) CheckInstall(arg0 context.Context, arg1 string) bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckInstall", arg0, arg1)
	ret0, _ := ret[0].(bool)
	return ret0
}

// CheckInstall indicates an expected call of CheckInstall
func (mr *MockBackendMockRecorder) CheckInstall(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckInstall", reflect.TypeOf((*MockBackend)(nil).CheckInstall), arg0, arg1)
}

// Command mocks base method
func (m *MockBackend) Command(arg0 context.Context, arg1 *backend.CommandParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Command", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Command indicates an expected call of Command
func (mr *MockBackendMockRecorder) Command(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Command", reflect.TypeOf((*MockBackend)(nil).Command), arg0, arg1)
}

// FileCopy mocks base method
func (m *MockBackend) FileCopy(arg0 context.Context, arg1 *backend.FileCopyParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileCopy", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FileCopy indicates an expected call of FileCopy
func (mr *MockBackendMockRecorder) FileCopy(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileCopy", reflect.TypeOf((*MockBackend)(nil).FileCopy), arg0, arg1)
}

// FileMove mocks base method
func (m *MockBackend) FileMove(arg0 context.Context, arg1 *backend.FileMoveParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FileMove", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// FileMove indicates an expected call of FileMove
func (mr *MockBackendMockRecorder) FileMove(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FileMove", reflect.TypeOf((*MockBackend)(nil).FileMove), arg0, arg1)
}

// HTTPRequest mocks base method
func (m *MockBackend) HTTPRequest(arg0 context.Context, arg1 *backend.HTTPRequestParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "HTTPRequest", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// HTTPRequest indicates an expected call of HTTPRequest
func (mr *MockBackendMockRecorder) HTTPRequest(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "HTTPRequest", reflect.TypeOf((*MockBackend)(nil).HTTPRequest), arg0, arg1)
}

// Install mocks base method
func (m *MockBackend) Install(arg0 context.Context, arg1 *backend.InstallParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Install", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Install indicates an expected call of Install
func (mr *MockBackendMockRecorder) Install(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Install", reflect.TypeOf((*MockBackend)(nil).Install), arg0, arg1)
}

// Symlink mocks base method
func (m *MockBackend) Symlink(arg0 context.Context, arg1 *backend.SymlinkParams) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Symlink", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Symlink indicates an expected call of Symlink
func (mr *MockBackendMockRecorder) Symlink(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Symlink", reflect.TypeOf((*MockBackend)(nil).Symlink), arg0, arg1)
}

// Uninstall mocks base method
func (m *MockBackend) Uninstall(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Uninstall", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Uninstall indicates an expected call of Uninstall
func (mr *MockBackendMockRecorder) Uninstall(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Uninstall", reflect.TypeOf((*MockBackend)(nil).Uninstall), arg0, arg1)
}