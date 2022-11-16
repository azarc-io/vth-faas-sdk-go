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
	buildChainCalledCount int
	stageCalledCount      int
	completeCalledCount   int
}

func (s *basicSpark) BuildChain(b Builder) ChainNodeFinalizer {
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
