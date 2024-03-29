// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 (interfaces: Connector)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockConnector is a mock of Connector interface.
type MockConnector struct {
	ctrl     *gomock.Controller
	recorder *MockConnectorMockRecorder
}

// MockConnectorMockRecorder is the mock recorder for MockConnector.
type MockConnectorMockRecorder struct {
	mock *MockConnector
}

// NewMockConnector creates a new mock instance.
func NewMockConnector(ctrl *gomock.Controller) *MockConnector {
	mock := &MockConnector{ctrl: ctrl}
	mock.recorder = &MockConnectorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockConnector) EXPECT() *MockConnectorMockRecorder {
	return m.recorder
}

// Start mocks base method.
func (m *MockConnector) Start(arg0 connectorv1.StartContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Start", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Start indicates an expected call of Start.
func (mr *MockConnectorMockRecorder) Start(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Start", reflect.TypeOf((*MockConnector)(nil).Start), arg0)
}

// Stop mocks base method.
func (m *MockConnector) Stop(arg0 connectorv1.StopContext) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Stop", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Stop indicates an expected call of Stop.
func (mr *MockConnectorMockRecorder) Stop(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Stop", reflect.TypeOf((*MockConnector)(nil).Stop), arg0)
}
