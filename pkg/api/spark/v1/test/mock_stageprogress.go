// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1 (interfaces: StageProgressHandler)

// Package spark_v1_mock is a generated GoMock package.
package spark_v1_mock

import (
	reflect "reflect"

	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
	gomock "github.com/golang/mock/gomock"
)

// MockStageProgressHandler is a mock of StageProgressHandler interface.
type MockStageProgressHandler struct {
	ctrl     *gomock.Controller
	recorder *MockStageProgressHandlerMockRecorder
}

// MockStageProgressHandlerMockRecorder is the mock recorder for MockStageProgressHandler.
type MockStageProgressHandlerMockRecorder struct {
	mock *MockStageProgressHandler
}

// NewMockStageProgressHandler creates a new mock instance.
func NewMockStageProgressHandler(ctrl *gomock.Controller) *MockStageProgressHandler {
	mock := &MockStageProgressHandler{ctrl: ctrl}
	mock.recorder = &MockStageProgressHandlerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStageProgressHandler) EXPECT() *MockStageProgressHandlerMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockStageProgressHandler) Get(arg0, arg1 string) (*sdk_v1.StageStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*sdk_v1.StageStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStageProgressHandlerMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStageProgressHandler)(nil).Get), arg0, arg1)
}

// GetResult mocks base method.
func (m *MockStageProgressHandler) GetResult(arg0, arg1 string) *sdk_v1.Result {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetResult", arg0, arg1)
	ret0, _ := ret[0].(*sdk_v1.Result)
	return ret0
}

// GetResult indicates an expected call of GetResult.
func (mr *MockStageProgressHandlerMockRecorder) GetResult(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetResult", reflect.TypeOf((*MockStageProgressHandler)(nil).GetResult), arg0, arg1)
}

// Set mocks base method.
func (m *MockStageProgressHandler) Set(arg0 *sdk_v1.SetStageStatusRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Set", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Set indicates an expected call of Set.
func (mr *MockStageProgressHandlerMockRecorder) Set(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Set", reflect.TypeOf((*MockStageProgressHandler)(nil).Set), arg0)
}

// SetJobStatus mocks base method.
func (m *MockStageProgressHandler) SetJobStatus(arg0 *sdk_v1.SetJobStatusRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetJobStatus", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetJobStatus indicates an expected call of SetJobStatus.
func (mr *MockStageProgressHandlerMockRecorder) SetJobStatus(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetJobStatus", reflect.TypeOf((*MockStageProgressHandler)(nil).SetJobStatus), arg0)
}

// SetResult mocks base method.
func (m *MockStageProgressHandler) SetResult(arg0 *sdk_v1.SetStageResultRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetResult", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetResult indicates an expected call of SetResult.
func (mr *MockStageProgressHandlerMockRecorder) SetResult(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetResult", reflect.TypeOf((*MockStageProgressHandler)(nil).SetResult), arg0)
}