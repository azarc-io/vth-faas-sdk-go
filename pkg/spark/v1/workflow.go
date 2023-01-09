package sparkv1

import (
	"context"
	module_runner2 "github.com/azarc-io/vth-faas-sdk-go/internal/module-runner"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/workflow"
	"time"
)

const JobGetStageResultQuery = "GET_STAGE_RESULT"

type RetryConfig struct {
	Times             uint          `json:"times" yaml:"times"`
	FirstBackoffWait  time.Duration `json:"first_backoff_wait" yaml:"first_backoff_wait"`
	BackoffMultiplier uint          `json:"backoff_multiplier" yaml:"backoff_multiplier"`
}

type jobWorkflow struct {
	Chain        *SparkChain
	SparkId      string
	ctx          context.Context
	sparkDataIO  SparkDataIO
	stageTracker InternalStageTracker
}

func (w *jobWorkflow) Run(ctx workflow.Context, jmd *JobMetadata) (*JobOutput, error) {
	state := &JobState{
		JobContext: jmd,
	}

	info := workflow.GetInfo(ctx)

	// handles queries to fetch a stage result
	err := workflow.SetQueryHandler(ctx, JobGetStageResultQuery, func(stageName string) (interface{}, error) {
		if v, ok := state.StageResults[stageName]; ok {
			return v, nil
		}

		return nil, module_runner2.ErrStageResultNotFound
	})

	if err != nil {
		return nil, err
	}

	var doNext func(next *Node) (executeSparkOutputs, error)

	doNext = func(next *Node) (executeSparkOutputs, error) {
		if next == nil {
			return nil, nil
		}

		for _, stage := range next.Stages {
			w.setStageStatus(stage.Name, StageStatus_STAGE_STARTED)
			res, err := w.executeStageActivity(ctx, stage.Name, info, state)
			if err != nil {
				w.setStageStatus(stage.Name, StageStatus_STAGE_FAILED)
				return nil, err
			}

			w.setStageStatus(stage.Name, StageStatus_STAGE_COMPLETED)

			if state.StageResults == nil {
				state.StageResults = make(map[string]Bindable)
			}

			state.StageResults[stage.Name] = res
			if w.stageTracker != nil {
				w.stageTracker.SetStageResult(stage.Name, res)
			}
		}

		if next.Complete != nil {
			w.setStageStatus(next.Complete.Name, StageStatus_STAGE_STARTED)
			v, err := w.executeCompleteActivity(ctx, next.Complete.Name, info, state)
			if err != nil {
				w.setStageStatus(next.Complete.Name, StageStatus_STAGE_FAILED)
				return nil, err
			}

			w.setStageStatus(next.Complete.Name, StageStatus_STAGE_COMPLETED)
			return v, err
		}

		return nil, module_runner2.ErrChainDoesNotHaveACompleteStage
	}

	output, err := doNext(w.Chain.RootNode)
	if err != nil {
		return nil, err
	}

	return &JobOutput{
		Outputs: output,
	}, nil
}

func (w *jobWorkflow) executeStageActivity(ctx workflow.Context, stageName string, info *workflow.Info, state *JobState) (Bindable, error) {
	var sr bindable // stage result

	options := DefaultActivityOptions.GetTemporalActivityOptions()
	options.ActivityID = stageName
	options.RetryPolicy = DefaultRetryPolicy.GetTemporalPolicy()

	var attempts uint = 0
	var waitTime *time.Duration
	for {
		c := workflow.WithActivityOptions(ctx, options)
		f := workflow.ExecuteActivity(c, w.ExecuteStageActivity, &ExecuteStageRequest{
			StageName:  stageName,
			Inputs:     state.JobContext.Inputs,
			WorkflowId: info.WorkflowExecution.ID,
			RunId:      info.WorkflowExecution.RunID,
		})
		if err := f.Get(ctx, &sr); err != nil {
			//TODO Compensate()
			return nil, err
		}

		if codec.MimeType(sr.MimeType) == codec.MimeTypeJson.WithType("error") {
			attempts++
			var se errorWrap
			if err := sr.Bind(&se); err != nil {
				return nil, err
			}

			if se.Retry.Times <= attempts {
				//TODO Compensate()
				return nil, se
			}

			if waitTime == nil {
				waitTime = &se.Retry.FirstBackoffWait
			} else {
				d := time.Duration(int64(*waitTime) * int64(se.Retry.BackoffMultiplier))
				waitTime = &d
			}

			log.Info().Msgf("stage error occurred, sleeping %s before retry attempt %d", waitTime, attempts)
			if err := workflow.Sleep(ctx, *waitTime); err != nil {
				return nil, err
			}
		} else {
			return &sr, nil
		}
	}
}

func (w *jobWorkflow) executeCompleteActivity(ctx workflow.Context, stageName string, info *workflow.Info, state *JobState) (executeSparkOutputs, error) {
	var sr executeSparkOutputs // stage result

	options := DefaultActivityOptions.GetTemporalActivityOptions()
	options.ActivityID = stageName
	options.RetryPolicy = DefaultRetryPolicy.GetTemporalPolicy()

	c := workflow.WithActivityOptions(ctx, options)
	f := workflow.ExecuteActivity(c, w.ExecuteCompleteActivity, &ExecuteStageRequest{
		StageName:  stageName,
		Inputs:     state.JobContext.Inputs,
		WorkflowId: info.WorkflowExecution.ID,
		RunId:      info.WorkflowExecution.RunID,
	})
	if err := f.Get(ctx, &sr); err != nil {
		// TODO ctx.Compensate()
		return nil, err
	}

	return sr, nil
}

func (w *jobWorkflow) ExecuteStageActivity(ctx context.Context, req *ExecuteStageRequest) (Bindable, StageError) {
	fn := w.Chain.GetStageFunc(req.StageName)
	sc := NewStageContext(req, w.sparkDataIO, req.WorkflowId, req.RunId, req.StageName, NewLogger(), req.Inputs)
	out, err := fn(sc)

	if err != nil {
		if se, ok := err.(StageError); ok {
			ew, _ := codec.Encode(errorWrap{
				ErrorMessage: se.Error(),
				Metadata:     se.Metadata(),
				Retry:        se.GetRetryConfig(),
			}, codec.MimeTypeJson)
			return NewBindable(Value{Value: ew, MimeType: string(codec.MimeTypeJson.WithType("error"))}), nil
		}

		return nil, err
	}

	stageValue, err2 := codec.Encode(out, codec.MimeTypeJson)
	if err2 != nil {
		return nil, NewStageError(err2)
	}

	return &bindable{Value: stageValue, Raw: out, MimeType: string(codec.MimeTypeJson)}, err
}

func (w *jobWorkflow) ExecuteCompleteActivity(ctx context.Context, req *ExecuteStageRequest) (executeSparkOutputs, StageError) {
	fn := w.Chain.GetStageCompleteFunc(req.StageName)
	cc := NewCompleteContext(req, w.sparkDataIO, req.WorkflowId, req.RunId, req.StageName, NewLogger(), req.Inputs)
	err := fn(cc)
	if err != nil {
		return nil, err
	}

	outputs := make(executeSparkOutputs)
	for _, output := range cc.(*completeContext).outputs {
		val, err := codec.Encode(output.Value, output.MimeType)
		if err != nil {
			return nil, NewStageError(err)
		}
		outputs[output.Name] = &bindable{Value: val, Raw: output.Value, MimeType: string(output.MimeType)}
	}
	return outputs, err
}

func (w *jobWorkflow) setStageStatus(name string, status StageStatus) {
	if w.stageTracker != nil {
		w.stageTracker.SetStageStatus(name, status)
	}
}

func NewJobWorkflow(ctx context.Context, sparkDataIO SparkDataIO, sparkId string, chain *SparkChain, opts ...WorkflowOption) JobWorkflow {
	wo := new(workflowOpts)
	for _, opt := range opts {
		wo = opt(wo)
	}
	return &jobWorkflow{ctx: ctx, sparkDataIO: sparkDataIO, SparkId: sparkId, Chain: chain, stageTracker: wo.stageTracker}
}

// errorWrap used to marshal errors between workflow and activities
type errorWrap struct {
	ErrorMessage string
	Metadata     map[string]any
	Retry        *RetryConfig
}

func (e errorWrap) Error() string {
	return e.ErrorMessage
}
