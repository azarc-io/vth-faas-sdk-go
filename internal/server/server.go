package server

import (
	"context"
	api_ctx "github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Server struct {
	config    *config.Config
	worker    api.JobWorker
	client    sdk_v1.ManagerServiceClient
	svr       *grpc.Server
	heartBeat *Heartbeat
}

func NewServer(cfg *config.Config, worker api.JobWorker, client sdk_v1.ManagerServiceClient) *Server {
	return &Server{config: cfg, worker: worker, client: client}
}

func (s Server) Start() error {

	// LOGGER SAMPLE >> add .Fields(fields) with the spark name on it
	logger := log.With().CallerWithSkipFrameCount(3).Stack().Logger()

	s.svr = grpc.NewServer(grpc.ConnectionTimeout(time.Second * 10)) // TODO env var
	sdk_v1.RegisterAgentServiceServer(s.svr, s)

	// TODO create an env var around this >> config.Grpc_reflection_enabled?
	reflection.Register(s.svr)
	// TODO create an env var around this >> config.Grpc_reflection_enabled?

	listener, err := net.Listen("tcp", "localhost:7777") // TODO env var
	if err != nil {
		logger.Error().Err(err).Msg("error setting up the listener")
		return err
	}

	s.heartBeat = NewHeartbeat(s.config, s.client)
	s.heartBeat.Start()

	if err = s.svr.Serve(listener); err != nil {
		logger.Error().Err(err).Msg("error starting the server")
		return err
	}
	return nil
}

func (s Server) Stop() {
	s.heartBeat.Stop()
	s.svr.GracefulStop()
}

func (s Server) ExecuteJob(ctx context.Context, request *sdk_v1.ExecuteJobRequest) (*sdk_v1.ExecuteJobResponse, error) {
	jobContext := api_ctx.NewJobMetadata(ctx, request.Key, request.CorrelationId, request.TransactionId, nil)
	go func() { // TODO goroutine pool
		err := s.worker.Run(jobContext)
		if err != nil {
			// we don't care about this error here, it is being sent to the manager service via grpc calls to update the job status
			// TODO fix me
		}
	}()
	return &sdk_v1.ExecuteJobResponse{AgentId: s.config.App.InstanceId}, nil
}
