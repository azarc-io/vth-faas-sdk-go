package spark_v1

import (
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateManagerServiceClient(config *config) (ManagerServiceClient, error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(config.Config.Agent.RetryBackoff)),
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	}
	cc, err := grpc.Dial(config.agentAddress(), opts...)
	if err != nil {
		return nil, err
	}
	return NewManagerServiceClient(cc), nil
}
