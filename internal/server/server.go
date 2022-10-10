package server

import (
	"context"
	api_ctx "github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

type Server struct {
	worker api.JobWorker
}

func NewServer(worker api.JobWorker) *Server {
	return &Server{worker: worker}
}

func (s Server) Start() error {

	// LOGGER SAMPLE >> add .Fields(fields) with the spark name on it
	logger := log.With().CallerWithSkipFrameCount(3).Stack().Logger()

	svr := grpc.NewServer(grpc.ConnectionTimeout(time.Second * 10)) // TODO env var
	sdk_v1.RegisterAgentServiceServer(svr, s)

	// TODO create an env var around this >> config.Grpc_reflection_enabled?
	reflection.Register(svr)
	// TODO create an env var around this >> config.Grpc_reflection_enabled?

	listener, err := net.Listen("tcp", "localhost:7777") // TODO env var
	if err != nil {
		logger.Error().Err(err).Msg("error setting up the listener")
		return err
	}
	if err = svr.Serve(listener); err != nil {
		logger.Error().Err(err).Msg("error starting the server")
		return err
	}
	return nil
}

func (s Server) ExecuteJob(ctx context.Context, request *sdk_v1.ExecuteJobRequest) (*sdk_v1.Void, error) {
	jobContext := api_ctx.NewJobMetadata(ctx, request.Key, request.CorrelationId, request.TransactionId, nil)
	err := s.worker.Run(jobContext)
	if err != nil {

		return nil, err
	}
	return &sdk_v1.Void{}, nil
}
