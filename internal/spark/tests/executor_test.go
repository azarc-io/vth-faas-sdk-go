package tests

import (
	ctx "context"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSparkExecutor(t *testing.T) {
	tests := []struct {
		name            string
		chainFn         func() (*spark.Chain, *stageBehaviour)
		assertFn        func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler)
		errorType       *sdk_v1.ErrorType
		lastActiveStage *sdk_v1.LastActiveStage
		prepare         func(sph *inmemory.StageProgressHandler)
	}{
		{
			name: "should execute single stage spark",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
			},
		},
		{
			name: "when stage return an error the status must be failed",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", nil, sdk_errors.NewStageError(errors.New("stage1")))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_FAILED), stage1Status, "'stage1' should be in 'failed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "should execute complete stage",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "complete")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Complete("complete", func(completeContext sdk_v1.CompleteContext) sdk_v1.StageError {
							sb.exec("complete")
							return sb.shouldErr("complete")
						}).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving stage status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				completeStatus, err := sph.Get("jobKey", "complete")
				assert.Nil(t, err, "error retrieving stage status for complete: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), completeStatus, "'complete' should be in 'completed' status, got: %s", completeStatus)
			},
		},
		{
			name: "should return an error if last active stage does not exist in the chain",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			lastActiveStage: &sdk_v1.LastActiveStage{Name: "non-existent"},
			errorType:       lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.False(t, sb.Executed("stage1"), "'stage1' must not have been executed")
			},
		},
		{
			name: "should skip stage1 return skip error and skip stage2 using conditional stage execution",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "stage2").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_SKIP)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Stage("stage2", stageFn("stage2", sb), spark.WithStageStatus("stage1", sdk_v1.StageStatus_STAGE_STATUS_COMPLETED)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_SKIPPED), stage1Status, "'stage1' should be in 'skipped' status, got: %s", stage1Status)
				assert.False(t, sb.Executed("stage2"), "'stage2' must NOT have been executed")
				stage2Status, err := sph.Get("jobKey", "stage2")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_SKIPPED), stage2Status, "'stage1' should be in 'skipped' status, got: %s", stage2Status)
			},
		},
		{
			name: "should fail when a error occur during conditional stage execution evaluation",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "stage2")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Stage("stage2", stageFn("stage2", sb), spark.WithStageStatus("stage1", sdk_v1.StageStatus_STAGE_STATUS_SKIPPED)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
				assert.False(t, sb.Executed("stage2"), "'stage2' must NOT have been executed")
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("stage2", sdk_v1.StageStatus_STAGE_STATUS_SKIPPED, errors.New("error updating status for stage 2"))
			},
		}, {
			name: "should fail when trying to update a stage status to starting",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.False(t, sb.Executed("stage1"), "'stage1' must have not been executed")
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("stage1", sdk_v1.StageStatus_STAGE_STATUS_STARTED, errors.New("error updating status for stage 1"))
			},
		},
		{
			name: "should fail and update the stage status to failed",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_FAILED), stage1Status, "'stage1' should be in 'failed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "if a stage fails and an error occur trying to update the stage, that error should be returned",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_CANCELED)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_STARTED), stage1Status, "'stage1' should be in 'started' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("stage1", sdk_v1.StageStatus_STAGE_STATUS_CANCELED, errors.New("error updating status for stage 1"))
			},
		},
		{
			name: "if a stage fails and the chain has a compensate node configured, it must run and returns the original error",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "compensate").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Compensate(spark.NewNode().Stage("compensate", stageFn("compensate", sb)).Build()).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("compensate"), "'compensate' must have been executed")
				stage1Status, err := sph.Get("jobKey", "compensate")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'compensate' should be in 'completed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "if a compensate stage fails, the status must be updated and the error should be returned accordingly",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "compensate").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED))).
					Change("compensate", nil,
						sdk_errors.NewStageError(errors.New("err-compensate"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_CANCELED)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Compensate(spark.NewNode().Stage("compensate", stageFn("compensate", sb)).Build()).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("compensate"), "'compensate' must have been executed")
				stageCompensate, err := sph.Get("jobKey", "compensate")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_CANCELED), stageCompensate, "'compensate' should be in 'completed' status, got: %s", stageCompensate)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_CANCELED),
		},
		{
			name: "if a stage returns a error type canceled and the chain has a cancel node configured, it must run and returns the original error",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "cancel").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_CANCELED)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Canceled(spark.NewNode().Stage("cancel", stageFn("cancel", sb)).Build()).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("cancel"), "'cancel' must have been executed")
				stage1Status, err := sph.Get("jobKey", "cancel")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'cancel' should be in 'completed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_CANCELED),
		},
		{
			name: "if a cancel stage fails, the status must be updated and the error should be returned accordingly",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "cancel").
					Change("stage1", nil,
						sdk_errors.NewStageError(errors.New("err-stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_CANCELED))).
					Change("cancel", nil,
						sdk_errors.NewStageError(errors.New("err-cancel"), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_RETRY)))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Canceled(spark.NewNode().Stage("cancel", stageFn("cancel", sb)).Build()).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("cancel"), "'cancel' must have been executed")
				stageCancel, err := sph.Get("jobKey", "cancel")
				assert.Nil(t, err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_FAILED), stageCancel, "'cancel' should be in 'completed' status, got: %s", stageCancel)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_RETRY),
		},
		{
			name: "when stage return an unsupported error it must be return error type fatal",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", nil, sdk_errors.NewStageError(errors.New("stage1"), sdk_errors.WithErrorType(sdk_v1.ErrorType(500))))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_FAILED), stage1Status, "'stage1' should be in 'failed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FATAL),
		},
		{
			name: "if a stage returns a result it must be stored properly",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", "a result", nil)
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
				raw, err := sph.GetResult("jobKey", "stage1").Raw()
				assert.Nil(t, err)
				assert.Equal(t, `"a result"`, string(raw), "result error: expected > 'a result', got: '%s'", string(raw))
			},
		},
		{
			name: "must return an error if a invalid result is returned from a stage",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", make(chan struct{}), nil)
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_FAILED), stage1Status, "'stage1' should be in 'failed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "must return an error if a invalid result is returned from a stage and we could not update the stage status",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", make(chan struct{}), nil)
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_STARTED), stage1Status, "'stage1' should be in 'started' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("stage1", sdk_v1.StageStatus_STAGE_STATUS_FAILED, errors.New("error updating status for stage 1"))
			},
		},
		{
			name: "must return an error if we can't call the api to store the stage result",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", "a valid result", nil)
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_FAILED), stage1Status, "'stage1' should be in 'failed' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().SetResult("jobKey", "stage1", errors.New("error calling set result api"))
			},
		},
		{
			name: "must return an fatal error if we can't call the api to store the stage result and the call to the update stage api also fails",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1").Change("stage1", "a valid result", nil)
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_STARTED), stage1Status, "'stage1' should be in 'started' status, got: %s", stage1Status)
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().SetResult("jobKey", "stage1", errors.New("error calling set result api"))
				sph.AddBehaviour().Set("stage1", sdk_v1.StageStatus_STAGE_STATUS_FAILED, errors.New("error updating stage1 status"))
			},
		},
		{
			name: "should fail when trying to update stage status to complete",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_STARTED), stage1Status, "'stage1' should be in 'started' status, got: %s", stage1Status)
			},
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("stage1", sdk_v1.StageStatus_STAGE_STATUS_COMPLETED, errors.New("error updating stage1 status"))
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "should fail when complete stage returns an error and we can't update complete stage status to failed",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "complete").Change("complete", nil, sdk_errors.NewStageError(errors.New("complete error")))
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Complete("complete", func(completeContext sdk_v1.CompleteContext) sdk_v1.StageError {
							sb.exec("complete")
							return sb.shouldErr("complete")
						}).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				assert.True(t, sb.Executed("complete"), "'complete' must have been executed")
				completeStatus, err := sph.Get("jobKey", "complete")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_STARTED), completeStatus, "'stage1' should be in 'started' status, got: %s", completeStatus)
			},
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("complete", sdk_v1.StageStatus_STAGE_STATUS_FAILED, errors.New("error updating stage1 status"))
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "should fail when can't update complete stage to starting",
			chainFn: func() (*spark.Chain, *stageBehaviour) {
				sb := NewStageBehaviour(t, "stage1", "complete")
				chain, err := spark.NewChain(
					spark.NewNode().
						Stage("stage1", stageFn("stage1", sb)).
						Complete("complete", func(completeContext sdk_v1.CompleteContext) sdk_v1.StageError {
							sb.exec("complete")
							return sb.shouldErr("complete")
						}).
						Build()).
					Build()
				assert.Nil(t, err, "error creating spark node chain: %v", err)
				return chain, sb
			},
			assertFn: func(t *testing.T, sb *stageBehaviour, sph sdk_v1.StageProgressHandler) {
				assert.True(t, sb.Executed("stage1"), "'stage1' must have been executed")
				stage1Status, err := sph.Get("jobKey", "stage1")
				assert.Nil(t, err, "error retrieving spark status for stage1: %v", err)
				assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_STAGE_STATUS_COMPLETED), stage1Status, "'stage1' should be in 'completed' status, got: %s", stage1Status)
			},
			prepare: func(sph *inmemory.StageProgressHandler) {
				sph.AddBehaviour().Set("complete", sdk_v1.StageStatus_STAGE_STATUS_STARTED, errors.New("error updating stage1 status"))
			},
			errorType: lo.ToPtr(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			chain, sb := test.chainFn()
			sph := inmemory.NewStageProgressHandler(t)
			if test.prepare != nil {
				test.prepare(sph)
			}
			zerolog.SetGlobalLevel(zerolog.Disabled)
			worker := v1.NewSparkTestWorker(t, chain, v1.WithIOHandler(inmemory.NewIOHandler(t)), v1.WithStageProgressHandler(sph))
			err := worker.Execute(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", test.lastActiveStage))
			if err != nil && test.errorType == nil {
				t.Errorf("a unexpected error occured: %v", err)
			}
			if test.errorType != nil {
				if err == nil {
					t.Errorf("error '%s' is expected from chain execution, got none", test.errorType)
				} else if *test.errorType != err.ErrorType() {
					t.Errorf("error expected: %v; got: %v;", test.errorType, err.ErrorType())
				}
			}
			test.assertFn(t, sb, sph)
			sph.ResetBehaviour()
		})
	}
}
