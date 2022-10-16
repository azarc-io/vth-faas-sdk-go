package context_test

import (
	ctx "context"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/mock"
	"github.com/azarc-io/vth-faas-sdk-go/internal/worker"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"github.com/golang/mock/gomock"
	"strings"
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
	stageProgressHandler := inmemory.NewStageProgressHandler(t,
		sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage2", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_stage3", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage1", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_stage2", sdk_v1.StageStatus_StagePending),
		sdk_v1.NewSetStageStatusReq("jobKey", "compensate_compensate_compensate_stage1", sdk_v1.StageStatus_StagePending),
	)
	variableHandler := inmemory.NewVariableHandler(t, nil)
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

func NewInitExecutor() *InitExecutor {
	return &InitExecutor{}
}

func (i *InitExecutor) Initialize() error {
	if i.m == nil {
		i.m = map[string]uint{}
	}
	i.m["initialize"] += 1
	return nil
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
		stageProgressHandler := inmemory.NewStageProgressHandler(t, sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending))
		variableHandler := inmemory.NewVariableHandler(t, nil)
		return context.NewJobContext(jobMetadata, stageProgressHandler, variableHandler)
	}

	initExecutor := &InitExecutor{}

	defer func() {
		if r := recover(); r != nil {
			err := initExecutor.Initialize() // now the map is initialized
			if err != nil {
				t.Error("error initializing the map: ", err)
			}
			initExecutor.Execute(newCxt())
			if initExecutor.m["stage1"] != 1 {
				t.Errorf("counter expected to be 1 and we got: %d", initExecutor.m["stage1"])
			}
			initExecutor.Execute(newCxt())
			if initExecutor.m["stage1"] != 2 {
				t.Errorf("counter expected to be 2 and we got: %d", initExecutor.m["stage1"])
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
	stageProgressHandler.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage2", sdk_v1.StageStatus_StageStarted))
	stageProgressHandler.EXPECT().Set(sdk_v1.NewSetStageStatusReq("jobKey", "stage2", sdk_v1.StageStatus_StageCompleted))

	var executed bool

	jobContext.Stage("stage2", func(context api.StageContext) (any, api.StageError) {
		executed = true
		return nil, nil
	}, context.WithStageStatus("stage1", sdk_v1.StageStatus_StagePending))

	if !executed {
		t.Error("conditional stage execution should be triggered")
	}
}

func TestShouldSkipConditionalStageExecution(t *testing.T) {
	jobContext, mockCtrlFinish, stageProgressHandler, _ := newMockJobContext(t)
	defer mockCtrlFinish()

	stageProgressHandler.EXPECT().Get("jobKey", "stage1").Return(sdk_v1.Ptr(sdk_v1.StageStatus_StagePending), nil)

	var executed bool

	jobContext.Stage("stage2", func(context api.StageContext) (any, api.StageError) {
		executed = true
		return nil, nil
	}, context.WithStageStatus("stage1", sdk_v1.StageStatus_StageCanceled))

	if executed {
		t.Error("conditional stage execution should be skipped")
	}

	if jobContext.Err().Error() != "conditional stage execution skipped this stage" {
		t.Errorf("error message expected: '%s'; got: '%s'", "conditional stage execution skipped this stage", jobContext.Err().Error())
	}

}

type T struct {
	A string
	B struct {
		C int
		D struct {
			E []string
		}
	}
}

// TODO move to the right place
func TestJsonVariableTest(t *testing.T) {
	sample := T{
		A: "a",
		B: struct {
			C int
			D struct {
				E []string
			}
		}{
			C: 1,
			D: struct {
				E []string
			}{E: []string{"1", "2", "3"}},
		},
	}

	v, err := sdk_v1.NewVariable("test_var", "application/json", sample)
	if err != nil {
		t.Error("error serializing sdk variable: ", err)
	}
	var fromValue T
	err = v.Bind(&fromValue)
	if err != nil {
		t.Error("error deserializing sdk variable: ", err)
	}
	if fromValue.A != sample.A {
		t.Errorf("error serializing: expect: %v' got: %v", fromValue.A, sample.A)
	}
	if fromValue.B.C != sample.B.C {
		t.Errorf("error serializing: expect: %v' got: %v", fromValue.B.C, sample.B.C)
	}
	if strings.Join(fromValue.B.D.E, "") != strings.Join(sample.B.D.E, "") {
		t.Errorf("error serializing: expect: %v' got: %v", fromValue.B.D.E, sample.B.D.E)
	}
}

// TODO move to the right place
func TestRawVariableValue(t *testing.T) {
	v, err := sdk_v1.NewVariable("test_var", "application/json", "test")
	if err != nil {
		t.Error("error creating a variable: ", err)
	}
	if b, er := v.Raw(); er != nil || string(b) != `"test"` {
		t.Errorf("error getting the raw value from a variable. expected: '%s' got: '%s', error: %s", "test", string(b), err)
	}

	v, err = sdk_v1.NewVariable("test_var", "application/json", []byte{1, 2, 3, 3, 2, 5, 5, 8, 4, 5, 7, 6, 4, 0, 2, 8, 9, 5, 7, 4, 9, 8, 5, 7, 4, 2, 3})
	if err != nil {
		t.Error("error creating a variable: ", err)
	}
	if b, er := v.Raw(); er != nil || string(b) != `"AQIDAwIFBQgEBQcGBAACCAkFBwQJCAUHBAID"` {
		t.Errorf("error getting the raw value from a variable. expected: '%s' got: '%s', error: %s", "test", string(b), err)
	}
}

func TestJobWorker(t *testing.T) {
	stageProgressHandler := inmemory.NewStageProgressHandler(t, sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending))
	variablesHandler := inmemory.NewVariableHandler(t, nil)
	job := NewInitExecutor()
	cfg, err := config.NewMock(map[string]string{"APP_ENVIRONMENT": "test", "AGENT_SERVER_PORT": "0", "MANAGER_SERVER_PORT": "0"})
	if err != nil {
		t.Error(err)
	}
	jobWorker, err := worker.NewJobWorker(cfg, job, worker.WithStageProgressHandler(stageProgressHandler), worker.WithVariableHandler(variablesHandler))
	if err != nil {
		t.Error("error instantiating the job worker: ", err)
	}
	for range []int{5, 5, 5, 5, 5} {
		err = jobWorker.Run(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))
		if err != nil {
			t.Error("error running the job worker: ", err)
		}
	}
	if job.m["initialize"] != 1 {
		t.Errorf("initialize should run only once; got: %d", job.m["initialize"])
	}
	if job.m["stage1"] != 5 {
		t.Errorf("stage1 should run 5 times; got: %d", job.m["stage1"])
	}
}

type ResumeTestJob struct {
	stageExecutionCounter map[string]uint
}

func (r ResumeTestJob) inc(stageName string) {
	r.stageExecutionCounter[stageName] += 1
}

func (r ResumeTestJob) Initialize() error {
	if r.stageExecutionCounter == nil {
		r.stageExecutionCounter = map[string]uint{}
	}
	r.inc("initialize")
	return nil
}

func (r ResumeTestJob) Execute(jobContext api.JobContext) {
	jobContext.Stage("stage-skip-1", func(stageContext api.StageContext) (any, api.StageError) {
		r.inc("stage-skip-1")
		return nil, nil
	}).Stage("stage-skip-2", func(stageContext api.StageContext) (any, api.StageError) {
		r.inc("stage-skip-1")
		return nil, nil
	}).Stage("stage-execute-1", func(stageContext api.StageContext) (any, api.StageError) {
		r.inc("stage-execute-1")
		return nil, nil
	}).Stage("stage-execute-2", func(stageContext api.StageContext) (any, api.StageError) {
		r.inc("stage-execute-2")
		return nil, nil
	}).Complete(func(completionContext api.CompletionContext) api.StageError {
		r.inc("stage-complete")
		return nil
	})
}

func TestResumeRetryLastActiveStageCompleted(t *testing.T) {
	stageProgressHandler := inmemory.NewStageProgressHandler(t, sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending))
	variablesHandler := inmemory.NewVariableHandler(t, nil)

	jobMetadata := context.NewJobMetadataFromGrpcRequest(nil, &sdk_v1.ExecuteJobRequest{
		Key:           "jobKey",
		TransactionId: "transactionId",
		CorrelationId: "correlationId",
		LastActiveStage: &sdk_v1.LastActiveStage{
			Name:   "stage-skip-2",
			Status: sdk_v1.StageStatus_StageCompleted,
		},
	})
	jobContext := context.NewJobContext(jobMetadata, stageProgressHandler, variablesHandler)
	resumeTestJob := ResumeTestJob{map[string]uint{}}
	_ = resumeTestJob.Initialize()
	resumeTestJob.Execute(jobContext)

	if resumeTestJob.stageExecutionCounter["initialize"] != 1 {
		t.Error("initialize executed more than 1 time")
	}

	if resumeTestJob.stageExecutionCounter["stage-skip-1"] != 0 {
		t.Error("stage-skip-1 executed")
	}

	if resumeTestJob.stageExecutionCounter["stage-skip-2"] != 0 {
		t.Error("stage-skip-1 executed")
	}

	if resumeTestJob.stageExecutionCounter["stage-execute-1"] != 1 {
		t.Errorf("stage-execute-1 executed: %d times; expected: 1 time", resumeTestJob.stageExecutionCounter["stage-execute-1"])
	}

	if resumeTestJob.stageExecutionCounter["stage-execute-2"] != 1 {
		t.Errorf("stage-execute-2 executed: %d times; expected: 1 time", resumeTestJob.stageExecutionCounter["stage-execute-2"])
	}
}

func TestResumeRetry(t *testing.T) {
	stageProgressHandler := inmemory.NewStageProgressHandler(t, sdk_v1.NewSetStageStatusReq("jobKey", "stage1", sdk_v1.StageStatus_StagePending))
	variablesHandler := inmemory.NewVariableHandler(t, nil)

	jobMetadata := context.NewJobMetadataFromGrpcRequest(nil, &sdk_v1.ExecuteJobRequest{
		Key:           "jobKey",
		TransactionId: "transactionId",
		CorrelationId: "correlationId",
		LastActiveStage: &sdk_v1.LastActiveStage{
			Name:   "stage-execute-1",
			Status: sdk_v1.StageStatus_StageFailed,
		},
	})
	jobContext := context.NewJobContext(jobMetadata, stageProgressHandler, variablesHandler)
	resumeTestJob := ResumeTestJob{map[string]uint{}}
	_ = resumeTestJob.Initialize()
	resumeTestJob.Execute(jobContext)

	if resumeTestJob.stageExecutionCounter["initialize"] != 1 {
		t.Error("initialize executed more than 1 time")
	}

	if resumeTestJob.stageExecutionCounter["stage-skip-1"] != 0 {
		t.Error("stage-skip-1 executed")
	}

	if resumeTestJob.stageExecutionCounter["stage-skip-2"] != 0 {
		t.Error("stage-skip-1 executed")
	}

	if resumeTestJob.stageExecutionCounter["stage-execute-1"] != 1 {
		t.Errorf("stage-execute-1 executed: %d times; expected: 1 time", resumeTestJob.stageExecutionCounter["stage-execute-1"])
	}

	if resumeTestJob.stageExecutionCounter["stage-execute-2"] != 1 {
		t.Errorf("stage-execute-2 executed: %d times; expected: 1 time", resumeTestJob.stageExecutionCounter["stage-execute-2"])
	}
}
