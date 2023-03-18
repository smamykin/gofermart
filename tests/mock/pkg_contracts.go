// Code generated by MockGen. DO NOT EDIT.
// Source: ./pkg/contracts/logger_interface.go

// Package mock_contracts is a generated GoMock package.
package mock_contracts

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockLoggerInterface is a mock of LoggerInterface interface.
type MockLoggerInterface struct {
	ctrl     *gomock.Controller
	recorder *MockLoggerInterfaceMockRecorder
}

// MockLoggerInterfaceMockRecorder is the mock recorder for MockLoggerInterface.
type MockLoggerInterfaceMockRecorder struct {
	mock *MockLoggerInterface
}

// NewMockLoggerInterface creates a new mock instance.
func NewMockLoggerInterface(ctrl *gomock.Controller) *MockLoggerInterface {
	mock := &MockLoggerInterface{ctrl: ctrl}
	mock.recorder = &MockLoggerInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoggerInterface) EXPECT() *MockLoggerInterfaceMockRecorder {
	return m.recorder
}

// Debug mocks base method.
func (m *MockLoggerInterface) Debug(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Debug", msg)
}

// Debug indicates an expected call of Debug.
func (mr *MockLoggerInterfaceMockRecorder) Debug(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Debug", reflect.TypeOf((*MockLoggerInterface)(nil).Debug), msg)
}

// Err mocks base method.
func (m *MockLoggerInterface) Err(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Err", err)
}

// Err indicates an expected call of Err.
func (mr *MockLoggerInterfaceMockRecorder) Err(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Err", reflect.TypeOf((*MockLoggerInterface)(nil).Err), err)
}

// Fatal mocks base method.
func (m *MockLoggerInterface) Fatal(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Fatal", err)
}

// Fatal indicates an expected call of Fatal.
func (mr *MockLoggerInterfaceMockRecorder) Fatal(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Fatal", reflect.TypeOf((*MockLoggerInterface)(nil).Fatal), err)
}

// Info mocks base method.
func (m *MockLoggerInterface) Info(msg string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Info", msg)
}

// Info indicates an expected call of Info.
func (mr *MockLoggerInterfaceMockRecorder) Info(msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Info", reflect.TypeOf((*MockLoggerInterface)(nil).Info), msg)
}

// Warn mocks base method.
func (m *MockLoggerInterface) Warn(err error) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Warn", err)
}

// Warn indicates an expected call of Warn.
func (mr *MockLoggerInterfaceMockRecorder) Warn(err interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Warn", reflect.TypeOf((*MockLoggerInterface)(nil).Warn), err)
}