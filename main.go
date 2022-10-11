package main

import (
	ctx "context"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/lithammer/shortuuid/v4"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"time"
)

func main() {
	jobMetadata := context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", "payload")
	stageProgressHandler := inmemory.NewMockStageProgressHandler(nil, sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending))
	variableHandler := inmemory.NewMockVariableHandler(nil)
	jobContext := context.NewJobContext(jobMetadata, stageProgressHandler, variableHandler)
	jobContext.Stage("stage1", func(stageContext api.StageContext) (any, api.StageError) {
		println("stage1 exec")
		return nil, nil
	})

	svr := grpc.NewServer(grpc.ConnectionTimeout(time.Second * 10))

	sdk_v1.RegisterAgentServiceServer(svr, AgentService{})
	reflection.Register(svr)

	listener, err := net.Listen("tcp", "localhost:7777")
	if err != nil {
		println("err: ", err.Error())
		return
	}
	if err := svr.Serve(listener); err != nil {
		println("err: ", err.Error())
	}
}

type AgentService struct{}

func (a AgentService) ExecuteJob(c ctx.Context, req *sdk_v1.ExecuteJobRequest) (*sdk_v1.ExecuteJobResponse, error) {
	println("ExecuteJob req!, ", fmt.Sprintf("job: %s, tr: %s, co: %s", req.Key, req.TransactionId, req.CorrelationId))
	return &sdk_v1.ExecuteJobResponse{AgentId: shortuuid.New()}, nil
}
