package api_test

//import (
//	"errors"
//	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/errors"
//	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
//	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
//	"testing"
//)
//
//func (j *_Job) Execute(ctx api.JobContext) {
//	ctx.Stage("takeMoney", func(ctx api.StageContext) (interface{}, api.StageError) {
//		var myData map[string]interface{}
//		if v, err := ctx.GetVariable("data"); err != nil {
//			if err := v.Bind(&myData); err != nil {
//				return nil, sdk_errors.NewFromOptions(sdk_errors.WithReasonFromError(err))
//			}
//		}
//		return "receipt_id", nil
//	}).Stage("generateEmailReceipt", func(ctx api.StageContext) (interface{}, api.StageError) {
//		var myData map[string]interface{}
//		if v, err := ctx.GetVariable("data"); err != nil {
//			if err := v.Bind(&myData); err != nil {
//				return nil, sdk_errors.NewFromOptions(sdk_errors.WithReasonFromError(err))
//			}
//		}
//		return "receipt_id", nil
//	}).Stage("generatePdfReceipt", func(ctx api.StageContext) (interface{}, api.StageError) {
//		var myData map[string]interface{}
//		if v, err := ctx.GetVariable("data"); err != nil {
//			if err := v.Bind(&myData); err != nil {
//				return nil, sdk_errors.NewFromOptions(sdk_errors.WithReasonFromError(err))
//			}
//		}
//		return "receipt_id", nil
//	}).Complete(func(ctx api.CompletionContext) (any, api.StageError) {
//		var takeMoneyValue *sdk_v1.Stage
//		var err error
//		if takeMoneyValue, err = ctx.GetStage("takeMoney"); err == nil {
//			if err := ctx.SetVariable(sdk_v1.NewVariable("rfd", "application/xml", takeMoneyValue.Raw())); err != nil {
//				return nil, sdk_errors.New().Fail(err)
//			}
//		}
//		err = takeMoneyValue.Bind(&takeMoneyValue)
//		if err != nil {
//			return nil, sdk_errors.New().Fail(err)
//		}
//		err = ctx.SetVariable(sdk_v1.NewVariable("rfd", "application/xml", []byte("test")))
//		if err != nil {
//			return nil, sdk_errors.New().Fail(err)
//		}
//		return nil, nil
//	}).Compensate(func(ctx api.CompensationContext) (any, api.StageError) {
//		ctx.Stage("refundMoney", func(context api.StageContext) (any, api.StageError) {
//			return nil, nil
//		}).Stage("compensateSomethingElse", func(ctx api.StageContext) (interface{}, api.StageError) {
//			return nil, nil
//		})
//		return nil, nil
//	}).Canceled(func(ctx api.CancelContext) (any, api.StageError) {
//		return nil, nil
//	})
//	println("##job executed")
//}
//
//func TestContracts(t *testing.T) {
//	jobCtx := _JobContext{
//		Args: JobArgs{
//			JobKey:        "1",
//			CorrelationId: "2",
//			TransactionId: "3",
//			Retries:       struct{}{},
//		},
//		StageChain: _StageChain{stages: make(map[string]api.StageDefinitionFn)},
//	}
//	job := _Job{jobCtx}
//	job.Execute(&jobCtx)
//}
//
//type (
//	JobArgs struct {
//		JobKey        string
//		CorrelationId string
//		TransactionId string
//		Retries       any
//	}
//
//	_Job struct {
//		JobContext _JobContext
//	}
//
//	_JobContext struct {
//		Args       JobArgs
//		StageChain _StageChain
//	}
//
//	_StageChain struct {
//		stages map[string]api.StageDefinitionFn
//	}
//
//	_CompleteChain struct {
//	}
//
//	_CompensationChain struct {
//	}
//
//	_CancelChain struct {
//	}
//
//	_CompensationContext struct{}
//	_CancelContext       struct{}
//	_StageContext        struct{}
//	_StageVariable       struct{}
//	_CompletionContext   struct{}
//)
//
//// CompletionContext
//
//func (cc _CompletionContext) GetStage(name string) sdk_v1.Stage {
//	return sdk_v1.Stage{}
//}
//
//func (cc _CompletionContext) SetVariable(name string, value any, mimeType string) error {
//	return nil
//}
//
//// StageVariable
//
//func (sv _StageVariable) GetName() string {
//	return ""
//}
//
//func (sv _StageVariable) Raw() any {
//	return struct{}{}
//}
//
//func (sv _StageVariable) Bind(any) error {
//	return nil
//}
//
//// StageContext
//
//func (sc _StageContext) GetVariable(string) sdk_v1.Variable {
//	return sdk_v1.Variable{}
//}
//
//// CompensationContext
//
//func (cc _CompensationContext) Stage(string, api.StageDefinitionFn) api.StageChain {
//	return &_StageChain{}
//}
//
//func (cc _CompensationContext) WithStageStatus([]string, any) bool {
//	return true
//}
//
//func (cc _CompensationContext) GetVariable(s string) sdk_v1.Variable {
//	return sdk_v1.Variable{}
//}
//
//func (cc _CompensationContext) SetVariable(name string, value any) error {
//	if name == "error" {
//		return errors.New("error")
//	}
//	return nil
//}
//
//// CancelContext
//
//func (cc _CancelContext) Stage(string, api.StageDefinitionFn) api.StageChain {
//	return &_StageChain{}
//}
//
//// StageChain
//func (sc *_StageChain) Stage(name string, sdf api.StageDefinitionFn) api.StageChain {
//	sc.stages[name] = sdf
//	return sc
//}
//
//func (sc *_StageChain) Complete(fn api.CompletionDefinitionFn) api.CompleteChain {
//	return _CompleteChain{}
//}
//
//func (sc *_StageChain) Compensate(fn api.CompensateDefinitionFn) api.CompensateChain {
//	return _CompensationChain{}
//}
//
//func (sc *_StageChain) Canceled(fn api.CancelDefinitionFn) api.CanceledChain {
//	return _CancelChain{}
//}
//
//func (sc *_StageChain) Run() {
//	return
//}
//
//// CompleteChain
//
//func (_ _CompleteChain) Compensate(fn api.CompensateDefinitionFn) api.CompensateChain {
//	return _CompensationChain{}
//}
//
//func (_ _CompleteChain) Canceled(fn api.CancelDefinitionFn) api.CanceledChain {
//	return _CancelChain{}
//}
//
//func (_ _CompleteChain) Run() {
//	return
//}
//
//// Compensation Chain
//
//func (_ _CompensationChain) Canceled(fn api.CancelDefinitionFn) api.CanceledChain {
//	return _CancelChain{}
//}
//
//func (_ _CompensationChain) Complete(fn api.CompletionDefinitionFn) api.CompleteChain {
//	return _CompleteChain{}
//}
//
//func (_ _CompensationChain) Run() {
//	return
//}
//
//// CancelChain
//
//func (_ _CancelChain) Compensate(fn api.CompensateDefinitionFn) api.CompensateChain {
//	return _CompensationChain{}
//}
//
//func (_ _CancelChain) Complete(fn api.CompletionDefinitionFn) api.CompleteChain {
//	return _CompleteChain{}
//}
//
//func (_ _CancelChain) Run() {
//	return
//}
//
//// JobContext
//
//func (jc _JobContext) JobKey() string {
//	return jc.Args.JobKey
//}
//
//func (jc _JobContext) CorrelationID() string {
//	return jc.Args.CorrelationId
//}
//
//func (jc _JobContext) TransactionID() string {
//	return jc.Args.TransactionId
//}
//
//func (jc _JobContext) Retries() any {
//	return jc.Args.Retries
//}
//
//func (jc _JobContext) Stage(name string, sdf api.StageDefinitionFn) api.StageChain {
//	jc.StageChain.Stage(name, sdf)
//	return &jc.StageChain
//}
