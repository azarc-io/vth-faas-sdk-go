package api

import (
	"errors"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/errors"
	"testing"
)

func (j *_Job) Execute(ctx JobContext) {
	ctx.Stage("takeMoney", func(ctx StageContext) (interface{}, StageError) {
		var myData map[string]interface{}
		if err := ctx.GetVariable("data").Bind(&myData); err != nil {

		}
		return "receipt_id", nil
	}).Stage("generateEmailReceipt", func(ctx StageContext) (interface{}, StageError) {
		var myData map[string]interface{}
		if err := ctx.GetVariable("data").Bind(&myData); err != nil {

		}
		return "receipt_id", nil
	}).Stage("generatePdfReceipt", func(ctx StageContext) (interface{}, StageError) {
		var myData map[string]interface{}
		if err := ctx.GetVariable("data").Bind(&myData); err != nil {

		}
		return "receipt_id", nil
	}).Complete(func(ctx CompletionContext) StageError {
		var takeMoneyValue Stage
		var err error
		if takeMoneyValue := ctx.GetStage("takeMoney").BindValue(&takeMoneyValue); err == nil {
			if err := ctx.SetVariable("rfd", takeMoneyValue, "application/xml"); err != nil {
				return err
			}
		}
		var myValue any
		if err := ctx.GetStage("takeMoney").BindValue(&myValue).Context().SetVariable("rfd", myValue, "application/json"); err != nil {
			if errors.Is(err, sdk_errors.BindValueFailed) {
				return err
			}
		}
		ctx.SetVariable("report", "some json data", "application/json")
		ctx.SetVariable("document", "some pdf data", "pdf?")
		return nil
	}).Compensate(func(ctx CompensationContext) StageError {
		ctx.Stage("refundMoney", func(context StageContext) (any, StageError) {
			return nil, nil
		}).Stage("compensateSomethingElse", func(ctx StageContext) (interface{}, StageError) {
			return nil, nil
		})
		return nil
	}).Cancel(func(ctx CancelContext) StageError {
		return nil
	})
	println("##job executed")
}

func TestContracts(t *testing.T) {
	jobCtx := _JobContext{
		Args: JobArgs{
			JobKey:        "1",
			CorrelationId: "2",
			TransactionId: "3",
			Retries:       struct{}{},
		},
		StageChain: _StageChain{stages: make(map[string]StageDefinitionFn)},
	}
	job := _Job{jobCtx}
	job.Execute(&jobCtx)
}

type (
	JobArgs struct {
		JobKey        string
		CorrelationId string
		TransactionId string
		Retries       any
	}

	_Job struct {
		JobContext _JobContext
	}

	_JobContext struct {
		Args       JobArgs
		StageChain _StageChain
	}

	_StageChain struct {
		stages map[string]StageDefinitionFn
	}

	_CompleteChain struct {
	}

	_CompensationChain struct {
	}

	_CancelChain struct {
	}

	_CompensationContext struct{}
	_CancelContext       struct{}
	_StageContext        struct{}
	_StageVariable       struct{}
	_CompletionContext   struct{}
)

// CompletionContext

func (cc _CompletionContext) GetStage(name string) Stage {
	return nil
}

func (cc _CompletionContext) SetVariable(name string, value any, mimeType string) error {
	return nil
}

// StageVariable

func (sv _StageVariable) Raw() any {
	return struct{}{}
}

func (sv _StageVariable) Bind(any) error {
	return nil
}

// StageContext

func (sc _StageContext) GetVariable(string) StageVariable {
	return _StageVariable{}
}

// CompensationContext

func (cc _CompensationContext) Stage(string, StageDefinitionFn) StageChain {
	return &_StageChain{}
}

func (cc _CompensationContext) WithStageStatus([]string, any) bool {
	return true
}

func (cc _CompensationContext) GetVariable(s string) StageVariable {
	return _StageVariable{}
}

func (cc _CompensationContext) SetVariable(name string, value any) error {
	if name == "error" {
		return errors.New("error")
	}
	return nil
}

// CancelContext

func (cc _CancelContext) Stage(string, StageDefinitionFn) StageChain {
	return &_StageChain{}
}

// StageChain
func (sc *_StageChain) Stage(name string, sdf StageDefinitionFn) StageChain {
	sc.stages[name] = sdf
	return sc
}

func (sc *_StageChain) Complete(fn CompletionDefinitionFn) CompleteChain {
	return _CompleteChain{}
}

func (sc *_StageChain) Compensate(fn CompensateDefinitionFn) CompensateChain {
	return _CompensationChain{}
}

func (sc *_StageChain) Canceled(fn CancelDefinitionFn) CanceledChain {
	return _CancelChain{}
}

func (sc *_StageChain) Run() {
	return
}

// CompleteChain

func (_ _CompleteChain) Compensate(fn CompensateDefinitionFn) CompensateChain {
	return _CompensationChain{}
}

func (_ _CompleteChain) Cancel(fn CancelDefinitionFn) CanceledChain {
	return _CancelChain{}
}

func (_ _CompleteChain) Run() {
	return
}

// Compensation Chain

func (_ _CompensationChain) Cancel(fn CancelDefinitionFn) CanceledChain {
	return _CancelChain{}
}

func (_ _CompensationChain) Complete(fn CompletionDefinitionFn) CompleteChain {
	return _CompleteChain{}
}

func (_ _CompensationChain) Run() {
	return
}

// CancelChain

func (_ _CancelChain) Compensate(fn CompensateDefinitionFn) CompensateChain {
	return _CompensationChain{}
}

func (_ _CancelChain) Complete(fn CompletionDefinitionFn) CompleteChain {
	return _CompleteChain{}
}

func (_ _CancelChain) Run() {
	return
}

// JobContext

func (jc _JobContext) JobKey() string {
	return jc.Args.JobKey
}

func (jc _JobContext) CorrelationID() string {
	return jc.Args.CorrelationId
}

func (jc _JobContext) TransactionID() string {
	return jc.Args.TransactionId
}

func (jc _JobContext) Retries() any {
	return jc.Args.Retries
}

func (jc _JobContext) Stage(name string, sdf StageDefinitionFn) StageChain {
	jc.StageChain.Stage(name, sdf)
	return &jc.StageChain
}
