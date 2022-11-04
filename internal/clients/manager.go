package clients

import (
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func CreateManagerServiceClient(config *config.Config) (sdk_v1.ManagerServiceClient, error) {
	retryOpts := []grpc_retry.CallOption{
		grpc_retry.WithBackoff(grpc_retry.BackoffExponential(config.ManagerService.RetryBackoff)),
	}
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpc_retry.StreamClientInterceptor(retryOpts...)),
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(retryOpts...)),
	}
	cc, err := grpc.Dial(config.ManagerService.HostPort(), opts...)
	if err != nil {
		return nil, err
	}
	return sdk_v1.NewManagerServiceClient(cc), nil
}
