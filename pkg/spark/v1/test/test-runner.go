package module_test_runner

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/google/uuid"
	gnats "github.com/nats-io/nats-server/v2/server"
	gnatsTest "github.com/nats-io/nats-server/v2/test"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"net"
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

type runnerTestOutput struct {
	Outputs map[string]*sparkv1.BindableValue `json:"outputs,omitempty"`
	Error   *sparkv1.ExecuteSparkError        `json:"error,omitempty"`
}

func (r *runnerTest) Execute(ctx *sparkv1.JobContext, opts ...sparkv1.Option) (*Outputs, error) {
	// Create the spark chain
	builder := sparkv1.NewBuilder()
	r.spark.BuildChain(builder)
	chain := builder.BuildChain()

	tmpDir := r.t.TempDir()
	log.Info().Msgf("tmp dir %s %s", ctx.Metadata.JobKeyValue, tmpDir)

	port, err := GetFreeTCPPort()
	if err != nil {
		return nil, err
	}

	s, err := RunServerOnPort(port, tmpDir)
	if err != nil {
		return nil, err
	}
	defer s.Shutdown()
	s.Start()

	sUrl := fmt.Sprintf("nats://127.0.0.1:%d", port)
	nc, err := nats.Connect(sUrl)
	if err != nil {
		return nil, err
	}
	defer nc.Close()

	if !nc.IsConnected() {
		errorMsg := fmt.Errorf("could not establish connection to nats-server")
		return nil, errorMsg
	}

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}

	_, err = js.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:      "AGENT_JOB_RES",
		Retention: jetstream.WorkQueuePolicy,
		Subjects: []string{
			fmt.Sprintf("agent.v1.job.a.b.%s.%s", "test", ctx.Metadata.JobKeyValue),
		},
	})
	if err != nil {
		return nil, err
	}

	c, err := js.CreateOrUpdateConsumer(ctx, "AGENT_JOB_RES", jetstream.ConsumerConfig{
		FilterSubject: fmt.Sprintf("agent.v1.job.a.b.%s.%s", "test", ctx.Metadata.JobKeyValue),
	})
	if err != nil {
		return nil, err
	}

	//Initialise spark
	so := new(sparkv1.SparkOpts)
	for _, opt := range opts {
		so = opt(so)
	}
	if err := r.spark.Init(sparkv1.NewInitContext(so)); err != nil {
		return nil, err
	}

	// Create new workflow
	wf, _ := sparkv1.NewJobWorkflow(
		ctx, uuid.NewString(), chain,
		sparkv1.WithStageTracker(r.InternalStageTracker),
		sparkv1.WithNatsClient(nc),
		sparkv1.WithConfig(&sparkv1.Config{
			NatsResponseSubject: "agent.v1.job.a.b.test." + ctx.Metadata.JobKeyValue,
		}),
	)

	b, _ := json.Marshal(ctx.Metadata)
	wf.Run(&nats.Msg{
		Data: b,
	})

	msgs, err := c.Fetch(1, jetstream.FetchMaxWait(time.Second*5))
	if err != nil {
		return nil, err
	}

	var res *runnerTestOutput
	for msg := range msgs.Messages() {
		if err := json.Unmarshal(msg.Data(), &res); err != nil {
			return nil, err
		}
	}

	if res == nil {
		return nil, errors.New("timed out")
	}

	if res.Error != nil {
		return nil, res.Error
	}

	output := sparkv1.ExecuteSparkOutput{
		Outputs: make(sparkv1.BindableMap),
		Error:   res.Error,
	}

	for k, v := range res.Outputs {
		output.Outputs[k] = v
	}

	return &Outputs{ExecuteSparkOutput: output}, nil
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

func RunServerOnPort(port int, dir string) (*gnats.Server, error) {
	opts := gnatsTest.DefaultTestOptions
	opts.Port = port
	opts.JetStream = true
	opts.StoreDir = dir
	return RunServerWithOptions(&opts)
}

func RunServerWithOptions(opts *gnats.Options) (*gnats.Server, error) {
	return gnats.NewServer(opts)
}

// GetFreeTCPPort returns free open TCP port
func GetFreeTCPPort() (port int, err error) {
	ln, err := net.Listen("tcp", "[::]:0")
	if err != nil {
		return 0, err
	}
	port = ln.Addr().(*net.TCPAddr).Port
	err = ln.Close()
	return
}
