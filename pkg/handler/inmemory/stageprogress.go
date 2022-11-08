package inmemory

import (
	"fmt"
	"reflect"
	"testing"

	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"github.com/stretchr/testify/assert"
)

type StageProgressHandler struct {
	t                  *testing.T
	stages             map[string]*v1.SetStageStatusRequest
	results            map[string]*v1.SetStageResultRequest
	jobs               map[string]*v1.SetJobStatusRequest
	behaviourSet       map[string]StageBehaviourParams
	behaviourSetResult map[string]ResultBehaviourParams
}

func NewStageProgressHandler(t *testing.T, seeds ...any) *StageProgressHandler {
	handler := StageProgressHandler{t,
		map[string]*v1.SetStageStatusRequest{}, map[string]*v1.SetStageResultRequest{},
		map[string]*v1.SetJobStatusRequest{}, map[string]StageBehaviourParams{},
		map[string]ResultBehaviourParams{}}
	for _, seed := range seeds {
		switch seed := seed.(type) {
		case *v1.SetStageStatusRequest:
			handler.stages[handler.key(seed.JobKey, seed.Name)] = seed
		case *v1.SetStageResultRequest:
			handler.results[handler.key(seed.JobKey, seed.Name)] = seed
		default:
			handler.t.Fatalf("invalid seed type. accepted values are: *sdk_v1.SetStageStatusRequest, *sdk_v1.SetStageResultRequest, but got: %s", reflect.TypeOf(seed).String())
		}
	}
	return &handler
}

func (i *StageProgressHandler) Get(jobKey, name string) (*v1.StageStatus, error) {
	if stage, ok := i.stages[i.key(jobKey, name)]; ok {
		return &stage.Status, nil
	}
	i.t.Fatalf("stage status no found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil, nil
}

func (i *StageProgressHandler) Set(stageStatus *v1.SetStageStatusRequest) error {
	if bp, ok := i.behaviourSet[stageStatus.Name]; ok {
		if bp.status == stageStatus.Status && bp.err != nil {
			return bp.err
		}
	}
	i.stages[i.key(stageStatus.JobKey, stageStatus.Name)] = stageStatus
	return nil
}

func (i *StageProgressHandler) GetResult(jobKey, name string) *v1.Result {
	if variable, ok := i.results[i.key(jobKey, name)]; ok {
		return v1.NewResult(nil, variable.Result)
	}
	i.t.Fatalf("stage result no found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil
}

func (i *StageProgressHandler) SetResult(result *v1.SetStageResultRequest) error {
	if br, ok := i.behaviourSetResult[result.Name]; ok {
		if br.jobKey == result.GetJobKey() && br.name == result.Name && br.err != nil {
			return br.err
		}
	}
	i.results[i.key(result.JobKey, result.Name)] = result
	return nil
}

func (i *StageProgressHandler) SetJobStatus(jobStatus *v1.SetJobStatusRequest) error {
	i.jobs[jobStatus.Key] = jobStatus
	return nil
}

func (i *StageProgressHandler) AddBehaviour() *Behaviour {
	return &Behaviour{i: i}
}

func (i *StageProgressHandler) ResetBehaviour() {
	i.behaviourSet = map[string]StageBehaviourParams{}
}

func (i *StageProgressHandler) AssertStageStatus(jobKey, stageName string, expectedStatus v1.StageStatus) {
	status, err := i.Get(jobKey, stageName)
	if err != nil {
		i.t.Error(err)
		return
	}
	assert.Equal(i.t, expectedStatus, status, "spark status expected: '%s' got: '%s'", expectedStatus, status)
}

func (i *StageProgressHandler) AssertStageResult(jobKey, stageName string, expectedStageResult any) {
	r := i.GetResult(jobKey, stageName)
	resB, err := r.Raw()
	if err != nil {
		i.t.Error(err)
		return
	}
	req, err := v1.NewSetStageResultReq(jobKey, api.MimeTypeJSON, expectedStageResult)
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

func (i *StageProgressHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}

type Behaviour struct {
	i *StageProgressHandler
}

func (b *Behaviour) Set(stageName string, status v1.StageStatus, err error) *StageProgressHandler {
	b.i.behaviourSet[stageName] = StageBehaviourParams{err: err, status: status}
	return b.i
}

func (b *Behaviour) SetResult(jobKey, stageName string, err error) *StageProgressHandler {
	b.i.behaviourSetResult[stageName] = ResultBehaviourParams{jobKey: jobKey, name: stageName, err: err}
	return b.i
}

type StageBehaviourParams struct {
	err    error
	status v1.StageStatus
}

type ResultBehaviourParams struct {
	jobKey string
	name   string
	err    error
}
