package sparkv1

import (
	"context"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"time"
)

var (
	ErrInvalidStageResultMimeType = errors.New("stage result expects mime-type of application/json")
	ErrTemporalIoNotSupported     = errors.New("temporal IO provider does not support input/output referencing")
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkPlugin struct {
	config *Config
	chain  *SparkChain
	ctx    context.Context
	nc     *nats.Conn
}

/************************************************************************/
// SERVER
/************************************************************************/

func newSparkPlugin(ctx context.Context, cfg *Config, chain *SparkChain) *sparkPlugin {
	return &sparkPlugin{ctx: ctx, config: cfg, chain: chain}
}

func (s *sparkPlugin) start() error {
	nc, err := s.createNatsClient()
	if err != nil {
		return err
	}
	s.nc = nc

	js, err := jetstream.New(nc)
	if err != nil {
		return err
	}
	store, err := js.ObjectStore(s.ctx, s.config.NatsBucket)
	if err != nil {
		return err
	}

	wf, err := NewJobWorkflow(s.ctx, uuid.NewString(), s.chain,
		WithConfig(s.config), WithNatsClient(nc), WithObjectStore(store))
	if err != nil {
		return err
	}

	stream, err := js.Stream(s.ctx, s.config.NatsRequestStreamName)
	if err != nil {
		return err
	}

	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name:          s.config.Id,
		FilterSubject: s.config.NatsRequestSubject,
		AckPolicy:     jetstream.AckExplicitPolicy,
		AckWait:       s.config.Timeout,
		MaxDeliver:    1,
		MaxAckPending: 1,
	})
	if err != nil {
		log.Error().Err(err).Msgf("could not create consumer for subject %s", s.config.NatsRequestSubject)
		return err
	}

	go func() {
	loop:
		for {
			select {
			// on consumer stopped
			case <-s.ctx.Done():
				log.Info().Msgf("stopping consumer")
				break loop
			default:
				batch, err := consumer.Fetch(15, jetstream.FetchMaxWait(time.Second*15))
				if err != nil {
					log.Error().Err(err).Msgf("failed to fetch job request messages, will retry shortly")
					continue
				}

				for msg := range batch.Messages() {
					go func(m jetstream.Msg) {
						m.Ack()
						wf.Run(m)
					}(msg)
				}
			}
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
	if s.nc != nil {
		_ = s.nc.Drain()
	}
}

func (s *sparkPlugin) createNatsClient() (*nats.Conn, error) {
	return nats.Connect(s.config.Nats.Address)
}

type temporalDataProvider struct {
	ctx context.Context
	c   client.Client
}

func (tdp *temporalDataProvider) NewInput(_ string, value *BindableValue) Bindable {
	return value
}

func (tdp *temporalDataProvider) NewOutput(_ string, value *BindableValue) (Bindable, error) {
	return value, nil
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
	return &BindableValue{Value: stageValue, MimeType: string(codec.MimeTypeJson)}, nil
}
