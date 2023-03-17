// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 (interfaces: Forwarder)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockForwarder is a mock of Forwarder interface.
type MockForwarder struct {
	ctrl     *gomock.Controller
	recorder *MockForwarderMockRecorder
}

// MockForwarderMockRecorder is the mock recorder for MockForwarder.
type MockForwarderMockRecorder struct {
	mock *MockForwarder
}

// NewMockForwarder creates a new mock instance.
func NewMockForwarder(ctrl *gomock.Controller) *MockForwarder {
	mock := &MockForwarder{ctrl: ctrl}
	mock.recorder = &MockForwarderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockForwarder) EXPECT() *MockForwarderMockRecorder {
	return m.recorder
}

// Forward mocks base method.
func (m *MockForwarder) Forward(arg0 string, arg1 []byte, arg2 map[string]string) (connectorv1.InboundResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Forward", arg0, arg1, arg2)
	ret0, _ := ret[0].(connectorv1.InboundResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Forward indicates an expected call of Forward.
func (mr *MockForwarderMockRecorder) Forward(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Forward", reflect.TypeOf((*MockForwarder)(nil).Forward), arg0, arg1, arg2)
}
