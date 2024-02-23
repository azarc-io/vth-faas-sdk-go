package module_test_runner

import (
	"context"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
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
	buildChainCalledCount int32
	stageCalledCount      int32
	completeCalledCount   int32
}

func (s *basicSpark) Init(ctx sparkv1.InitContext) error {
	return nil
}

func (s *basicSpark) Stop() {

}

func (s *basicSpark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	s.buildChainCalledCount += 1
	return b.NewChain("test-0").
		Stage("Stage-0", func(_ sparkv1.StageContext) (any, sparkv1.StageError) {
			atomic.AddInt32(&s.stageCalledCount, 1)
			return nil, nil
		}).
		Complete(func(context sparkv1.CompleteContext) sparkv1.StageError {
			atomic.AddInt32(&s.completeCalledCount, 1)
			return nil
		})
}

/************************************************************************/
// TYPES SLOW SPARK
/************************************************************************/

type slowSpark struct {
	stageCalledCount      int32
	completeCalledCount   int32
	buildChainCalledCount int32
}

func (s *slowSpark) Init(ctx sparkv1.InitContext) error {
	return nil
}

func (s *slowSpark) Stop() {

}

func (s *slowSpark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	s.buildChainCalledCount += 1
	return b.NewChain("test-0").
		Stage("Stage-0", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			ctx.Log().Info("Stage-0 called")
			atomic.AddInt32(&s.stageCalledCount, 1)
			time.Sleep(time.Millisecond * 500)
			return nil, nil
		}).
		Stage("Stage-1", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			atomic.AddInt32(&s.stageCalledCount, 1)
			ctx.Log().Info("Stage-2 called")
			return nil, nil
		}).
		Complete(func(ctx sparkv1.CompleteContext) sparkv1.StageError {
			ctx.Log().Info("Complete called")
			atomic.AddInt32(&s.completeCalledCount, 1)
			return nil
		})
}

/************************************************************************/
// TESTS
/************************************************************************/

func (s *WorkerSuite) Test_Should_Call_BuildChain_On_Registration() {
	jobKey := "call_buildchain_on_register"
	spark := new(basicSpark)
	worker := s.createWorker(spark)
	ctx := NewTestJobContext(context.Background(), jobKey, "cid", "tid", Inputs{})

	out, err := worker.Execute(ctx)
	s.Require().NoError(err)
	s.Require().NotNil(out)

	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.buildChainCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageCalledCount))
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.completeCalledCount))
}

func (s *WorkerSuite) Test_Should_Drain_Running_Stages_During_Shutdown_When_Context_Is_Cancelled() {
	oc, cancel := context.WithCancel(context.Background())
	jobKey := "drain_during_shutdown"

	spark := new(slowSpark)
	worker := s.createWorker(spark)
	done := make(chan struct{})

	go func() {
		time.Sleep(time.Millisecond * 100)
		// Cancel the context
		log.Info().Msg("cancelling")
		cancel()
	}()

	go func() {
		ctx := NewTestJobContext(oc, jobKey, "cid", "tid", Inputs{})
		_, err := worker.Execute(ctx)
		s.ErrorContains(err, "canceled")
		close(done)
	}()

	<-done

	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.buildChainCalledCount), "the SparkChain should run at least once")
	s.Require().Equal(int32(1), atomic.LoadInt32(&spark.stageCalledCount), "only the first Stage should have run")
	s.Require().Equal(int32(0), atomic.LoadInt32(&spark.completeCalledCount), "completion should not be called")
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (s *WorkerSuite) createWorker(spark sparkv1.Spark) RunnerTest {
	worker, err := NewTestRunner(s.T(), spark)
	s.Require().Nil(err)
	return worker
}

func (s *WorkerSuite) createSlowWorker(ctx context.Context) (sparkv1.Worker, *slowSpark) {
	spark := &slowSpark{}
	worker, err := sparkv1.NewSparkWorker(ctx, spark)
	s.Require().Nil(err)
	return worker, spark
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
