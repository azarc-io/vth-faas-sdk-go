package spark_v1

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/************************************************************************/
// TYPES
/************************************************************************/

type WorkerSuite struct {
	suite.Suite
}

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
// TESTS
/************************************************************************/

func (s *WorkerSuite) Test_Should_Call_BuildChain_On_Registration() {
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())

	spark := &basicSpark{}
	worker, err := NewSparkWorker(nil, spark, WithIOHandler(vh), WithStageProgressHandler(sph))
	s.Require().Nil(err)

	err = worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(1, spark.buildChainCalledCount)
	s.Require().Equal(1, spark.stageCalledCount)
	s.Require().Equal(1, spark.completeCalledCount)
}

func (s *WorkerSuite) Test_Should_Delegate_Stage_Execution_If_Option_Provided() {
	jobKey := "test"
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())

	spark := &basicSpark{}
	worker, err := NewSparkWorker(nil, spark,
		WithIOHandler(vh), WithStageProgressHandler(sph), WithDelegateStage(spark.delegateStage))
	s.Require().Nil(err)

	err = worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
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
	sph := NewInMemoryStageProgressHandler(s.T())
	vh := NewInMemoryIOHandler(s.T())

	spark := &basicSpark{}
	worker, err := NewSparkWorker(nil, spark,
		WithIOHandler(vh), WithStageProgressHandler(sph), WithDelegateCompletion(spark.delegateCompletion))
	s.Require().Nil(err)

	err = worker.Execute(worker.LocalContext(jobKey, "cid", "tid"))
	s.Require().Nil(err)

	s.Require().Equal(1, spark.buildChainCalledCount)
	s.Require().Equal(1, spark.stageCalledCount)
	s.Require().Equal(1, spark.completeCalledCount)
	s.Require().Equal(0, spark.stageDelegatedCount)
	s.Require().Equal(1, spark.completeDelegatedCount)
	s.Require().Equal([]string{"test-0_complete"}, spark.delegatedCompleteNames)
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
