package sparkv1

import (
	"context"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"github.com/google/uuid"
	"github.com/hashicorp/go-plugin"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/rs/zerolog/log"
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

	err = s.createEventConsumer(js, wf)
	if err != nil {
		return err
	}

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

func (s *sparkPlugin) createEventConsumer(js jetstream.JetStream, wf JobWorkflow) error {
	stream, err := js.CreateOrUpdateStream(s.ctx, jetstream.StreamConfig{
		Name: s.config.NatsRequestStreamName,
	})
	if err != nil {
		return err
	}

	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name:              s.config.Id,
		FilterSubject:     s.config.NatsRequestSubject,
		AckPolicy:         jetstream.AckExplicitPolicy,
		AckWait:           s.config.Timeout,
		MaxDeliver:        1,
		MaxAckPending:     1,
		InactiveThreshold: maxInactiveConsumerDuration,
	})
	if err != nil {
		log.Error().Err(err).Msgf("could not create consumer for subject %s", s.config.NatsRequestSubject)
		return err
	}

	go func() {
		var lastConsumedTime = time.Now()
	loop:
		for {
			select {
			// on consumer stopped
			case <-s.ctx.Done():
				log.Info().Msgf("stopping consumer")
				return
			default:
				batch, err := consumer.Fetch(15, jetstream.FetchMaxWait(maxConsumerFetchWait))
				if err != nil {
					log.Error().Err(err).Msgf("failed to fetch job request messages, will retry shortly")
					continue
				}

				var received bool
				for msg := range batch.Messages() {
					received = true
					go func(m jetstream.Msg) {
						m.Ack()
						wf.Run(m)
					}(msg)
				}

				if received {
					lastConsumedTime = time.Now()
				} else if time.Since(lastConsumedTime) > maxInactiveResetConsumerDuration {
					err := backoff.Retry(func() error {
						return s.createEventConsumer(js, wf)
					}, backoff.WithMaxRetries(backoff.NewExponentialBackOff(), 3))

					if err != nil {
						log.Error().Err(err).Msgf(
							"failed to re-subscribe audit event consumer after it became idle for too long")
					}

					break loop
				}
			}
		}
	}()

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
