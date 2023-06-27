package sparkv1

import (
	"context"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
	"time"
)

var (
	ErrInvalidStageResultMimeType = errors.New("stage result expects mime-type of application/json")
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkPlugin struct {
	config *Config
	chain  *SparkChain
	ctx    context.Context
	wrk    worker.Worker
}

/************************************************************************/
// SERVER
/************************************************************************/

func newSparkPlugin(ctx context.Context, cfg *Config, chain *SparkChain) *sparkPlugin {
	return &sparkPlugin{ctx: ctx, config: cfg, chain: chain}
}

func (s *sparkPlugin) start() error {
	tc, err := s.createTemporalClient()
	if err != nil {
		return err
	}

	s.wrk = worker.New(tc, s.config.QueueGroup, worker.Options{})

	var sparkIO SparkDataIO
	if s.config.IOServer != nil {
		log.Info().Msgf("IO Data Provider Enabled")
		sparkIO = &ioDataProvider{ctx: s.ctx, baseUrl: s.config.IOServer.Url, apiKey: s.config.IOServer.ApiKey}
	} else {
		log.Info().Msgf("Temporal Data Provider Enabled")
		sparkIO = &temporalDataProvider{ctx: s.ctx, c: tc}
	}

	wf := NewJobWorkflow(s.ctx, sparkIO, uuid.NewString(), s.chain)
	s.wrk.RegisterActivity(wf.ExecuteStageActivity)
	s.wrk.RegisterActivity(wf.ExecuteCompleteActivity)
	s.wrk.RegisterWorkflowWithOptions(wf.Run, workflow.RegisterOptions{
		Name: s.config.Id,
	})

	if err := s.wrk.Start(); err != nil {
		return err
	}

	// TODO Remove once this is deployed. Testing the workflow runs
	go func() {
		return

		time.Sleep(3 * time.Second)
		o := client.StartWorkflowOptions{
			TaskQueue: s.config.QueueGroup,
		}
		_, err := tc.ExecuteWorkflow(context.Background(), o, s.config.Id, &JobMetadata{
			SparkId: s.config.Id,
			Inputs: map[string]*bindable{
				"name": NewBindableValue("Jono", "application/text"),
			},
		})
		if err != nil {
			log.Error().Err(err).Msgf("workflow run errored")
		}
	}()

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: s.config.Id,
		},
		Plugins: make(map[string]plugin.Plugin),
	})

	return nil
}

func (s *sparkPlugin) stop() {
	s.wrk.Stop()
}

func (s *sparkPlugin) createTemporalClient() (client.Client, error) {
	opts := client.Options{
		HostPort:  s.config.Temporal.Address,
		Namespace: s.config.Temporal.Namespace,
		//TODO Create logger
		Logger:   &TemporalLogger{},
		Identity: s.config.Id,
	}

	return client.Dial(opts)
}

type temporalDataProvider struct {
	ctx context.Context
	c   client.Client
}

func (tdp *temporalDataProvider) GetStageResult(workflowID, runID, stageName, correlationID string) (Bindable, error) {
	res, err := tdp.c.QueryWorkflow(tdp.ctx, workflowID, runID, JobGetStageResultQuery, stageName)
	if err != nil {
		return nil, err
	}

	var val Value
	if err := res.Get(&val); err != nil {
		return nil, err
	}

	if val.MimeType != "application/json" {
		return nil, ErrInvalidStageResultMimeType
	}

	return NewBindable(val), nil
}

func (tdp *temporalDataProvider) PutStageResult(workflowID, runID, stageName, correlationID string, stageValue []byte) (Bindable, error) {
	return &bindable{Value: stageValue, MimeType: string(codec.MimeTypeJson)}, nil
}
