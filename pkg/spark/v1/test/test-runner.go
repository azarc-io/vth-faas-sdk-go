package module_test_runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1/util"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
)

var (
	ErrInvalidStageResultMimeType = errors.New("stage result expects mime-type of application/json")
)

// RunnerTest Test Helper
type RunnerTest interface {
	sparkv1.StageTracker
	Execute(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error)
	ExecuteWithoutStageRetryOverride(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error)
}

type Input struct {
	Value    any
	MimeType codec.MimeType
}
type Inputs map[string]*Input

type Outputs struct {
	sparkv1.ExecuteSparkOutput
}

func (o *Outputs) Bind(varName string, target any) error {
	err := fmt.Errorf("%w: %s", ErrNoOutput, varName)
	if o == nil || o.Outputs == nil {
		return err
	}

	if b := o.Outputs[varName]; b != nil {
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

type runnerTestOutput struct {
	VariablesKey  string                            `json:"variables_key,omitempty"`
	JobKey        string                            `json:"job_key,omitempty"`
	CorrelationId string                            `json:"correlation_id,omitempty"`
	TransactionId string                            `json:"transaction_id,omitempty"`
	Model         string                            `json:"model,omitempty"`
	Outputs       map[string]*sparkv1.BindableValue `json:"outputs,omitempty"`
	Error         *sparkv1.ExecuteSparkError        `json:"error,omitempty"`
}

func (r *runnerTest) Execute(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error) {
	return r.execute(ctx, true, opts...)
}

func (r *runnerTest) ExecuteWithoutStageRetryOverride(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error) {
	return r.execute(ctx, false, opts...)
}

func (r *runnerTest) execute(ctx *sparkv1.JobContext, addStageOverride bool, opts ...sparkv1.Option) (*Outputs, error) {
	// Create the spark chain
	builder := sparkv1.NewBuilder()
	r.spark.BuildChain(builder)
	chain := builder.BuildChain()
	jmd := ctx.Metadata
	outputs := make(sparkv1.BindableMap)

	jmd.VariablesKey = jmd.JobKeyValue

	tmpDir := r.t.TempDir()
	log.Info().Msgf("tmp dir %s %s", ctx.Metadata.JobKeyValue, tmpDir)

	port, err := util.GetFreeTCPPort()
	if err != nil {
		return nil, fmt.Errorf("error getting free tcp port: %w", err)
	}

	s, err := util.RunServerOnPort(port, tmpDir)
	if err != nil {
		return nil, fmt.Errorf("error running nats server: %w", err)
	}
	defer s.Shutdown()
	s.Start()

	nc, js := util.GetNatsClient(port)
	defer nc.Close()

	store, err := js.CreateObjectStore(context.Background(), jetstream.ObjectStoreConfig{
		Bucket: "test",
	})
	if err != nil {
		return nil, fmt.Errorf("error creating object store: %w", err)
	}

	//Initialise spark
	so := new(sparkv1.SparkOpts)
	for _, opt := range opts {
		so = opt(so)
	}
	if err := r.spark.Init(sparkv1.NewInitContext(so)); err != nil {
		return nil, fmt.Errorf("error init spark: %w", err)
	}

	var stageRetryOverride *sparkv1.RetryConfig
	if addStageOverride {
		stageRetryOverride = &sparkv1.RetryConfig{
			Times:             2,
			FirstBackoffWait:  time.Millisecond * 100,
			BackoffMultiplier: 1,
		}
	}

	// Create new workflow
	wf, _ := sparkv1.NewJobWorkflow(
		ctx, uuid.NewString(), chain,
		sparkv1.WithStageTracker(r.InternalStageTracker),
		sparkv1.WithNatsClient(nc),
		sparkv1.WithObjectStore(store),
		sparkv1.WithInputs(ctx.Metadata.Inputs),
		sparkv1.WithConfig(&sparkv1.Config{
			NatsResponseSubject: "agent.v1.job.a.b.test." + ctx.Metadata.JobKeyValue,
		}),
		sparkv1.WithStageRetryOverride(stageRetryOverride),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating new workflow: %w", err)
	}

	// start the request consumer
	requestSubject, err := r.startRequestConsumer(ctx, js, wf)
	if err != nil {
		return nil, fmt.Errorf("error starting request consumer: %w", err)
	}

	// start the response consumer
	responseConsumer, _, err := r.startResponseConsumer(ctx, js)
	if err != nil {
		return nil, fmt.Errorf("error starting response consumer: %w", err)
	}

	b, err := json.Marshal(ctx.Metadata)
	if err != nil {
		return nil, fmt.Errorf("error marshaling metadata: %w", err)
	}

	if _, err := js.Publish(ctx, requestSubject, b); err != nil {
		return nil, fmt.Errorf("error publishing to a stream: %w", err)
	}

	msgs, err := responseConsumer.Fetch(1, jetstream.FetchMaxWait(time.Second*120))
	if err != nil {
		return nil, fmt.Errorf("error fetching messages: %w", err)
	}

	var res *runnerTestOutput
	for msg := range msgs.Messages() {
		if err := json.Unmarshal(msg.Data(), &res); err != nil {
			return nil, fmt.Errorf("error unmarshaling message: %w", err)
		}
	}

	if res == nil {
		return nil, errors.New("timed out")
	}

	if res.Error != nil {
		return nil, res.Error
	}

	output := sparkv1.ExecuteSparkOutput{
		Error:         res.Error,
		JobPid:        jmd.JobPid,
		JobKey:        jmd.JobKeyValue,
		CorrelationId: jmd.CorrelationIdValue,
		TransactionId: jmd.TransactionIdValue,
		Model:         jmd.Model,
	}

	ob, err := store.GetBytes(ctx, res.VariablesKey)
	if err != nil {
		if !errors.Is(err, jetstream.ErrObjectNotFound) {
			return nil, fmt.Errorf("error object not found: %w", err)
		}
	} else {
		if err := json.Unmarshal(ob, &res.Outputs); err != nil {
			return nil, fmt.Errorf("error unmarshaling output: %w", err)
		}
	}

	for k, v := range res.Outputs {
		outputs[k] = v
	}

	output.Outputs = outputs

	return &Outputs{
		ExecuteSparkOutput: output,
	}, nil
}

func (r *runnerTest) startRequestConsumer(ctx *sparkv1.JobContext, js jetstream.JetStream, wf sparkv1.JobWorkflow) (string, error) {
	subject := fmt.Sprintf("agent.v1.job.request.%s", ctx.Metadata.JobKeyValue)

	_, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      "AGENT_JOB_REQ",
		Retention: jetstream.WorkQueuePolicy,
		Subjects: []string{
			subject,
		},
	})
	if err != nil {
		return "", err
	}

	consumer, err := js.CreateOrUpdateConsumer(context.Background(), "AGENT_JOB_REQ", jetstream.ConsumerConfig{
		FilterSubject: subject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       time.Second * 240,
		MaxDeliver:    3,
		MaxAckPending: 3,
	})

	if err != nil {
		return "", err
	}

	go func() {
	loop:
		for {
			select {
			// on consumer stopped
			case <-ctx.Done():
				log.Info().Msgf("stopping consumer")
				break loop
			default:
				batch, err := consumer.Fetch(15, jetstream.FetchMaxWait(time.Second*15))
				if err != nil {
					if errors.Is(err, nats.ErrConnectionClosed) {
						break loop
					}
					continue
				}

				for msg := range batch.Messages() {
					go func(m jetstream.Msg) {
						wf.Run(m)
					}(msg)
				}
			}
		}
	}()

	return subject, nil
}

func (r *runnerTest) startResponseConsumer(ctx *sparkv1.JobContext, js jetstream.JetStream) (jetstream.Consumer, string, error) {
	subject := fmt.Sprintf("agent.v1.job.a.b.%s.%s", "test", ctx.Metadata.JobKeyValue)

	_, err := js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      "AGENT_JOB_RES",
		Retention: jetstream.WorkQueuePolicy,
		Subjects: []string{
			subject,
		},
	})
	if err != nil {
		return nil, "", err
	}

	consumer, err := js.CreateOrUpdateConsumer(ctx, "AGENT_JOB_RES", jetstream.ConsumerConfig{
		FilterSubject: subject,
	})
	if err != nil {
		return nil, "", err
	}

	return consumer, subject, nil
}

func NewTestRunner(t *testing.T, spark sparkv1.Spark, options ...Option) (RunnerTest, error) {
	var to testOpts
	for _, option := range options {
		option(&to)
	}

	st := newStageTracker(t)
	return &runnerTest{spark: spark, testOpts: &to, InternalStageTracker: st, StageTracker: st, t: t}, nil
}

func NewTestJobContext(ctx context.Context, jobKey, correlationId, transactionId string, inputs Inputs) *sparkv1.JobContext {
	ins := make(sparkv1.ExecuteSparkInputs)
	for name, bindable := range inputs {
		ins[name] = sparkv1.NewBindable(sparkv1.Value{Value: MustEncode(bindable.Value), MimeType: string(bindable.MimeType)})
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
