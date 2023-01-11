package spark

import (
	"fmt"

	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

type Config struct {
	Retry *sparkv1.RetryConfig `json:"retry"`
	Panic bool                 `json:"panic"`
}

type Spark struct {
	cfg            *Config
	totalFailures  int
	currentFailure int
}

func (s *Spark) Init(ctx sparkv1.InitContext) error {
	s.cfg = new(Config)
	if err := ctx.Config().Bind(s.cfg); err != nil {
		return err
	}
	return nil
}

func (s *Spark) Stop() {

}

func (s *Spark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	return b.NewChain("chain-1").
		Stage("stage-1", func(_ sparkv1.StageContext) (any, sparkv1.StageError) {
			if s.cfg.Panic {
				panic("i was forced to panic :)")
			}
			return nil, nil
		}).
		Stage("stage-2", func(_ sparkv1.StageContext) (any, sparkv1.StageError) {
			if s.totalFailures <= s.currentFailure {
				// allow response to pass now
				return fmt.Sprintf("finally I can pass after %d failures", s.currentFailure), nil
			}
			s.currentFailure++
			return nil, sparkv1.NewStageError(
				fmt.Errorf("failures %d of %d", s.currentFailure, s.totalFailures),
				sparkv1.WithRetry(s.cfg.Retry.Times, s.cfg.Retry.BackoffMultiplier, s.cfg.Retry.FirstBackoffWait),
			)
		}).
		Complete(func(ctx sparkv1.CompleteContext) sparkv1.StageError {
			var (
				stg1Res string
				err     error
			)

			// get the result of the stages
			if err = ctx.StageResult("stage-2").Bind(&stg1Res); err != nil {
				return sparkv1.NewStageError(err)
			}

			// write the output of the spark
			if err = ctx.Output(
				sparkv1.NewVar("message", "application/text", stg1Res),
			); err != nil {
				return sparkv1.NewStageError(err)
			}

			return nil
		})
}

// NewSpark creates a Spark
func NewSpark(totalFailures int) sparkv1.Spark {
	return &Spark{totalFailures: totalFailures}
}
