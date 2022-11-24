package spark_v1

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

type InMemoryStageProgressHandler struct {
	t                  *testing.T
	stages             map[string]*sparkv1.SetStageStatusRequest
	results            map[string]*sparkv1.SetStageResultRequest
	jobs               map[string]*sparkv1.SetJobStatusRequest
	behaviourSet       map[string]StageBehaviourParams
	behaviourSetResult map[string]ResultBehaviourParams
}

func NewInMemoryStageProgressHandler(t *testing.T, seeds ...any) TestStageProgressHandler {
	handler := InMemoryStageProgressHandler{t,
		map[string]*sparkv1.SetStageStatusRequest{}, map[string]*sparkv1.SetStageResultRequest{},
		map[string]*sparkv1.SetJobStatusRequest{}, map[string]StageBehaviourParams{},
		map[string]ResultBehaviourParams{}}
	for _, seed := range seeds {
		switch seed := seed.(type) {
		case *sparkv1.SetStageStatusRequest:
			handler.stages[handler.key(seed.JobKey, seed.Name)] = seed
		case *sparkv1.SetStageResultRequest:
			handler.results[handler.key(seed.JobKey, seed.Name)] = seed
		default:
			handler.t.Fatalf("invalid seed type. accepted values are: *sdk_v1.SetStageStatusRequest, "+
				"*sdk_v1.SetStageResultRequest, but got: %s", reflect.TypeOf(seed).String())
		}
	}
	return &handler
}

func (i *InMemoryStageProgressHandler) Get(jobKey, name string) (*sparkv1.StageStatus, error) {
	if stage, ok := i.stages[i.key(jobKey, name)]; ok {
		return &stage.Status, nil
	}
	i.t.Fatalf("stage status no found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil, nil
}

func (i *InMemoryStageProgressHandler) Set(stageStatus *sparkv1.SetStageStatusRequest) error {
	if bp, ok := i.behaviourSet[stageStatus.Name]; ok {
		if bp.status == stageStatus.Status && bp.err != nil {
			return bp.err
		}
	}
	i.stages[i.key(stageStatus.JobKey, stageStatus.Name)] = stageStatus
	return nil
}

func (i *InMemoryStageProgressHandler) GetResult(jobKey, name string) Bindable {
	if variable, ok := i.results[i.key(jobKey, name)]; ok {
		return newResult(nil, variable.Result)
	}
	i.t.Fatalf("stage result not found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil
}

func (i *InMemoryStageProgressHandler) SetResult(result *sparkv1.SetStageResultRequest) error {
	if br, ok := i.behaviourSetResult[result.Name]; ok {
		if br.jobKey == result.GetJobKey() && br.name == result.Name && br.err != nil {
			return br.err
		}
	}
	i.results[i.key(result.JobKey, result.Name)] = result
	return nil
}

func (i *InMemoryStageProgressHandler) SetJobStatus(jobStatus *sparkv1.SetJobStatusRequest) error {
	i.jobs[jobStatus.Key] = jobStatus
	return nil
}

func (i *InMemoryStageProgressHandler) AddBehaviour() *Behaviour {
	return &Behaviour{i: i}
}

func (i *InMemoryStageProgressHandler) ResetBehaviour() {
	i.behaviourSet = map[string]StageBehaviourParams{}
}

func (i *InMemoryStageProgressHandler) AssertStageCompleted(jobKey, stageName string) {
	i.assertStageStatus(jobKey, stageName, sparkv1.StageStatus_STAGE_STATUS_COMPLETED)
}

func (i *InMemoryStageProgressHandler) AssertStageStarted(jobKey, stageName string) {
	i.assertStageStatus(jobKey, stageName, sparkv1.StageStatus_STAGE_STATUS_STARTED)
}

func (i *InMemoryStageProgressHandler) AssertStageSkipped(jobKey, stageName string) {
	i.assertStageStatus(jobKey, stageName, sparkv1.StageStatus_STAGE_STATUS_SKIPPED)
}

func (i *InMemoryStageProgressHandler) AssertStageCancelled(jobKey, stageName string) {
	i.assertStageStatus(jobKey, stageName, sparkv1.StageStatus_STAGE_STATUS_CANCELLED)
}

func (i *InMemoryStageProgressHandler) AssertStageFailed(jobKey, stageName string) {
	i.assertStageStatus(jobKey, stageName, sparkv1.StageStatus_STAGE_STATUS_FAILED)
}

func (i *InMemoryStageProgressHandler) AssertStageUnspecified(jobKey, stageName string) {
	i.assertStageStatus(jobKey, stageName, sparkv1.StageStatus_STAGE_STATUS_PENDING_UNSPECIFIED)
}

func (i *InMemoryStageProgressHandler) AssertStageResult(jobKey, stageName string, expectedStageResult any) {
	r := i.GetResult(jobKey, stageName)
	resB, err := r.Raw()
	if err != nil {
		i.t.Error(err)
		return
	}
	req, err := newSetStageResultReq(jobKey, MimeTypeJSON, expectedStageResult)
	if err != nil {
		i.t.Error(err)
		return
	}
	reqB, err := req.Result.Raw()
	if err != nil {
		i.t.Error(err)
		return
	}
	assert.Equal(i.t, reqB, resB)
}

func (i *InMemoryStageProgressHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}

func (i *InMemoryStageProgressHandler) assertStageStatus(jobKey, stageName string, expectedStatus sparkv1.StageStatus) {
	status, err := i.Get(jobKey, stageName)
	if err != nil {
		i.t.Error(err)
		return
	}
	assert.Equal(i.t, &expectedStatus, status, "spark status expected: '%s' got: '%s'", expectedStatus, status)
}

type Behaviour struct {
	i *InMemoryStageProgressHandler
}

func (b *Behaviour) Set(stageName string, status sparkv1.StageStatus, err error) *InMemoryStageProgressHandler {
	b.i.behaviourSet[stageName] = StageBehaviourParams{err: err, status: status}
	return b.i
}

func (b *Behaviour) SetResult(jobKey, stageName string, err error) *InMemoryStageProgressHandler {
	b.i.behaviourSetResult[stageName] = ResultBehaviourParams{jobKey: jobKey, name: stageName, err: err}
	return b.i
}

type StageBehaviourParams struct {
	err    error
	status sparkv1.StageStatus
}

type ResultBehaviourParams struct {
	jobKey string
	name   string
	err    error
}
