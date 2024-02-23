package sparkv1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	module_runner2 "github.com/azarc-io/vth-faas-sdk-go/internal/module-runner"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/cenkalti/backoff/v4"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"runtime"
	"time"
)

const JobGetStageResultQuery = "GET_STAGE_RESULT"

type RetryConfig struct {
	Times             uint          `json:"times" yaml:"times"`
	FirstBackoffWait  time.Duration `json:"first_backoff_wait" yaml:"first_backoff_wait"`
	BackoffMultiplier uint          `json:"backoff_multiplier" yaml:"backoff_multiplier"`
}

type sparkOutput struct {
	Outputs map[string]*BindableValue `json:"outputs,omitempty"`
	Error   *ExecuteSparkError        `json:"error,omitempty"`
}

type jobWorkflow struct {
	Chain        *SparkChain
	SparkId      string
	ctx          context.Context
	stageTracker InternalStageTracker
	cfg          *Config
	nc           *nats.Conn
}

func (w *jobWorkflow) Run(msg *nats.Msg) {
	var jmd *JobMetadata
	if err := json.Unmarshal(msg.Data, &jmd); err != nil {
		w.publishError(err)
		return
	}

	state := &JobState{
		JobContext: jmd,
	}

	var sparkIO SparkDataIO = &ioDataProvider{
		ctx:          w.ctx,
		bucket:       w.cfg.NatsBucket,
		nc:           w.nc,
		stageResults: map[string]*BindableValue{},
	}

	var doNext func(next *Node) *ExecuteSparkOutput

	doNext = func(next *Node) *ExecuteSparkOutput {
		if next == nil {
			return nil
		}

		for _, stage := range next.Stages {
			select {
			case <-w.ctx.Done():
				w.setStageStatus(stage.Name, StageStatus_STAGE_FAILED)
				return getSparkErrorOutput(errors.New("canceled"))
			default:
				w.setStageStatus(stage.Name, StageStatus_STAGE_STARTED)
				res, err := w.executeStageActivity(w.ctx, stage.Name, state, sparkIO)
				if err != nil {
					w.setStageStatus(stage.Name, StageStatus_STAGE_FAILED)
					return getSparkErrorOutput(err)
				}

				select {
				case <-w.ctx.Done():
					w.setStageStatus(stage.Name, StageStatus_STAGE_FAILED)
				default:
					w.setStageStatus(stage.Name, StageStatus_STAGE_COMPLETED)

					if state.StageResults == nil {
						state.StageResults = make(map[string]Bindable)
					}

					state.StageResults[stage.Name] = res
					if w.stageTracker != nil {
						w.stageTracker.SetStageResult(stage.Name, res)
					}
				}
			}
		}

		if next.Complete != nil {
			w.setStageStatus(next.Complete.Name, StageStatus_STAGE_STARTED)
			v, err := w.executeCompleteActivity(w.ctx, next.Complete.Name, state, sparkIO)
			if err != nil {
				w.setStageStatus(next.Complete.Name, StageStatus_STAGE_FAILED)
				return getSparkErrorOutput(err)
			}

			w.setStageStatus(next.Complete.Name, StageStatus_STAGE_COMPLETED)
			return v
		}

		return getSparkErrorOutput(module_runner2.ErrChainDoesNotHaveACompleteStage)
	}

	result := doNext(w.Chain.RootNode)
	b, err := json.Marshal(result)
	if err != nil {
		w.publishError(err)
		return
	}

	w.publish(b)
}

func (w *jobWorkflow) executeStageActivity(ctx context.Context, stageName string, state *JobState, io SparkDataIO) (Bindable, error) {
	var (
		sr  Bindable // stage result
		err error
	)

	var attempts uint = 0
	var waitTime *time.Duration
	for {
		sr, err = w.ExecuteStageActivity(ctx, &ExecuteStageRequest{
			StageName:     stageName,
			Inputs:        state.JobContext.Inputs,
			JobKey:        state.JobContext.JobKeyValue,
			TransactionId: state.JobContext.TransactionIdValue,
			CorrelationId: state.JobContext.CorrelationIdValue,
		}, io)
		if err != nil {
			return nil, err
		}

		if codec.MimeType(sr.GetMimeType()) == codec.MimeTypeJson.WithType("error") {
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
			runtime.Gosched()
			time.Sleep(*waitTime)
		} else {
			return sr, nil
		}
	}
}

func (w *jobWorkflow) executeCompleteActivity(ctx context.Context, stageName string, state *JobState, io SparkDataIO) (*ExecuteSparkOutput, error) {
	var (
		out *ExecuteSparkOutput
		err error
	)

	out, err = w.ExecuteCompleteActivity(ctx, &ExecuteStageRequest{
		StageName:     stageName,
		Inputs:        state.JobContext.Inputs,
		JobKey:        state.JobContext.JobKeyValue,
		TransactionId: state.JobContext.TransactionIdValue,
		CorrelationId: state.JobContext.CorrelationIdValue,
	}, io)
	if err != nil {
		return nil, err
	}

	return out, nil
}

func (w *jobWorkflow) ExecuteStageActivity(ctx context.Context, req *ExecuteStageRequest, io SparkDataIO) (Bindable, StageError) {
	fn := w.Chain.GetStageFunc(req.StageName)

	newInputs := make(map[string]Bindable)
	for k, v := range req.Inputs {
		newInputs[k] = io.NewInput(req.CorrelationId, v)
	}

	sc := NewStageContext(ctx, req, io, req.StageName, NewLogger(), newInputs)

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

	res, err2 := io.PutStageResult(req.StageName, stageValue)
	if err2 != nil {
		return getTransferableError(err2), nil
	}
	return res, nil
}

func (w *jobWorkflow) ExecuteCompleteActivity(ctx context.Context, req *ExecuteStageRequest, io SparkDataIO) (*ExecuteSparkOutput, StageError) {
	fn := w.Chain.GetStageCompleteFunc(req.StageName)

	newInputs := make(map[string]Bindable)
	for k, v := range req.Inputs {
		newInputs[k] = io.NewInput(req.CorrelationId, v)
	}

	cc := NewCompleteContext(ctx, req, io, req.StageName, NewLogger(), newInputs)

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
		Outputs: map[string]Bindable{},
	}
	for _, output := range cc.(*completeContext).outputs {
		var err error
		if res.Outputs[output.Name], err = io.NewOutput(req.StageName, &BindableValue{Value: output.Value, MimeType: string(output.MimeType)}); err != nil {
			return nil, NewStageError(fmt.Errorf("error occured setting output: %w", err))
		}
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

func (w *jobWorkflow) publish(b []byte) {
	pb := backoff.NewExponentialBackOff()
	if err := backoff.Retry(func() error {
		if err := w.nc.Publish(w.cfg.NatsResponseSubject, b); err != nil {
			log.Error().Err(err).Msgf("failed to publish result to broker, will retry")
			return err
		}
		return nil
	}, backoff.WithContext(pb, context.Background())); err != nil {
		log.Error().Err(err).Msgf("failed to publish result to broker, result is lost")
	}
}

func (w *jobWorkflow) publishError(err error) {
	result := getSparkErrorOutput(err)
	b, err := json.Marshal(result)
	if err != nil {
		log.Error().Err(err).Msgf("spark errored but could not marshal error response, result will be lost")
	}
	w.publish(b)
}

func NewJobWorkflow(ctx context.Context, sparkId string, chain *SparkChain, opts ...WorkflowOption) JobWorkflow {
	wo := new(workflowOpts)
	for _, opt := range opts {
		wo = opt(wo)
	}
	return &jobWorkflow{
		ctx:          ctx,
		SparkId:      sparkId,
		Chain:        chain,
		stageTracker: wo.stageTracker,
		cfg:          wo.config,
		nc:           wo.nc,
	}
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
