package spark_v1

import (
	"context"
	"net"
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

type Server struct {
	config *Config
	worker Worker
	svr    *grpc.Server
}

/************************************************************************/
// SERVER
/************************************************************************/

func newServer(cfg *Config, worker Worker) *Server {
	return &Server{config: cfg, worker: worker}
}

func (s *Server) start() error {
	// LOGGER SAMPLE >> add .Fields(fields) with the spark name on it
	log := NewLogger()

	// nosemgrep
	s.svr = grpc.NewServer(grpc.ConnectionTimeout(connectionTimeout))
	RegisterAgentServiceServer(s.svr, s)

	reflection.Register(s.svr)

	listener, err := net.Listen("tcp", s.config.ServerAddress())
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

func (s *Server) stop() {
	if s.svr != nil {
		s.svr.GracefulStop()
	}
}

/************************************************************************/
// RPC IMPLEMENTATIONS
/************************************************************************/

func (s *Server) ExecuteJob(ctx context.Context, request *ExecuteJobRequest) (*ExecuteJobResponse, error) {
	jobContext := NewSparkMetadata(ctx, request.Key, request.CorrelationId, request.TransactionId, nil)
	go func() { // TODO goroutine pool
		_ = s.worker.Execute(jobContext)
	}()
	return &ExecuteJobResponse{AgentId: s.config.Config.App.InstanceID}, nil
}
