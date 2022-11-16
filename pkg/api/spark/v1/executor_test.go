package sdk_v1

import (
	"context"
	"github.com/stretchr/testify/suite"
	"testing"
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

	c := b.BuildChain()
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
