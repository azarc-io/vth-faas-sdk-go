package sdk_v1

import (
	"context"
	"net"
	"time"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	config    *config.Config
	worker    Worker
	client    ManagerServiceClient
	svr       *grpc.Server
	heartBeat *Heartbeat
}

func NewServer(cfg *config.Config, worker Worker, client ManagerServiceClient) *Server {
	return &Server{config: cfg, worker: worker, client: client}
}

var connectionTimeout = time.Second * 10

func (s Server) Start() error {
	// LOGGER SAMPLE >> add .Fields(fields) with the spark name on it
	log := NewLogger()

	// nosemgrep
	s.svr = grpc.NewServer(grpc.ConnectionTimeout(connectionTimeout)) // TODO env var
	RegisterAgentServiceServer(s.svr, s)

	// TODO create an env var around this >> config.Grpc_reflection_enabled?
	reflection.Register(s.svr)
	// TODO create an env var around this >> config.Grpc_reflection_enabled?

	listener, err := net.Listen("tcp", "localhost:7777") // TODO env var
	if err != nil {
		log.Error(err, "error setting up the listener")
		return err
	}

	s.heartBeat = NewHeartbeat(s.config, s.client)
	s.heartBeat.Start()

	// nosemgrep
	if err = s.svr.Serve(listener); err != nil {
		log.Error(err, "error starting the server")
		return err
	}
	return nil
}

func (s Server) Stop() {
	s.heartBeat.Stop()
	s.svr.GracefulStop()
}

func (s Server) ExecuteJob(ctx context.Context, request *ExecuteJobRequest) (*ExecuteJobResponse, error) {
	jobContext := NewSparkMetadata(ctx, request.Key, request.CorrelationId, request.TransactionId, nil)
	go func() { // TODO goroutine pool
		_ = s.worker.Execute(jobContext)
	}()
	return &ExecuteJobResponse{AgentId: s.config.App.InstanceID}, nil
}
