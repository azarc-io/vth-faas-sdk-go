package spark_v1

import (
	"context"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"github.com/rs/zerolog/log"
	"net"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

/************************************************************************/
// CONSTANTS
/************************************************************************/

var connectionTimeout = time.Second * 10

/************************************************************************/
// TYPES
/************************************************************************/

type server struct {
	config *config
	worker Worker
	svr    *grpc.Server
	ctx    context.Context
	client StageProgressHandler
}

/************************************************************************/
// SERVER
/************************************************************************/

func newServer(cfg *config, worker Worker, ctx context.Context, client StageProgressHandler) *server {
	return &server{config: cfg, worker: worker, ctx: ctx, client: client}
}

func (s *server) start() error {
	// LOGGER SAMPLE >> add .Fields(fields) with the spark name on it
	log := NewLogger()

	// nosemgrep
	s.svr = grpc.NewServer(grpc.ConnectionTimeout(connectionTimeout))
	sparkv1.RegisterAgentServiceServer(s.svr, s)

	reflection.Register(s.svr)

	listener, err := net.Listen("tcp", s.config.serverAddress())
	if err != nil {
		log.Error(err, "error setting up the listener")
		return err
	}

	// nosemgrep
	if err = s.svr.Serve(listener); err != nil {
		log.Error(err, "error starting the server")
		return err
	}
	return nil
}

func (s *server) stop() {
	if s.svr != nil {
		s.svr.GracefulStop()
	}
}

/************************************************************************/
// RPC IMPLEMENTATIONS
/************************************************************************/

func (s *server) ExecuteJob(ctx context.Context, request *sparkv1.ExecuteJobRequest) (*sparkv1.Void, error) {
	log.Info().Msgf("execute with metadata: %v", request.RequestMetadata.Metadata)
	jobContext := NewSparkMetadata(s.ctx, request.Key, request.CorrelationId, request.TransactionId, nil,
		WithPerRequestMetadata(request.RequestMetadata.Metadata))
	accepted := make(chan error)
	go func() { // TODO goroutine pool
		log.Info().Msgf("executing job: %s", request.Key)
		err := s.client.JobStarting(&sparkv1.JobStartingRequest{
			Key:       request.Key,
			ServiceId: os.Getenv("SERVICE_ID"),
		})

		if err != nil {
			log.Error().Err(err).Msgf("could not start job")
			accepted <- err
			close(accepted)
			return
		}

		close(accepted)

		_ = s.worker.Execute(jobContext)
	}()
	return &sparkv1.Void{}, <-accepted
}
