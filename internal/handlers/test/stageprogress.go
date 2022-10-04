package test

import (
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type inMemoryStageProgressHandler struct {
	stages  map[string]*sdk_v1.Stage
	results map[string]*sdk_v1.StageResult
	jobs    map[string]*sdk_v1.Job
}

// TODO add t *testing.T
// TODO add options to the constructor
func NewMockStageProgressHandler(stages ...*sdk_v1.Stage) api.StageProgressHandler {
	handler := inMemoryStageProgressHandler{map[string]*sdk_v1.Stage{}, map[string]*sdk_v1.StageResult{}, map[string]*sdk_v1.Job{}}
	if stages == nil {
		return handler
	}
	for _, stage := range stages {
		handler.stages[stage.Name] = stage
	}
	return handler
}

func (i inMemoryStageProgressHandler) Get(name string) (*sdk_v1.Stage, error) {
	if variable, ok := i.stages[name]; ok {
		return variable, nil
	}
	return nil, sdk_errors.VariableNotFound
}

func (i inMemoryStageProgressHandler) Set(stage *sdk_v1.Stage) error {
	i.stages[stage.Name] = stage
	return nil
}

func (i inMemoryStageProgressHandler) GetResult(stage *sdk_v1.Stage) (*sdk_v1.StageResult, error) {
	return i.results[stage.Name], nil
}

func (i inMemoryStageProgressHandler) SetResult(result *sdk_v1.StageResult) error {
	i.results[result.Stage.Name] = result
	return nil
}

func (i inMemoryStageProgressHandler) GetJob(jobKey string) (*sdk_v1.Job, error) {
	return i.jobs[jobKey], nil
}

func (i inMemoryStageProgressHandler) SetJobStatus(jobKey string, status sdk_v1.JobStatus) error {
	if job, ok := i.jobs[jobKey]; ok {
		job.Status = status
	}
	return nil
}
