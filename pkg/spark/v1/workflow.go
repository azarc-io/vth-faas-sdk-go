package sparkv1

import (
	"context"
	"errors"
	"fmt"
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

func (w *jobWorkflow) Run(ctx workflow.Context, jmd *JobMetadata) (*ExecuteSparkOutput, error) {
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

	var doNext func(next *Node) *ExecuteSparkOutput

	doNext = func(next *Node) *ExecuteSparkOutput {
		if next == nil {
			return nil
		}

		for _, stage := range next.Stages {
			w.setStageStatus(stage.Name, StageStatus_STAGE_STARTED)
			res, err := w.executeStageActivity(ctx, stage.Name, info, state)
			if err != nil {
				w.setStageStatus(stage.Name, StageStatus_STAGE_FAILED)
				return getSparkErrorOutput(err)
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
				return getSparkErrorOutput(err)
			}

			w.setStageStatus(next.Complete.Name, StageStatus_STAGE_COMPLETED)
			return v
		}

		return getSparkErrorOutput(module_runner2.ErrChainDoesNotHaveACompleteStage)
	}

	return doNext(w.Chain.RootNode), nil
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
			StageName:     stageName,
			Inputs:        state.JobContext.Inputs,
			WorkflowId:    info.WorkflowExecution.ID,
			RunId:         info.WorkflowExecution.RunID,
			JobKey:        state.JobContext.JobKeyValue,
			TransactionId: state.JobContext.TransactionIdValue,
			CorrelationId: state.JobContext.CorrelationIdValue,
		})
		if err := f.Get(ctx, &sr); err != nil {
			//TODO Compensate()
			return nil, err
		}

		if codec.MimeType(sr.MimeType) == codec.MimeTypeJson.WithType("error") {
			attempts++
			se := errorWrap{StageName: stageName}
			if err := sr.Bind(&se); err != nil {
				return nil, err
			}

			if se.Retry == nil {
				// no retries set, exit
				return nil, se
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

func (w *jobWorkflow) executeCompleteActivity(ctx workflow.Context, stageName string, info *workflow.Info, state *JobState) (*ExecuteSparkOutput, error) {
	var res ExecuteSparkOutput // stage result

	options := DefaultActivityOptions.GetTemporalActivityOptions()
	options.ActivityID = stageName
	options.RetryPolicy = DefaultRetryPolicy.GetTemporalPolicy()

	c := workflow.WithActivityOptions(ctx, options)
	f := workflow.ExecuteActivity(c, w.ExecuteCompleteActivity, &ExecuteStageRequest{
		StageName:     stageName,
		Inputs:        state.JobContext.Inputs,
		WorkflowId:    info.WorkflowExecution.ID,
		RunId:         info.WorkflowExecution.RunID,
		JobKey:        state.JobContext.JobKeyValue,
		TransactionId: state.JobContext.TransactionIdValue,
		CorrelationId: state.JobContext.CorrelationIdValue,
	})
	if err := f.Get(ctx, &res); err != nil {
		// TODO ctx.Compensate()
		return nil, err
	}

	return &res, nil
}

func (w *jobWorkflow) ExecuteStageActivity(ctx context.Context, req *ExecuteStageRequest) (Bindable, StageError) {
	fn := w.Chain.GetStageFunc(req.StageName)
	sc := NewStageContext(req, w.sparkDataIO, req.WorkflowId, req.RunId, req.StageName, NewLogger(), req.Inputs)

	var err StageError
	out := w.executeFn(func() (any, StageError) {
		return fn(sc)
	}, &err)
	if err != nil {
		return getTransferableError(err.(error)), nil
	}

	stageValue, err2 := codec.Encode(out)
	if err2 != nil {
		return getTransferableError(err2), nil
	}

	return &bindable{Value: stageValue, MimeType: string(codec.MimeTypeJson)}, nil
}

func (w *jobWorkflow) ExecuteCompleteActivity(ctx context.Context, req *ExecuteStageRequest) (*ExecuteSparkOutput, StageError) {
	fn := w.Chain.GetStageCompleteFunc(req.StageName)
	cc := NewCompleteContext(req, w.sparkDataIO, req.WorkflowId, req.RunId, req.StageName, NewLogger(), req.Inputs)

	var err StageError
	_ = w.executeFn(func() (any, StageError) {
		err = fn(cc)
		return nil, err
	}, &err)
	if err != nil {
		return &ExecuteSparkOutput{
			Error: &ExecuteSparkError{
				StageName:    req.StageName,
				ErrorCode:    err.ErrorCode(),
				ErrorMessage: err.Error(),
				Metadata:     err.Metadata(),
				StackTrace:   getStackTrace(err),
			},
		}, nil
	}

	res := &ExecuteSparkOutput{
		Outputs: map[string]*bindable{},
	}
	for _, output := range cc.(*completeContext).outputs {
		res.Outputs[output.Name] = &bindable{Value: output.Value, MimeType: string(output.MimeType)}
	}
	return res, err
}

func (w *jobWorkflow) executeFn(executor func() (any, StageError), se *StageError) any {
	var v any
	defer func() {
		if err := recover(); err != nil {
			switch ec := err.(type) {
			case error:
				*se = NewStageErrorWithCode(errorCodeInternal, ec)
			default:
				*se = NewStageErrorWithCode(errorCodeInternal, errors.New(fmt.Sprint(ec)))
			}
		}
	}()

	// execute stage
	v, *se = executor()

	return v
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
	StageName    string           `json:"stage_name,omitempty"`
	ErrorCode    ErrorCode        `json:"error_code"`
	ErrorMessage string           `json:"error_message,omitempty"`
	Metadata     map[string]any   `json:"metadata,omitempty"`
	Retry        *RetryConfig     `json:"retry,omitempty"`
	StackTrace   []StackTraceItem `json:"stack_trace"`
}

func (e errorWrap) Error() string {
	return e.ErrorMessage
}

func getTransferableError(err error) Bindable {
	var ew []byte
	if se, ok := err.(StageError); ok {
		ew, _ = codec.Encode(errorWrap{
			StageName:    se.StageName(),
			ErrorMessage: se.Error(),
			ErrorCode:    se.ErrorCode(),
			Metadata:     se.Metadata(),
			Retry:        se.GetRetryConfig(),
			StackTrace:   getStackTrace(se),
		})
	} else {
		ew, _ = codec.Encode(errorWrap{
			ErrorMessage: err.Error(),
			ErrorCode:    errorCodeInternal,
		})
	}

	return NewBindable(Value{Value: ew, MimeType: string(MimeJsonError)})
}

func getSparkErrorOutput(err error) *ExecuteSparkOutput {
	if e, ok := err.(errorWrap); ok {
		return &ExecuteSparkOutput{
			Error: &ExecuteSparkError{
				StageName:    e.StageName,
				ErrorCode:    e.ErrorCode,
				ErrorMessage: e.ErrorMessage,
				Metadata:     e.Metadata,
				StackTrace:   e.StackTrace,
			},
		}
	}

	var stackTrace []StackTraceItem
	if st, ok := err.(stackTracer); ok {
		stackTrace = getStackTrace(st)
	}

	return &ExecuteSparkOutput{
		Error: &ExecuteSparkError{
			ErrorMessage: err.Error(),
			StackTrace:   stackTrace,
			ErrorCode:    errorCodeInternal,
		},
	}
}
