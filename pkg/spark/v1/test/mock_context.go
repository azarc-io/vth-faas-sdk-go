// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1 (interfaces: Context)

// Package spark_v1_mock is a generated GoMock package.
package spark_v1_mock

import (
	context "context"
	reflect "reflect"

	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockContext is a mock of Context interface.
type MockContext struct {
	ctrl     *gomock.Controller
	recorder *MockContextMockRecorder
}

// MockContextMockRecorder is the mock recorder for MockContext.
type MockContextMockRecorder struct {
	mock *MockContext
}

// NewMockContext creates a new mock instance.
func NewMockContext(ctrl *gomock.Controller) *MockContext {
	mock := &MockContext{ctrl: ctrl}
	mock.recorder = &MockContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContext) EXPECT() *MockContextMockRecorder {
	return m.recorder
}

// CorrelationID mocks base method.
func (m *MockContext) CorrelationID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CorrelationID")
	ret0, _ := ret[0].(string)
	return ret0
}

// CorrelationID indicates an expected call of CorrelationID.
func (mr *MockContextMockRecorder) CorrelationID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CorrelationID", reflect.TypeOf((*MockContext)(nil).CorrelationID))
}

// Ctx mocks base method.
func (m *MockContext) Ctx() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ctx")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Ctx indicates an expected call of Ctx.
func (mr *MockContextMockRecorder) Ctx() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ctx", reflect.TypeOf((*MockContext)(nil).Ctx))
}

// JobKey mocks base method.
func (m *MockContext) JobKey() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "JobKey")
	ret0, _ := ret[0].(string)
	return ret0
}

// JobKey indicates an expected call of JobKey.
func (mr *MockContextMockRecorder) JobKey() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "JobKey", reflect.TypeOf((*MockContext)(nil).JobKey))
}

// LastActiveStage mocks base method.
func (m *MockContext) LastActiveStage() *sparkv1.LastActiveStage {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LastActiveStage")
	ret0, _ := ret[0].(*sparkv1.LastActiveStage)
	return ret0
}

// LastActiveStage indicates an expected call of LastActiveStage.
func (mr *MockContextMockRecorder) LastActiveStage() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LastActiveStage", reflect.TypeOf((*MockContext)(nil).LastActiveStage))
}

// TransactionID mocks base method.
func (m *MockContext) TransactionID() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TransactionID")
	ret0, _ := ret[0].(string)
	return ret0
}

// TransactionID indicates an expected call of TransactionID.
func (mr *MockContextMockRecorder) TransactionID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TransactionID", reflect.TypeOf((*MockContext)(nil).TransactionID))
}
