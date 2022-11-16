package sdk_v1

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"github.com/stretchr/testify/assert"
)

type InMemoryStageProgressHandler struct {
	t                  *testing.T
	stages             map[string]*SetStageStatusRequest
	results            map[string]*SetStageResultRequest
	jobs               map[string]*SetJobStatusRequest
	behaviourSet       map[string]StageBehaviourParams
	behaviourSetResult map[string]ResultBehaviourParams
}

func NewInMemoryStageProgressHandler(t *testing.T, seeds ...any) *InMemoryStageProgressHandler {
	handler := InMemoryStageProgressHandler{t,
		map[string]*SetStageStatusRequest{}, map[string]*SetStageResultRequest{},
		map[string]*SetJobStatusRequest{}, map[string]StageBehaviourParams{},
		map[string]ResultBehaviourParams{}}
	for _, seed := range seeds {
		switch seed := seed.(type) {
		case *SetStageStatusRequest:
			handler.stages[handler.key(seed.JobKey, seed.Name)] = seed
		case *SetStageResultRequest:
			handler.results[handler.key(seed.JobKey, seed.Name)] = seed
		default:
			handler.t.Fatalf("invalid seed type. accepted values are: *sdk_v1.SetStageStatusRequest, *sdk_v1.SetStageResultRequest, but got: %s", reflect.TypeOf(seed).String())
		}
	}
	return &handler
}

func (i *InMemoryStageProgressHandler) Get(jobKey, name string) (*StageStatus, error) {
	if stage, ok := i.stages[i.key(jobKey, name)]; ok {
		return &stage.Status, nil
	}
	i.t.Fatalf("stage status no found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil, nil
}

func (i *InMemoryStageProgressHandler) Set(stageStatus *SetStageStatusRequest) error {
	if bp, ok := i.behaviourSet[stageStatus.Name]; ok {
		if bp.status == stageStatus.Status && bp.err != nil {
			return bp.err
		}
	}
	i.stages[i.key(stageStatus.JobKey, stageStatus.Name)] = stageStatus
	return nil
}

func (i *InMemoryStageProgressHandler) GetResult(jobKey, name string) *Result {
	if variable, ok := i.results[i.key(jobKey, name)]; ok {
		return NewResult(nil, variable.Result)
	}
	i.t.Fatalf("stage result not found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil
}

func (i *InMemoryStageProgressHandler) SetResult(result *SetStageResultRequest) error {
	if br, ok := i.behaviourSetResult[result.Name]; ok {
		if br.jobKey == result.GetJobKey() && br.name == result.Name && br.err != nil {
			return br.err
		}
	}
	i.results[i.key(result.JobKey, result.Name)] = result
	return nil
}

func (i *InMemoryStageProgressHandler) SetJobStatus(jobStatus *SetJobStatusRequest) error {
	i.jobs[jobStatus.Key] = jobStatus
	return nil
}

func (i *InMemoryStageProgressHandler) AddBehaviour() *Behaviour {
	return &Behaviour{i: i}
}

func (i *InMemoryStageProgressHandler) ResetBehaviour() {
	i.behaviourSet = map[string]StageBehaviourParams{}
}

func (i *InMemoryStageProgressHandler) AssertStageStatus(jobKey, stageName string, expectedStatus StageStatus) {
	status, err := i.Get(jobKey, stageName)
	if err != nil {
		i.t.Error(err)
		return
	}
	assert.Equal(i.t, &expectedStatus, status, "spark status expected: '%s' got: '%s'", expectedStatus, status)
}

func (i *InMemoryStageProgressHandler) AssertStageResult(jobKey, stageName string, expectedStageResult any) {
	r := i.GetResult(jobKey, stageName)
	resB, err := r.Raw()
	if err != nil {
		i.t.Error(err)
		return
	}
	req, err := NewSetStageResultReq(jobKey, api.MimeTypeJSON, expectedStageResult)
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

type Behaviour struct {
	i *InMemoryStageProgressHandler
}

func (b *Behaviour) Set(stageName string, status StageStatus, err error) *InMemoryStageProgressHandler {
	b.i.behaviourSet[stageName] = StageBehaviourParams{err: err, status: status}
	return b.i
}

func (b *Behaviour) SetResult(jobKey, stageName string, err error) *InMemoryStageProgressHandler {
	b.i.behaviourSetResult[stageName] = ResultBehaviourParams{jobKey: jobKey, name: stageName, err: err}
	return b.i
}

type StageBehaviourParams struct {
	err    error
	status StageStatus
}

type ResultBehaviourParams struct {
	jobKey string
	name   string
	err    error
}
