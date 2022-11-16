//go:build e2e

package spark_v1

import (
	"context"
	"errors"
	"sync"
	"time"
)

func (s *ExecutorSuite) Test_Should_Retry_Stage_If_Stage_Returns_Retry_Option() {
	wg := sync.WaitGroup{}
	wg.Add(3)

	b := NewBuilder()
	b.NewChain("test-0").
		Stage("stage-0", func(ctx StageContext) (any, StageError) {
			wg.Done()
			ctx.Log().Info("stage-0 called")
			return nil, NewStageError(errors.New("unstable"), WithRetry(1, time.Millisecond*10))
		}).
		Compensate(
			b.NewChain("compensate").
				Stage("stage-1", func(ctx StageContext) (any, StageError) {
					wg.Done()
					ctx.Log().Info("compensate stage-1 called")
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

	if WaitTimeout(&wg, time.Second) {
		s.FailNow("time out waiting for all steps to complete")
	}
}
