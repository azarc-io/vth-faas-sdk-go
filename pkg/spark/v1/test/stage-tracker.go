package module_test_runner

import (
	"errors"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"reflect"
	"testing"

	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/stretchr/testify/assert"
)

var (
	ErrNoStageResult = errors.New("stage result not found")
	ErrNoOutput      = errors.New("unable to find output")
)

type stageTracker struct {
	sparkv1.StageTracker
	sparkv1.InternalStageTracker
	results     map[string]*result
	t           *testing.T
	resultOrder []string
}

type result struct {
	value       sparkv1.Bindable
	stageStatus sparkv1.StageStatus
}

func (st *stageTracker) GetStageResult(name string) (data any, mime codec.MimeType, err sparkv1.StageError) {
	res, ok := st.results[name]
	if !ok {
		return nil, "", sparkv1.NewStageError(fmt.Errorf("%w: %s", ErrNoStageResult, name))
	}

	raw, err2 := res.value.GetValue()
	if err2 != nil {
		return nil, "", sparkv1.NewStageError(err2)
	}

	//TODO Remove
	//switch v := raw.(type) {
	//case sparkv1.StageError:
	//	return v, codec.MimeTypeJson.WithType("error"), nil
	//}
	return raw, codec.MimeType(res.value.GetMimeType()), nil
}

func (st *stageTracker) SetStageStatus(name string, status sparkv1.StageStatus) {
	res, ok := st.results[name]
	if !ok {
		res = &result{}
		st.results[name] = res
	}

	res.stageStatus = status
}

func (st *stageTracker) SetStageResult(name string, val sparkv1.Bindable) {
	res, ok := st.results[name]
	if !ok {
		res = &result{}
		st.results[name] = res
	}

	res.value = val
	st.resultOrder = append(st.resultOrder, name)
}

func (st *stageTracker) AssertStageCompleted(stageName string) {
	st.assertStageStatus(stageName, sparkv1.StageStatus_STAGE_COMPLETED)
}

func (st *stageTracker) AssertStageStarted(stageName string) {
	st.assertStageStatus(stageName, sparkv1.StageStatus_STAGE_STARTED)
}

func (st *stageTracker) AssertStageSkipped(stageName string) {
	st.assertStageStatus(stageName, sparkv1.StageStatus_STAGE_SKIPPED)
}

func (st *stageTracker) AssertStageCancelled(stageName string) {
	st.assertStageStatus(stageName, sparkv1.StageStatus_STAGE_CANCELED)
}

func (st *stageTracker) AssertStageFailed(stageName string) {
	st.assertStageStatus(stageName, sparkv1.StageStatus_STAGE_FAILED)
}

func (st *stageTracker) AssertStageResult(stageName string, expectedStageResult any) {
	res, ok := st.results[stageName]
	if !ok {
		st.t.Error(fmt.Errorf("%w: %s", ErrNoStageResult, stageName))
	}

	// create new instance of the expected value so
	newValPtr := reflect.New(reflect.TypeOf(expectedStageResult))
	newVal := newValPtr.Elem().Interface()
	err := res.value.Bind(&newVal)
	if !assert.NoError(st.t, err) {
		return
	}

	assert.Equal(st.t, expectedStageResult, newVal)
}

func (st *stageTracker) AssertStageOrder(stageNames ...string) {
	if len(stageNames) > len(st.resultOrder) {
		st.t.Fatalf("more stage names provided than were executed")
		return
	}

	actual := make([]string, len(stageNames))
	for ind := range stageNames {
		actual[ind] = st.resultOrder[ind]
	}
	assert.Equal(st.t, stageNames, actual, fmt.Sprintf("actual stages: %v", st.resultOrder))
}

func (st *stageTracker) assertStageStatus(stageName string, expectedStatus sparkv1.StageStatus) {
	res, ok := st.results[stageName]
	if !ok {
		st.t.Error(fmt.Errorf("%w: %s", ErrNoStageResult, stageName))
		return
	}

	assert.Equal(st.t, expectedStatus, res.stageStatus, "spark status expected: '%s' got: '%s'", expectedStatus, res.stageStatus)
}

func newStageTracker(t *testing.T) *stageTracker {
	return &stageTracker{
		t:       t,
		results: make(map[string]*result),
	}
}
