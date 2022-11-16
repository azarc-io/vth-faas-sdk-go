package sdk_v1

import (
	"context"
	"errors"
	"github.com/stretchr/testify/suite"
	"strconv"
	"sync"
	"testing"
	"time"
)

/************************************************************************/
// TYPES
/************************************************************************/

type ExecutorSuite struct {
	suite.Suite
}

/************************************************************************/
// TESTS
/************************************************************************/

func (s *ExecutorSuite) Test_Execute_Single_Stage_Then_Complete_With_No_Logic_Should_Not_Error() {
	b := NewBuilder()
	b.NewChain("test-0").
		Stage("stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Complete(func(context CompleteContext) StageError {
			return nil
		})

	c := b.buildChain()
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())
	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
	jobContext := NewJobContext(metadata, sph, vh, NewLogger())

	err := c.Execute(jobContext)

	s.Require().Nil(err)
	sph.AssertStageStatus(jobKey, "test-0_complete", StageStatus_STAGE_STATUS_COMPLETED)
	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_COMPLETED)
}

// TODO confirm this with Jono
//func (s *ExecutorSuite) Test_Fetch_Nil_Stage_Result_Should_Return_Nil() {
//	b := NewBuilder()
//	b.NewChain("test-0").
//		Stage("stage-0", func(_ StageContext) (any, StageError) {
//			return nil, nil
//		}).
//		Complete(func(ctx CompleteContext) StageError {
//			raw, err := ctx.StageResult("stage-0").Raw()
//			s.Require().Nil(err)
//			s.Nil(raw)
//			return nil
//		})
//
//	c := b.buildChain()
//	jobKey := "test"
//	sph := NewInMemoryStageProgressHandler(s.T())
//	vh := NewInMemoryIOHandler(s.T())
//	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
//	jobContext := NewJobContext(metadata, sph, vh, NewLogger())
//
//	err := c.Execute(jobContext)
//
//	s.Require().Nil(err)
//	sph.AssertStageStatus(jobKey, "test-0_complete", StageStatus_STAGE_STATUS_COMPLETED)
//	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_COMPLETED)
//}

func (s *ExecutorSuite) Test_Complete_Can_Fetch_String_Stage_Result() {
	b := NewBuilder()
	b.NewChain("test-0").
		Stage("stage-0", func(_ StageContext) (any, StageError) {
			return "test", nil
		}).
		Complete(func(ctx CompleteContext) StageError {
			// Test Raw
			raw, err := ctx.StageResult("stage-0").Raw()
			s.Require().Nil(err)
			s.Equal("test", string(raw))

			// Test Bind
			var res string
			err = ctx.StageResult("stage-0").Bind(&res)
			s.Require().Nil(err)
			s.Equal("test", res)
			return nil
		})

	c := b.buildChain()
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())
	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
	jobContext := NewJobContext(metadata, sph, vh, NewLogger())

	err := c.Execute(jobContext)

	s.Require().Nil(err)
	sph.AssertStageStatus(jobKey, "test-0_complete", StageStatus_STAGE_STATUS_COMPLETED)
	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_COMPLETED)
}

func (s *ExecutorSuite) Test_Complete_Can_Fetch_Numeric_Stage_Result() {
	b := NewBuilder()
	b.NewChain("test-0").
		Stage("stage-0", func(_ StageContext) (any, StageError) {
			return 1, nil
		}).
		Stage("stage-1", func(_ StageContext) (any, StageError) {
			return -1, nil
		}).
		Complete(func(ctx CompleteContext) StageError {
			// Test Raw
			raw, err := ctx.StageResult("stage-0").Raw()
			byteToInt, _ := strconv.Atoi(string(raw))
			s.Require().Nil(err)
			s.Equal(1, byteToInt)
			raw, err = ctx.StageResult("stage-1").Raw()
			byteToInt, _ = strconv.Atoi(string(raw))
			s.Require().Nil(err)
			s.Equal(-1, byteToInt)

			// Test Bind
			var res int
			err = ctx.StageResult("stage-0").Bind(&res)
			s.Require().Nil(err)
			s.Equal(1, res)
			err = ctx.StageResult("stage-1").Bind(&res)
			s.Require().Nil(err)
			s.Equal(-1, res)
			return nil
		})

	c := b.buildChain()
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())
	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
	jobContext := NewJobContext(metadata, sph, vh, NewLogger())

	err := c.Execute(jobContext)

	s.Require().Nil(err)
	sph.AssertStageStatus(jobKey, "test-0_complete", StageStatus_STAGE_STATUS_COMPLETED)
	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_COMPLETED)
}

func (s *ExecutorSuite) Test_Should_Compensate_If_Stage_Return_Error() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	b := NewBuilder()
	b.NewChain("test-0").
		Stage("stage-0", func(_ StageContext) (any, StageError) {
			return nil, NewStageError(errors.New("unstable"))
		}).
		Compensate(
			b.NewChain("compensate").
				Stage("stage-1", func(ctx StageContext) (any, StageError) {
					wg.Done()
					return "compensated", nil
				}).
				Complete(CompleteSuccess),
		).
		Complete(CompleteError)

	c := b.buildChain()
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())
	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
	jobContext := NewJobContext(metadata, sph, vh, NewLogger())

	err := c.Execute(jobContext)

	s.Require().NotNil(err)
	s.Require().Equal("unstable", err.Error())
	sph.AssertStageStatus(jobKey, "compensate_complete", StageStatus_STAGE_STATUS_COMPLETED)
	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_FAILED)

	if WaitTimeout(&wg, time.Second) {
		s.FailNow("time out waiting for compensate")
	}
}

// TODO turn this back on when cancellation is implemented
//func (s *ExecutorSuite) Test_Should_Cancel_If_Stage_Return_Error() {
//	wg := sync.WaitGroup{}
//	wg.Add(1)
//
//	b := NewBuilder()
//	b.NewChain("test-0").
//		Stage("stage-0", func(ctx StageContext) (any, StageError) {
//			time.Sleep(time.Second)
//			return nil, NewStageError(errors.New("unstable"))
//		}).
//		Cancelled(
//			b.NewChain("cancel").
//				Stage("stage-1", func(ctx StageContext) (any, StageError) {
//					wg.Done()
//					return "cancelled", nil
//				}).
//				Complete(CompleteSuccess),
//		).
//		Complete(CompleteError)
//
//	c := b.buildChain()
//	jobKey := "test"
//	sph := NewInMemoryStageProgressHandler(s.T())
//	vh := NewInMemoryIOHandler(s.T())
//	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
//	jobContext := NewJobContext(metadata, sph, vh, NewLogger())
//
//	err := c.Execute(jobContext)
//
//	s.Require().NotNil(err)
//	s.Require().Equal("unstable", err.Error())
//	sph.AssertStageStatus(jobKey, "compensate_complete", StageStatus_STAGE_STATUS_COMPLETED)
//	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_FAILED)
//
//	WaitTimeout(&wg, time.Second)
//}

func (s *ExecutorSuite) Test_Should_Skip_Stage_If_Stage_Returns_Skip_Option() {
	wg := sync.WaitGroup{}
	wg.Add(2)

	b := NewBuilder()
	b.NewChain("test-0").
		Stage("stage-0", func(ctx StageContext) (any, StageError) {
			wg.Done()
			ctx.Log().Info("stage-0 called")
			return nil, NewStageError(errors.New("unstable"), WithSkip())
		}).
		Stage("stage-1", func(ctx StageContext) (any, StageError) {
			wg.Done()
			ctx.Log().Info("stage-1 called")
			return "pass", nil
		}).
		Compensate(
			b.NewChain("compensate").
				Stage("stage-2", func(ctx StageContext) (any, StageError) {
					s.Require().FailNow("compensate should not be called")
					return nil, nil
				}).
				Complete(CompleteSuccess),
		).
		Cancelled(
			b.NewChain("cancel").
				Stage("stage-3", func(ctx StageContext) (any, StageError) {
					s.Require().FailNow("cancelled should not be called")
					return nil, nil
				}).
				Complete(CompleteSuccess),
		).
		Complete(CompleteSuccess)

	c := b.buildChain()
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())
	metadata := NewSparkMetadata(context.Background(), jobKey, "cid", "tid", nil)
	jobContext := NewJobContext(metadata, sph, vh, NewLogger())

	err := c.Execute(jobContext)
	s.Require().Nil(err)
	sph.AssertStageStatus(jobKey, "stage-0", StageStatus_STAGE_STATUS_SKIPPED)
	sph.AssertStageStatus(jobKey, "stage-1", StageStatus_STAGE_STATUS_COMPLETED)
	sph.AssertStageStatus(jobKey, "test-0_complete", StageStatus_STAGE_STATUS_COMPLETED)

	if WaitTimeout(&wg, time.Second) {
		s.FailNow("time out waiting for all steps to complete")
	}
}

/************************************************************************/
// LIFECYCLE
/************************************************************************/

func (s *ExecutorSuite) TearDownSuite() {

}

func (s *ExecutorSuite) SetupSuite() {

}

/************************************************************************/
// SUITE
/************************************************************************/

func TestExecutorSuite(t *testing.T) {
	suite.Run(t, new(ExecutorSuite))
}
