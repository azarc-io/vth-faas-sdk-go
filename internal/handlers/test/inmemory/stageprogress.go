package inmemory

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"reflect"
	"testing"
)

type inMemoryStageProgressHandler struct {
	t       *testing.T
	stages  map[string]*sdk_v1.SetStageStatusRequest
	results map[string]*sdk_v1.SetStageResultRequest
	jobs    map[string]*sdk_v1.SetJobStatusRequest
}

func NewStageProgressHandler(t *testing.T, seeds ...any) api.StageProgressHandler {
	handler := inMemoryStageProgressHandler{t, map[string]*sdk_v1.SetStageStatusRequest{}, map[string]*sdk_v1.SetStageResultRequest{}, map[string]*sdk_v1.SetJobStatusRequest{}}
	for _, seed := range seeds {
		switch seed.(type) {
		case *sdk_v1.SetStageStatusRequest:
			s := seed.(*sdk_v1.SetStageStatusRequest)
			handler.stages[handler.key(s.JobKey, s.Name)] = s
		case *sdk_v1.SetStageResultRequest:
			r := seed.(*sdk_v1.SetStageResultRequest)
			handler.results[handler.key(r.JobKey, r.Name)] = r
		default:
			handler.t.Fatalf("invalid seed type. accepted values are: *sdk_v1.SetStageStatusRequest, *sdk_v1.SetStageResultRequest, but got: %s", reflect.TypeOf(seed).String())
		}
	}
	return handler
}

func (i inMemoryStageProgressHandler) Get(jobKey, name string) (*sdk_v1.StageStatus, error) {
	if stage, ok := i.stages[i.key(jobKey, name)]; ok {
		return &stage.Status, nil
	}
	i.t.Fatalf("stage status no found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil, nil
}

func (i inMemoryStageProgressHandler) Set(stageStatus *sdk_v1.SetStageStatusRequest) error {
	i.stages[i.key(stageStatus.JobKey, stageStatus.Name)] = stageStatus
	return nil
}

func (i inMemoryStageProgressHandler) GetResult(jobKey, name string) (*sdk_v1.StageResult, error) {
	if variable, ok := i.results[i.key(jobKey, name)]; ok {
		return variable.Result, nil
	}
	i.t.Fatalf("stage result no found for params >> jobKey: %s, stageName: %s", jobKey, name)
	return nil, nil
}

func (i inMemoryStageProgressHandler) SetResult(result *sdk_v1.SetStageResultRequest) error {
	i.results[i.key(result.JobKey, result.Name)] = result
	return nil
}

func (i inMemoryStageProgressHandler) SetJobStatus(jobStatus *sdk_v1.SetJobStatusRequest) error {
	i.jobs[jobStatus.Key] = jobStatus
	return nil
}

func (i inMemoryStageProgressHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}
