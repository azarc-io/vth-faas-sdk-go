package sparkv1

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

/************************************************************************/
// TYPES
/************************************************************************/

type BuilderSuite struct {
	suite.Suite
}

/************************************************************************/
// TESTS
/************************************************************************/

func (s *BuilderSuite) Test_Should_Create_Root_Node_With_No_Children() {
	b := NewBuilder()
	n := b.NewChain("test-0").
		Stage("Stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Complete(func(context CompleteContext) StageError {
			return nil
		}).
		build()

	s.Require().True(n.HasCompletionStage(), "must have a completion Stage")
	s.Require().False(n.HasCompensationStage(), "must not have a compensation Stage")
	s.Require().False(n.HasCancellationStage(), "must not have a cancellation Stage")
	s.Equal("test-0", n.ChainName())
	s.Equal("test-0_complete", n.CompletionName())
	s.Equal(1, n.CountOfStages(), "must have only 1 Stage")
	s.True(n.IsRoot())
	s.False(n.IsCompensate())
	s.False(n.IsCancel())
}

func (s *BuilderSuite) Test_Should_Create_Root_Node_With_Child_Node() {
	b := NewBuilder()
	n := b.NewChain("test-0").
		Stage("Stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Compensate(
			b.NewChain("test-1").
				Stage("Stage-1", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Cancelled(
			b.NewChain("test-3").
				Stage("Stage-3", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Complete(func(context CompleteContext) StageError {
			return nil
		}).
		build()

	// generate report for validation
	r := generateReportForChain(b.BuildChain())
	s.Require().NotNil(r)
	s.Require().Len(r.Errors, 0)

	s.Require().True(n.HasCompletionStage(), "must have a completion Stage")
	s.Require().True(n.HasCompensationStage(), "must not have a compensation Stage")
	s.Require().True(n.HasCancellationStage(), "must have a cancellation Stage")
	s.Equal("test-0", n.ChainName())
	s.Equal("test-0_complete", n.CompletionName())
	s.Equal(1, n.CountOfStages(), "must have only 1 Stage")
}

func (s *BuilderSuite) Test_Report_Should_Generate_Single_Error_On_Single_Duplicate_Stage_Names() {
	b := NewBuilder()
	n := b.NewChain("test-0").
		Stage("Stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Compensate(
			b.NewChain("test-1").
				Stage("Stage-0", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Cancelled(
			b.NewChain("test-3").
				Stage("Stage-3", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Complete(func(context CompleteContext) StageError {
			return nil
		}).
		build()

	// generate report for validation
	r := generateReportForChain(b.BuildChain())
	s.Require().NotNil(r)
	s.Require().Len(r.Errors, 1)
	s.Equal("duplicate Stage names are not permitted [SparkChain]: Stage-0 [at]: root > Compensate", r.Errors[0].Error())

	s.Require().True(n.HasCompletionStage(), "must have a completion Stage")
	s.Require().True(n.HasCompensationStage(), "must not have a compensation Stage")
	s.Require().True(n.HasCancellationStage(), "must have a cancellation Stage")
	s.Equal("test-0", n.ChainName())
	s.Equal("test-0_complete", n.CompletionName())
	s.Equal(1, n.CountOfStages(), "must have only 1 Stage")
}

func (s *BuilderSuite) Test_Report_Should_Generate_Multiple_Error_On_Multiple_Duplicate_Stage_Names() {
	b := NewBuilder()
	n := b.NewChain("test-0").
		Stage("Stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Compensate(
			b.NewChain("test-1").
				Stage("Stage-0", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Cancelled(
			b.NewChain("test-3").
				Stage("Stage-0", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Complete(func(context CompleteContext) StageError {
			return nil
		}).
		build()

	// generate report for validation
	r := generateReportForChain(b.BuildChain())
	s.Require().NotNil(r)
	s.Require().Len(r.Errors, 2)
	s.Equal("duplicate Stage names are not permitted [SparkChain]: Stage-0 [at]: root > Compensate", r.Errors[0].Error())
	s.Equal("duplicate Stage names are not permitted [SparkChain]: Stage-0 [at]: root > canceled", r.Errors[1].Error())

	s.Require().True(n.HasCompletionStage(), "must have a completion Stage")
	s.Require().True(n.HasCompensationStage(), "must not have a compensation Stage")
	s.Require().True(n.HasCancellationStage(), "must have a cancellation Stage")
	s.Equal("test-0", n.ChainName())
	s.Equal("test-0_complete", n.CompletionName())
	s.Equal(1, n.CountOfStages(), "must have only 1 Stage")
}

func (s *BuilderSuite) Test_Report_Should_Generate_Errors_On_Duplicate_Chain_Names() {
	b := NewBuilder()
	n := b.NewChain("test-0").
		Stage("Stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Compensate(
			b.NewChain("test-0").
				Stage("Stage-1", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Cancelled(
			b.NewChain("test-0").
				Stage("Stage-2", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Complete(func(context CompleteContext) StageError {
			return nil
		}).
		build()

	// generate report for validation
	r := generateReportForChain(b.BuildChain())
	s.Require().NotNil(r)
	s.Require().Len(r.Errors, 2)
	s.Equal("duplicate SparkChain names are not permitted [Name]: test-0 [at]: root > Compensate", r.Errors[0].Error())
	s.Equal("duplicate SparkChain names are not permitted [Name]: test-0 [at]: root > canceled", r.Errors[1].Error())

	s.Require().True(n.HasCompletionStage(), "must have a completion Stage")
	s.Require().True(n.HasCompensationStage(), "must not have a compensation Stage")
	s.Require().True(n.HasCancellationStage(), "must have a cancellation Stage")
	s.Equal("test-0", n.ChainName())
	s.Equal("test-0_complete", n.CompletionName())
	s.Equal(1, n.CountOfStages(), "must have only 1 Stage")
}

func (s *BuilderSuite) Test_Report_Should_Generate_Errors_On_Empty_Names() {
	b := NewBuilder()
	n := b.NewChain("").
		Stage("Stage-0", func(_ StageContext) (any, StageError) {
			return nil, nil
		}).
		Compensate(
			b.NewChain("test-0").
				Stage("", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Cancelled(
			b.NewChain("test-1").
				Stage("Stage-2", func(context StageContext) (any, StageError) {
					return nil, nil
				}).
				Complete(func(context CompleteContext) StageError {
					return nil
				}),
		).
		Complete(func(context CompleteContext) StageError {
			return nil
		}).
		build()

	// generate report for validation
	r := generateReportForChain(b.BuildChain())
	s.Require().NotNil(r)
	s.Require().Len(r.Errors, 2)
	s.Equal("SparkChain Name can not be empty [at]: root", r.Errors[0].Error())
	s.Equal("Stage Name can not be empty [at]: root > Compensate", r.Errors[1].Error())

	s.Require().True(n.HasCompletionStage(), "must have a completion Stage")
	s.Require().True(n.HasCompensationStage(), "must not have a compensation Stage")
	s.Require().True(n.HasCancellationStage(), "must have a cancellation Stage")
	s.Equal("", n.ChainName())
	s.Equal("_complete", n.CompletionName())
	s.Equal(1, n.CountOfStages(), "must have only 1 Stage")
}

/************************************************************************/
// LIFECYCLE
/************************************************************************/

func (s *BuilderSuite) TearDownSuite() {

}

func (s *BuilderSuite) SetupSuite() {

}

/************************************************************************/
// SUITE
/************************************************************************/

func TestRoutingSuite(t *testing.T) {
	suite.Run(t, new(BuilderSuite))
}
