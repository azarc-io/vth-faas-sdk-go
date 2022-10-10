package context_test

import (
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/mock"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"github.com/golang/mock/gomock"
	"testing"
)

func newMockJobContext(t *testing.T) (api.JobContext, func(), *mock.MockStageProgressHandler, *mock.MockVariableHandler) {
	mockCtrl := gomock.NewController(t)
	stageProgressHandlerMock := mock.NewMockStageProgressHandler(mockCtrl)
	variableHandlerMock := mock.NewMockVariableHandler(mockCtrl)
	jobMetadata := context.NewJobMetadata(nil, "jobKey", "correlationId", "transactionId", "payload")
	jobContext := context.NewJobContext(jobMetadata, stageProgressHandlerMock, variableHandlerMock)
	return jobContext, mockCtrl.Finish, stageProgressHandlerMock, variableHandlerMock
}

//func newInMemoryJobContext(t *testing.T) (context.Job, api.StageProgressHandler, api.VariableHandler) {
//
//}

func TestRunSingleStageWithSuccess(t *testing.T) {
	jobContext, mockCtrlFinish, stageProgressHandler, _ := newMockJobContext(t)
	defer mockCtrlFinish()

	stageProgressHandler.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StageStarted))
	stageProgressHandler.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StageCompleted))

	jobContext.Stage("stage1", func(context api.StageContext) (any, api.StageError) {
		return nil, nil
	})
}

func TestStageBuilder(t *testing.T) {
	jobMetadata := context.NewJobMetadata(nil, "jobKey", "correlationId", "transactionId", "payload")
	stageProgressHandler := inmemory.NewMockStageProgressHandler(t,
		sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage2", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage3", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage2", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_compensate_stage1", sdk_v1.StageStatus_StagePending),
	)
	variableHandler := inmemory.NewMockVariableHandler(t)
	jobContext := context.NewJobContext(jobMetadata, stageProgressHandler, variableHandler)
	jobContext.Stage("stage1", func(context api.StageContext) (any, api.StageError) {
		println("---- stage1")
		return nil, sdk_errors.NewStageError(errors.New("err"))
	}).Compensate(func(context api.CompensationContext) api.StageError {
		context.Stage("compensate_stage1", func(context api.StageContext) (any, api.StageError) {
			println("---- compensation stage1")
			return nil, nil
		}).Stage("compensate_stage2", func(context api.StageContext) (any, api.StageError) {
			println("---- compensation stage2")
			return nil, nil
		}).Stage("compensate_stage3", func(context api.StageContext) (any, api.StageError) {
			println("---- compensation stage3")
			return nil, sdk_errors.NewStageError(errors.New("as"))
		}).Complete(func(context api.CompletionContext) api.StageError {
			println("---- [should not display this] compensation complete")
			return nil
		}).Compensate(func(context api.CompensationContext) api.StageError {
			println("---- compensation compensate")
			context.Stage("compensate_compensate_stage1", func(context api.StageContext) (any, api.StageError) {
				println("---- compensate compensate stage1")
				return nil, nil
			}).Stage("compensate_compensate_stage2", func(context api.StageContext) (any, api.StageError) {
				println("---- compensate compensate stage2")
				return nil, sdk_errors.NewStageError(errors.New("as"))
			}).Compensate(func(context api.CompensationContext) api.StageError {
				println("---- compensation compensate compensate")
				context.Stage("compensate_compensate_compensate_stage1", func(context api.StageContext) (any, api.StageError) {
					println("---- compensate compensate compensate stage1")
					return nil, nil
				})
				return nil
			})
			return nil
		})
		return nil
	})
}

func TestStageBuilderMock(t *testing.T) {

	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()

	stageProgressHandlerMock := mock.NewMockStageProgressHandler(mockCtrl)
	variableHandlerMock := mock.NewMockVariableHandler(mockCtrl)

	jobMetadata := context.NewJobMetadata(nil, "jobKey", "correlationId", "transactionId", "payload")

	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StageFailed, &sdk_v1.Error{Error: "err", ErrorCode: 0, ErrorType: sdk_v1.ErrorType_Failed})).Return(nil)
	stageProgressHandlerMock.EXPECT().SetJobStatus(sdk_v1.NewSetJobStatusReq("jobKey", sdk_v1.JobStatus_JobCompensationStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage1", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage1", sdk_v1.StageStatus_StageCompleted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage2", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage2", sdk_v1.StageStatus_StageCompleted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage3", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage3", sdk_v1.StageStatus_StageFailed, &sdk_v1.Error{Error: "as", ErrorCode: 0, ErrorType: sdk_v1.ErrorType_Failed})).Return(nil)
	stageProgressHandlerMock.EXPECT().SetJobStatus(sdk_v1.NewSetJobStatusReq("jobKey", sdk_v1.JobStatus_JobCompensationStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage1", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage1", sdk_v1.StageStatus_StageCompleted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage2", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage2", sdk_v1.StageStatus_StageFailed, &sdk_v1.Error{Error: "as", ErrorCode: 0, ErrorType: sdk_v1.ErrorType_Failed})).Return(nil)
	stageProgressHandlerMock.EXPECT().SetJobStatus(sdk_v1.NewSetJobStatusReq("jobKey", sdk_v1.JobStatus_JobCompensationStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_compensate_stage1", sdk_v1.StageStatus_StageStarted)).Return(nil)
	stageProgressHandlerMock.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_compensate_stage1", sdk_v1.StageStatus_StageCompleted)).Return(nil)
	stageProgressHandlerMock.EXPECT().SetJobStatus(sdk_v1.NewSetJobStatusReq("jobKey", sdk_v1.JobStatus_JobCompensationDone)).Return(nil)
	stageProgressHandlerMock.EXPECT().SetJobStatus(sdk_v1.NewSetJobStatusReq("jobKey", sdk_v1.JobStatus_JobCompensationDone)).Return(nil)
	stageProgressHandlerMock.EXPECT().SetJobStatus(sdk_v1.NewSetJobStatusReq("jobKey", sdk_v1.JobStatus_JobCompensationDone)).Return(nil)

	jobContext := context.NewJobContext(jobMetadata, stageProgressHandlerMock, variableHandlerMock)

	jobContext.Stage("stage1", func(context api.StageContext) (any, api.StageError) {
		println("---- stage1")
		return nil, sdk_errors.NewStageError(errors.New("err"))
	}).Compensate(func(context api.CompensationContext) api.StageError {
		context.Stage("compensate_stage1", func(context api.StageContext) (any, api.StageError) {
			println("---- compensation stage1")
			return nil, nil
		}).Stage("compensate_stage2", func(context api.StageContext) (any, api.StageError) {
			println("---- compensation stage2")
			return nil, nil
		}).Stage("compensate_stage3", func(context api.StageContext) (any, api.StageError) {
			println("---- compensation stage3")
			return nil, sdk_errors.NewStageError(errors.New("as"))
		}).Complete(func(context api.CompletionContext) api.StageError {
			println("---- [should not display this] compensation complete")
			return nil
		}).Compensate(func(context api.CompensationContext) api.StageError {
			println("---- compensation compensate")
			context.Stage("compensate_compensate_stage1", func(context api.StageContext) (any, api.StageError) {
				println("---- compensate compensate stage1")
				return nil, nil
			}).Stage("compensate_compensate_stage2", func(context api.StageContext) (any, api.StageError) {
				println("---- compensate compensate stage2")
				return nil, sdk_errors.NewStageError(errors.New("as"))
			}).Compensate(func(context api.CompensationContext) api.StageError {
				println("---- compensation compensate compensate")
				context.Stage("compensate_compensate_compensate_stage1", func(context api.StageContext) (any, api.StageError) {
					println("---- compensate compensate compensate stage1")
					return nil, nil
				})
				return nil
			})
			return nil
		})
		return nil
	})
}

type InitExecutor struct {
	m map[string]uint
}

var _ api.Job = &InitExecutor{}

func NewInitExecutor() InitExecutor {
	return InitExecutor{map[string]uint{}}
}

func (i *InitExecutor) Execute(ctx api.JobContext) {
	ctx.Stage("stage1", func(stageContext api.StageContext) (any, api.StageError) {
		i.m["stage1"] += 1
		return nil, nil
	})
	return
}

func TestInitialization(t *testing.T) {
	newCxt := func() api.JobContext {
		jobMetadata := context.NewJobMetadata(nil, "jobKey", "correlationId", "transactionId", "payload")
		stageProgressHandler := inmemory.NewMockStageProgressHandler(t, sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending))
		variableHandler := inmemory.NewMockVariableHandler(t)
		return context.NewJobContext(jobMetadata, stageProgressHandler, variableHandler)
	}

	initExecutor := InitExecutor{}

	defer func() {
		if r := recover(); r != nil {
			initExecutor = NewInitExecutor() // now the map is initialized
			initExecutor.Execute(newCxt())
			if initExecutor.m["stage1"] != 1 {
				t.Errorf("counter expected to be 1 and we got: %d", initExecutor.m["stage1"])
			}
			initExecutor.Execute(newCxt())
			if initExecutor.m["stage1"] != 2 {
				t.Errorf("counter expected to be 1 and we got: %d", initExecutor.m["stage1"])
			}
		}
	}()

	ctx := newCxt()
	initExecutor.Execute(ctx) // panic: assignment to entry in nil map
	t.Fatal("should have panicked because the map in the executor isn't initialized")
}

func TestShouldTriggerConditionalStageExecution(t *testing.T) {
	jobContext, mockCtrlFinish, stageProgressHandler, _ := newMockJobContext(t)
	defer mockCtrlFinish()

	stageProgressHandler.EXPECT().Get("jobKey", "stage1").Return(sdk_v1.Ptr(sdk_v1.StageStatus_StagePending), nil)
	stageProgressHandler.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StageStarted))
	stageProgressHandler.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StageCompleted))

	var executed bool

	jobContext.Stage("stage1", func(context api.StageContext) (any, api.StageError) {
		executed = true
		return nil, nil
	}, context.WithStageStatus(sdk_v1.StageStatus_StagePending))

	if !executed {
		t.Error("conditional stage execution should be triggered")
	}
}

func TestShouldSkipConditionalStageExecution(t *testing.T) {
	jobContext, mockCtrlFinish, stageProgressHandler, _ := newMockJobContext(t)
	defer mockCtrlFinish()

	stageProgressHandler.EXPECT().Get("jobKey", "stage1").Return(sdk_v1.Ptr(sdk_v1.StageStatus_StagePending), nil)

	var executed bool

	jobContext.Stage("stage1", func(context api.StageContext) (any, api.StageError) {
		executed = true
		return nil, nil
	}, context.WithStageStatus(sdk_v1.StageStatus_StageCanceled))

	if executed {
		t.Error("conditional stage execution should be skipped")
	}

	if jobContext.Err().Error() != "conditional stage execution skipped this stage" {
		t.Errorf("error message expected: '%s'; got: '%s'", "conditional stage execution skipped this stage", jobContext.Err().Error())
	}

}
