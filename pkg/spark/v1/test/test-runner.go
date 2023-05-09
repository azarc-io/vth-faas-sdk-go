package module_test_runner

import (
	"context"
	"errors"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/google/uuid"
	"go.temporal.io/sdk/testsuite"
	"testing"
	"time"
)

var (
	ErrInvalidStageResultMimeType = errors.New("stage result expects mime-type of application/json")
)

// RunnerTest Test Helper
type RunnerTest interface {
	sparkv1.StageTracker
	Execute(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error)
}

type Outputs struct {
	sparkv1.ExecuteSparkOutput
}

func (o *Outputs) Bind(varName string, target any) error {
	err := fmt.Errorf("%w: %s", ErrNoOutput, varName)
	if o == nil || o.Outputs == nil {
		return err
	}

	if b := o.ExecuteSparkOutput.Outputs[varName]; b != nil {
		return b.Bind(target)
	}
	return err
}

type runnerTest struct {
	sparkv1.StageTracker
	sparkv1.InternalStageTracker
	ctx      sparkv1.Context
	spark    sparkv1.Spark
	testOpts *testOpts
	t        *testing.T
}

func (r *runnerTest) Execute(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error) {
	// Execute new workflow using test client
	wts := testsuite.WorkflowTestSuite{}
	env := wts.NewTestWorkflowEnvironment()
	env.SetTestTimeout(10 * time.Minute)

	// Create the spark chain
	builder := sparkv1.NewBuilder()
	r.spark.BuildChain(builder)
	chain := builder.BuildChain()

	//Initialise spark
	so := new(sparkv1.SparkOpts)
	for _, opt := range opts {
		so = opt(so)
	}
	if err := r.spark.Init(sparkv1.NewInitContext(so)); err != nil {
		return nil, err
	}

	// Create new workflow
	wf := sparkv1.NewJobWorkflow(ctx, &temporalDataProvider{provider: env}, uuid.NewString(), chain, sparkv1.WithStageTracker(r.InternalStageTracker))

	env.RegisterActivity(wf.ExecuteStageActivity)
	env.RegisterActivity(wf.ExecuteCompleteActivity)

	// check for cancel
	go func() {
		<-ctx.Done()
		if !env.IsWorkflowCompleted() {
			env.CancelWorkflow()
		}
	}()

	env.ExecuteWorkflow(wf.Run, ctx.Metadata)

	if err := env.GetWorkflowError(); err != nil {
		return nil, err
	}

	res := sparkv1.ExecuteSparkOutput{}
	if err := env.GetWorkflowResult(&res); err != nil {
		return nil, err
	}

	if res.Error != nil {
		return nil, res.Error
	}

	return &Outputs{res}, nil
}

func NewTestRunner(t *testing.T, spark sparkv1.Spark, options ...Option) (RunnerTest, error) {
	var to testOpts
	for _, option := range options {
		option(&to)
	}

	st := newStageTracker(t)
	return &runnerTest{spark: spark, testOpts: &to, InternalStageTracker: st, StageTracker: st, t: t}, nil
}

func NewTestJobContext(ctx context.Context, jobKey, correlationId, transactionId string, inputs sparkv1.ExecuteSparkInputs) *sparkv1.JobContext {
	ins := make(sparkv1.ExecuteSparkInputs)
	for name, bindable := range inputs {
		data, err := codec.Encode(bindable.Value)
		if err != nil {
			panic(err)
		}

		ins[name] = sparkv1.NewBindableValue(data, bindable.MimeType)
	}

	return &sparkv1.JobContext{
		Context: ctx,
		Metadata: &sparkv1.JobMetadata{
			JobKeyValue:        jobKey,
			CorrelationIdValue: correlationId,
			TransactionIdValue: transactionId,
			Inputs:             ins,
		},
	}
}

/*********
Temporal Data Provider which wraps the temporal testsuite.TestWorkflowEnvironment
*/

type temporalDataProvider struct {
	provider *testsuite.TestWorkflowEnvironment
}

func (tdp *temporalDataProvider) GetStageResult(workflowId, runId, stageName string) (sparkv1.Bindable, error) {
	res, err := tdp.provider.QueryWorkflow(sparkv1.JobGetStageResultQuery, stageName)
	if err != nil {
		return nil, err
	}

	var val sparkv1.Value
	if err := res.Get(&val); err != nil {
		return nil, err
	}

	if val.MimeType != "application/json" {
		return nil, ErrInvalidStageResultMimeType
	}

	return sparkv1.NewBindable(val), nil
}
