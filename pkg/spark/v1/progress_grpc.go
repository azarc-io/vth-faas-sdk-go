package spark_v1

import (
	"context"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"time"
)

type stageProgressHandler struct {
	client sparkv1.ManagerServiceClient
}

func (g *stageProgressHandler) FinishJob(ctx SparkContext, req *sparkv1.FinishJobRequest) error {
	nctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	_, err := g.client.FinishJob(nctx, req,
		grpc_retry.WithMax(3), grpc_retry.WithPerRetryTimeout(time.Second*5))
	return err
}

func newGrpcStageProgressHandler(client sparkv1.ManagerServiceClient) StageProgressHandler {
	return &stageProgressHandler{client: client}
}

func (g *stageProgressHandler) Get(ctx SparkContext, name string) (*sparkv1.StageStatus, error) {
	nctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	resp, err := g.client.GetStageStatus(nctx, newGetStageStatusReq(ctx, name),
		grpc_retry.WithMax(3), grpc_retry.WithPerRetryTimeout(time.Second*5))
	return &resp.Status, err
}

func (g *stageProgressHandler) Set(ctx SparkContext, stageStatus *sparkv1.SetStageStatusRequest) error {
	nctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	_, err := g.client.SetStageStatus(nctx, stageStatus,
		grpc_retry.WithMax(3), grpc_retry.WithPerRetryTimeout(time.Second*5))
	return err
}

func (g *stageProgressHandler) GetResult(ctx SparkContext, name string) Bindable {
	nctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	result, err := g.client.GetStageResult(nctx, newStageResultReq(ctx, name),
		grpc_retry.WithMax(3), grpc_retry.WithPerRetryTimeout(time.Second*5))
	if err != nil {
		return newResult(err, nil)
	}
	return newResult(nil, result)
}

func (g *stageProgressHandler) SetResult(ctx SparkContext, result *sparkv1.SetStageResultRequest) error {
	nctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	_, err := g.client.SetStageResult(nctx, result,
		grpc_retry.WithMax(3), grpc_retry.WithPerRetryTimeout(time.Second*5))
	return err
}

func (g *stageProgressHandler) JobStarting(result *sparkv1.JobStartingRequest) error {
	ctx, _ := context.WithTimeout(context.TODO(), time.Second*30)
	_, err := g.client.JobStarting(ctx, result,
		grpc_retry.WithMax(3), grpc_retry.WithPerRetryTimeout(time.Second*5))
	return err
}
