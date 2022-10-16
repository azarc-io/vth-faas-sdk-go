package context

import (
	ctx "context"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

type Job struct {
	ctx                  ctx.Context
	metadata             JobMetadata
	stageProgressHandler api.StageProgressHandler
	variableHandler      api.VariableHandler
}

func NewJobContext(metadata api.Context, sph api.StageProgressHandler, vh api.VariableHandler) api.JobContext {
	m := JobMetadata{ctx: metadata.Ctx(), jobKey: metadata.JobKey(), correlationId: metadata.CorrelationID(), transactionId: metadata.TransactionID(), payload: metadata.Payload(), lastActiveStage: nil} // metadata.LastActiveStage() FIXME
	return &Job{metadata: m, stageProgressHandler: sph, variableHandler: vh}
}

func (j *Job) VariableHandler() api.VariableHandler {
	return j.variableHandler
}

func (j *Job) StageProgressHandler() api.StageProgressHandler {
	return j.stageProgressHandler
}

func (j *Job) Ctx() ctx.Context {
	return j.ctx
}

func (j *Job) JobKey() string {
	return j.metadata.jobKey
}

func (j *Job) CorrelationID() string {
	return j.metadata.correlationId
}

func (j *Job) TransactionID() string {
	return j.metadata.transactionId
}

func (j *Job) Payload() any {
	return j.metadata.payload
}

func (j *Job) LastActiveStage() api.LastActiveStatus {
	return j.metadata.lastActiveStage
}

func (j *Job) SetVariables(stage string, variables ...*sdk_v1.Variable) error {
	return j.variableHandler.Set(j.metadata.jobKey, stage, variables...)
}

func (j *Job) GetVariables(stage string, names ...string) (*sdk_v1.Variables, error) {
	return j.variableHandler.Get(j.metadata.jobKey, stage, names...)
}

// TODO move
type stageOptionParams struct {
	stageName string
	sph       api.StageProgressHandler
	vh        api.VariableHandler
	ctx       api.Context
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) StageProgressHandler() api.StageProgressHandler {
	return s.sph
}

func (s stageOptionParams) VariableHandler() api.VariableHandler {
	return s.vh
}

func (s stageOptionParams) Context() api.Context {
	return s.ctx
}

func newStageOptionParams(stageName string, job *Job) api.StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		sph:       job.stageProgressHandler,
		vh:        job.variableHandler,
		ctx:       job.metadata,
	}
}

func WithStageStatus(stageName string, status sdk_v1.StageStatus) api.StageOption {
	return func(sop api.StageOptionParams) api.StageError {
		stageStatus, err := sop.StageProgressHandler().Get(sop.Context().JobKey(), stageName)
		if err != nil {
			return sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_Skip)) // TODO confirm if that error type is the right one to return here
		}
		if *stageStatus != status {
			return sdk_errors.NewStageError(fmt.Errorf("conditional stage execution skipped this stage"), sdk_errors.WithErrorType(sdk_v1.ErrorType_Skip)) // TODO confirm if that error type is the right one to return here
		}
		return nil
	}
}
