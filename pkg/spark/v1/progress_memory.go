package spark_v1

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/internal/common"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"testing"

	"github.com/stretchr/testify/assert"
)

type InMemoryStageProgressHandler struct {
	t                  *testing.T
	stages             map[string]*sparkv1.SetStageStatusRequest
	results            map[string]*sparkv1.SetStageResultRequest
	resultOrder        map[string][]string
	behaviourSet       map[string]StageBehaviourParams
	behaviourSetResult map[string]ResultBehaviourParams
	finished           bool
}

func (i *InMemoryStageProgressHandler) FinishJob(ctx SparkContext, req *sparkv1.FinishJobRequest) error {
	i.finished = true
	return nil
}

func (i *InMemoryStageProgressHandler) JobStarting(result *sparkv1.JobStartingRequest) error {
	// not needed for in memory
	return nil
}

func NewInMemoryStageProgressHandler(t *testing.T, seeds ...any) TestStageProgressHandler {
	handler := InMemoryStageProgressHandler{t,
		map[string]*sparkv1.SetStageStatusRequest{},
		map[string]*sparkv1.SetStageResultRequest{},
		make(map[string][]string),
		map[string]StageBehaviourParams{},
		map[string]ResultBehaviourParams{},
		false,
	}

	return &handler
}

func (i *InMemoryStageProgressHandler) Get(ctx SparkContext, name string) (*sparkv1.StageStatus, error) {
	if stage, ok := i.stages[i.key(ctx, name)]; ok {
		return &stage.Status, nil
	}
	i.t.Fatalf("stage status no found for params >> jobKey: %s, stageName: %s", ctx.JobKey(), name)
	return nil, nil
}

func (i *InMemoryStageProgressHandler) Set(ctx SparkContext, stageStatus *sparkv1.SetStageStatusRequest) error {
	if bp, ok := i.behaviourSet[stageStatus.Name]; ok {
		if bp.status == stageStatus.Status && bp.err != nil {
			return bp.err
		}
	}
	i.stages[i.key(ctx, stageStatus.Name)] = stageStatus

	if stageStatus.Status == sparkv1.StageStatus_STAGE_STARTED {
		i.resultOrder[stageStatus.Key] = append(i.resultOrder[stageStatus.Key], stageStatus.Name)
	}
	return nil
}

func (i *InMemoryStageProgressHandler) GetResult(ctx SparkContext, name string) Bindable {
	if variable, ok := i.results[i.key(ctx, name)]; ok {
		return newResult(nil, &sparkv1.GetStageResultResponse{
			Data: variable.Data,
		})
	}
	i.t.Fatalf("stage result not found for params >> jobKey: %s, stageName: %s", ctx.JobKey(), name)
	return nil
}

func (i *InMemoryStageProgressHandler) SetResult(ctx SparkContext, result *sparkv1.SetStageResultRequest) error {
	if br, ok := i.behaviourSetResult[result.Name]; ok {
		if br.jobKey == result.GetKey() && br.name == result.Name && br.err != nil {
			return br.err
		}
	}
	i.results[i.key(ctx, result.Name)] = result
	return nil
}

func (i *InMemoryStageProgressHandler) AddBehaviour() *Behaviour {
	return &Behaviour{i: i}
}

func (i *InMemoryStageProgressHandler) ResetBehaviour() {
	i.behaviourSet = map[string]StageBehaviourParams{}
}

func (i *InMemoryStageProgressHandler) AssertJobFinished(value bool) {
	assert.Equal(i.t, value, i.finished)
}

func (i *InMemoryStageProgressHandler) AssertStageCompleted(ctx SparkContext, stageName string) {
	i.assertStageStatus(ctx, stageName, sparkv1.StageStatus_STAGE_COMPLETED)
}

func (i *InMemoryStageProgressHandler) AssertStageStarted(ctx SparkContext, stageName string) {
	i.assertStageStatus(ctx, stageName, sparkv1.StageStatus_STAGE_STARTED)
}

func (i *InMemoryStageProgressHandler) AssertStageSkipped(ctx SparkContext, stageName string) {
	i.assertStageStatus(ctx, stageName, sparkv1.StageStatus_STAGE_SKIPPED)
}

func (i *InMemoryStageProgressHandler) AssertStageCancelled(ctx SparkContext, stageName string) {
	i.assertStageStatus(ctx, stageName, sparkv1.StageStatus_STAGE_CANCELED)
}

func (i *InMemoryStageProgressHandler) AssertStageFailed(ctx SparkContext, stageName string) {
	i.assertStageStatus(ctx, stageName, sparkv1.StageStatus_STAGE_FAILED)
}

func (i *InMemoryStageProgressHandler) AssertStageResult(ctx SparkContext, stageName string, expectedStageResult any) {
	r := i.GetResult(ctx, stageName)
	resB, err := r.Raw()
	if err != nil {
		i.t.Error(err)
		return
	}
	req, err := newSetStageResultReq(ctx, common.MimeTypeJSON, expectedStageResult)
	if err != nil {
		i.t.Error(err)
		return
	}
	assert.Equal(i.t, req.Data, resB)
}

func (i *InMemoryStageProgressHandler) AssertStageOrder(ctx SparkContext, stageNames ...string) {
	sns := i.resultOrder[ctx.JobKey()]

	if len(stageNames) > len(sns) {
		i.t.Fatalf("more stage names provided than were executed")
		return
	}

	actual := make([]string, len(stageNames))
	for ind := range stageNames {
		actual[ind] = sns[ind]
	}
	assert.Equal(i.t, stageNames, actual, fmt.Sprintf("actual stages: %v", sns))
}

func (i *InMemoryStageProgressHandler) key(ctx SparkContext, name string) string {
	return fmt.Sprintf("%s_%s", ctx.JobKey(), name)
}

func (i *InMemoryStageProgressHandler) assertStageStatus(ctx SparkContext, stageName string, expectedStatus sparkv1.StageStatus) {
	status, err := i.Get(ctx, stageName)
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
