package spark_v1

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

/************************************************************************/
// TYPES SUITE
/************************************************************************/

type WorkerSuite struct {
	suite.Suite
}

/************************************************************************/
// TYPES BASIC SPARK
/************************************************************************/

type basicSpark struct {
	buildChainCalledCount  int
	stageCalledCount       int
	completeCalledCount    int
	stageDelegatedCount    int
	completeDelegatedCount int
	delegatedStageNames    []string
	delegatedCompleteNames []string
}

func (s *basicSpark) BuildChain(b Builder) Chain {
	s.buildChainCalledCount += 1
	return b.NewChain("test-0").
		Stage("stage-0", func(_ StageContext) (any, StageError) {
			s.stageCalledCount += 1
			return nil, nil
		}).
		Complete(func(context CompleteContext) StageError {
			s.completeCalledCount += 1
			return nil
		})
}

func (s *basicSpark) delegateStage(ctx StageContext, cb StageDefinitionFn) (any, StageError) {
	s.stageDelegatedCount += 1
	s.delegatedStageNames = append(s.delegatedStageNames, ctx.Name())
	return cb(ctx)
}

func (s *basicSpark) delegateCompletion(ctx CompleteContext, cb CompleteDefinitionFn) StageError {
	s.completeDelegatedCount += 1
	s.delegatedCompleteNames = append(s.delegatedCompleteNames, ctx.Name())
	return cb(ctx)
}

/************************************************************************/
// TYPES SLOW SPARK
/************************************************************************/

type slowSpark struct {
	stageCalledCount      int
	completeCalledCount   int
	buildChainCalledCount int
}

func (s *slowSpark) BuildChain(b Builder) Chain {
	s.buildChainCalledCount += 1
	return b.NewChain("test-0").
		Stage("stage-0", func(ctx StageContext) (any, StageError) {
			s.stageCalledCount += 1
			time.Sleep(time.Millisecond * 200)
			ctx.Log().Info("stage-0 called")
			return nil, nil
		}).
		Stage("stage-1", func(ctx StageContext) (any, StageError) {
			s.stageCalledCount += 1
			ctx.Log().Info("stage-2 called")
			return nil, nil
		}).
		Complete(func(ctx CompleteContext) StageError {
			ctx.Log().Info("complete called")
			s.completeCalledCount += 1
			return nil
		})
}

/************************************************************************/
// TESTS
/************************************************************************/

func (s *WorkerSuite) Test_Should_Call_BuildChain_On_Registration() {
	jobKey := "test"
	worker, _, _, spark := s.createWorker(context.Background())

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(1, spark.buildChainCalledCount)
	s.Require().Equal(1, spark.stageCalledCount)
	s.Require().Equal(1, spark.completeCalledCount)
}

func (s *WorkerSuite) Test_Should_Delegate_Stage_Execution_If_Option_Provided() {
	jobKey := "test"
	worker, _, _, spark := s.createWorker(context.Background())

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(1, spark.buildChainCalledCount)
	s.Require().Equal(1, spark.stageCalledCount)
	s.Require().Equal(1, spark.completeCalledCount)
	s.Require().Equal(1, spark.stageDelegatedCount)
	s.Require().Equal(0, spark.completeDelegatedCount)
	s.Require().Equal([]string{"stage-0"}, spark.delegatedStageNames)
}

func (s *WorkerSuite) Test_Should_Delegate_Completion_If_Option_Provided() {
	jobKey := "test"
	worker, _, _, spark := s.createWorker(context.Background())

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(1, spark.buildChainCalledCount)
	s.Require().Equal(1, spark.stageCalledCount)
	s.Require().Equal(1, spark.completeCalledCount)
	s.Require().Equal(0, spark.stageDelegatedCount)
	s.Require().Equal(1, spark.completeDelegatedCount)
	s.Require().Equal([]string{"test-0_complete"}, spark.delegatedCompleteNames)
}

func (s *WorkerSuite) Test_Should_Drain_Running_Stages_During_Shutdown_When_Context_Is_Cancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	jobKey := "test"
	worker, _, _, spark := s.createWorker(ctx)

	go func() {
		time.Sleep(time.Millisecond * 50)
		// cancel the context
		cancel()
	}()

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	// only stage 0 should have executed correctly while stage 1 did not execute

	s.Require().Equal(1, spark.buildChainCalledCount)
	s.Require().Equal(1, spark.stageCalledCount)
	s.Require().Equal(0, spark.completeCalledCount)
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (s *WorkerSuite) createWorker(ctx context.Context) (Worker, *InMemoryStageProgressHandler, IOHandler, *basicSpark) {
	spark := &basicSpark{}
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())

	worker, err := NewSparkWorker(ctx, spark,
		WithIOHandler(vh),
		WithStageProgressHandler(sph),
		WithDelegateCompletion(spark.delegateCompletion),
		WithLog(NewLogger()),
	)

	s.Require().Nil(err)

	return worker, sph, vh, spark
}

func (s *WorkerSuite) createSlowWorker(ctx context.Context) (Worker, *InMemoryStageProgressHandler, IOHandler, *slowSpark) {
	spark := &slowSpark{}
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())

	worker, err := NewSparkWorker(ctx, spark,
		WithIOHandler(vh),
		WithStageProgressHandler(sph),
		WithLog(NewLogger()),
	)

	s.Require().Nil(err)

	return worker, sph, vh, spark
}

/************************************************************************/
// LIFECYCLE
/************************************************************************/

func (s *WorkerSuite) TearDownSuite() {

}

func (s *WorkerSuite) SetupSuite() {

}

/************************************************************************/
// SUITE
/************************************************************************/

func TestWorkerSuite(t *testing.T) {
	suite.Run(t, new(WorkerSuite))
}
