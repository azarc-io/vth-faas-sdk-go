// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/azarc-io/vth-faas-sdk-go/pkg/api (interfaces: VariableHandler)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockVariableHandler is a mock of VariableHandler interface.
type MockVariableHandler struct {
	ctrl     *gomock.Controller
	recorder *MockVariableHandlerMockRecorder
}

// MockVariableHandlerMockRecorder is the mock recorder for MockVariableHandler.
type MockVariableHandlerMockRecorder struct {
	mock *MockVariableHandler
}

// NewMockVariableHandler creates a new mock instance.
func NewMockVariableHandler(ctrl *gomock.Controller) *MockVariableHandler {
	mock := &MockVariableHandler{ctrl: ctrl}
	mock.recorder = &MockVariableHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVariableHandler) EXPECT() *MockVariableHandlerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockVariableHandler) Get(arg0, arg1 string, arg2 ...string) ([]*sdk_v1.Variable, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Get", varargs...)
	ret0, _ := ret[0].([]*sdk_v1.Variable)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockVariableHandlerMockRecorder) Get(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockVariableHandler)(nil).Get), varargs...)
}

// Set mocks base method.
func (m *MockVariableHandler) Set(arg0, arg1 string, arg2 ...*sdk_v1.Variable) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Set", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockVariableHandlerMockRecorder) Set(arg0, arg1 interface{}, arg2 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockVariableHandler)(nil).Set), varargs...)
}