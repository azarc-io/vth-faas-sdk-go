package spark

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/rs/zerolog/log"
)

type Config struct {
	Foo string `json:"foo"`
}

type Spark struct {
	cfg *Config
}

func (s *Spark) Init(ctx sparkv1.InitContext) error {
	s.cfg = new(Config)
	if err := ctx.Config().Bind(s.cfg); err != nil {
		return err
	}
	log.Info().Fields(s.cfg).Msgf("config fields: ")
	return nil
}

func (s *Spark) Stop() {

}

func (s *Spark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	return b.NewChain("chain-1").
		Stage("stage-1", func(_ sparkv1.StageContext) (any, sparkv1.StageError) {
			return "hello", nil
		}).
		Stage("stage-2", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			ctx.Input("test")
			return "world", nil
		}).
		Stage("stage-3", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			return []byte("with bytes"), nil
		}).
		Stage("stage-4", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			// Test config is bound
			return s.cfg.Foo, nil
		}).
		Stage("stage-5", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			// Test config is bound
			return fmt.Sprintf("JobKey:%s; TransactionId:%s; CorrelationId:%s", ctx.JobKey(), ctx.TransactionID(), ctx.CorrelationID()), nil
		}).
		Complete(func(ctx sparkv1.CompleteContext) sparkv1.StageError {
			var (
				stg1Res, stg2Res, stg5Res string
				stg3Res                   []byte
				err                       error
			)

			// get the result of the 3 stages
			if err = ctx.StageResult("stage-1").Bind(&stg1Res); err != nil {
				return sparkv1.NewStageError(err)
			}
			if err = ctx.StageResult("stage-2").Bind(&stg2Res); err != nil {
				return sparkv1.NewStageError(err)
			}
			if err = ctx.StageResult("stage-3").Bind(&stg3Res); err != nil {
				return sparkv1.NewStageError(err)
			}
			if err = ctx.StageResult("stage-5").Bind(&stg5Res); err != nil {
				return sparkv1.NewStageError(err)
			}

			// write the output of the spark
			if err = ctx.Output(
				sparkv1.NewVar("message", "application/text", fmt.Sprintf("%s %s %s", stg1Res, stg2Res, string(stg3Res))),
				sparkv1.NewVar("contextValues", "application/text", stg5Res),
			); err != nil {
				return sparkv1.NewStageError(err)
			}

			return nil
		})
}

// NewSpark creates a Spark
// shutdown is a reference to the context cancel function, it can be used to gracefully stop the worker if needed
func NewSpark() sparkv1.Spark {
	return &Spark{}
}
