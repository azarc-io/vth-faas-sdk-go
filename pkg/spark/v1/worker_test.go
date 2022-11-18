package spark_v1

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
	"sync/atomic"
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
	buildChainCalledCount  int32
	stageCalledCount       int32
	completeCalledCount    int32
	stageDelegatedCount    int32
	completeDelegatedCount int32
	delegatedStageNames    []string
	delegatedCompleteNames []string
}

func (s *basicSpark) BuildChain(b Builder) Chain {
	s.buildChainCalledCount += 1
	return b.NewChain("test-0").
		Stage("stage-0", func(_ StageContext) (any, StageError) {
			atomic.AddInt32(&s.stageCalledCount, 1)
			return nil, nil
		}).
		Complete(func(context CompleteContext) StageError {
			atomic.AddInt32(&s.completeCalledCount, 1)
			return nil
		})
}

func (s *basicSpark) delegateStage(ctx StageContext, cb StageDefinitionFn) (any, StageError) {
	atomic.AddInt32(&s.stageDelegatedCount, 1)
	s.delegatedStageNames = append(s.delegatedStageNames, ctx.Name())
	return cb(ctx)
}

func (s *basicSpark) delegateCompletion(ctx CompleteContext, cb CompleteDefinitionFn) StageError {
	atomic.AddInt32(&s.completeDelegatedCount, 1)
	s.delegatedCompleteNames = append(s.delegatedCompleteNames, ctx.Name())
	return cb(ctx)
}

/************************************************************************/
// TYPES SLOW SPARK
/************************************************************************/

type slowSpark struct {
	stageCalledCount      int32
	completeCalledCount   int32
	buildChainCalledCount int32
}

func (s *slowSpark) BuildChain(b Builder) Chain {
	s.buildChainCalledCount += 1
	return b.NewChain("test-0").
		Stage("stage-0", func(ctx StageContext) (any, StageError) {
			ctx.Log().Info("stage-0 called")
			atomic.AddInt32(&s.stageCalledCount, 1)
			time.Sleep(time.Millisecond * 500)
			return nil, nil
		}).
		Stage("stage-1", func(ctx StageContext) (any, StageError) {
			atomic.AddInt32(&s.stageCalledCount, 1)
			ctx.Log().Info("stage-2 called")
			return nil, nil
		}).
		Complete(func(ctx CompleteContext) StageError {
			ctx.Log().Info("complete called")
			atomic.AddInt32(&s.completeCalledCount, 1)
			return nil
		})
}

/************************************************************************/
// TESTS
/************************************************************************/

func (s *WorkerSuite) Test_Should_Call_BuildChain_On_Registration() {
	jobKey := "test"
	worker, _, _, spark := s.createWorker(context.Background(), true, true)

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.buildChainCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.completeCalledCount))
}

func (s *WorkerSuite) Test_Should_Delegate_Stage_Execution_If_Option_Provided() {
	jobKey := "test"
	worker, _, _, spark := s.createWorker(context.Background(), true, false)

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.buildChainCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.completeCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageDelegatedCount))
	s.Require().Equal(int32(0), atomic.LoadInt32(&spark.completeDelegatedCount))
	s.Require().Equal([]string{"stage-0"}, spark.delegatedStageNames)
}

func (s *WorkerSuite) Test_Should_Delegate_Completion_If_Option_Provided() {
	jobKey := "test"
	worker, _, _, spark := s.createWorker(context.Background(), false, true)

	err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.buildChainCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.completeCalledCount))
	s.Require().Equal(int32(0), atomic.LoadInt32(&spark.stageDelegatedCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.completeDelegatedCount))
	s.Require().Equal([]string{"test-0_complete"}, spark.delegatedCompleteNames)
}

func (s *WorkerSuite) Test_Should_Drain_Running_Stages_During_Shutdown_When_Context_Is_Cancelled() {
	ctx, cancel := context.WithCancel(context.Background())
	jobKey := "test"
	worker, _, _, spark := s.createSlowWorker(ctx)
	done := make(chan struct{})

	go func() {
		worker.Run()
		close(done)
	}()

	go func() {
		time.Sleep(time.Millisecond * 100)
		// cancel the context
		log.Info().Msg("cancelling")
		cancel()
	}()

	go func() {
		err := worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
		s.Require().Nil(err)
	}()

	<-done

	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.buildChainCalledCount), "the chain should run at least once")
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageCalledCount), "only the first stage should have run")
	s.Require().Equal(int32(0), atomic.LoadInt32(&spark.completeCalledCount), "completion should not be called")
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (s *WorkerSuite) createWorker(ctx context.Context, delStage, delComplete bool) (Worker, *InMemoryStageProgressHandler, IOHandler, *basicSpark) {
	spark := &basicSpark{}
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())

	var options []Option
	options = append(options, WithIOHandler(vh), WithStageProgressHandler(sph), WithLog(NewLogger()))

	if delStage {
		options = append(options, WithDelegateStage(spark.delegateStage))
	}

	if delComplete {
		options = append(options, WithDelegateCompletion(spark.delegateCompletion))
	}

	worker, err := NewSparkWorker(ctx, spark, options...)

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
