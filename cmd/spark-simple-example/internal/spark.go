package spark

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

type Spark struct{}

func (s Spark) Init(ctx sparkv1.InitContext) error {
	return nil
}

func (s Spark) Stop() {

}

func (s Spark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	return b.NewChain("chain-1").
		Stage("stage-1", func(_ sparkv1.StageContext) (any, sparkv1.StageError) {
			return "hello", nil
		}).
		Stage("stage-2", func(_ sparkv1.StageContext) (any, sparkv1.StageError) {
			return "world", nil
		}).
		Complete(func(ctx sparkv1.CompleteContext) sparkv1.StageError {
			var (
				stg1Res, stg2Res string
				err              error
			)

			// get the result of the 2 stages
			if err = ctx.StageResult("stage-1").Bind(&stg1Res); err != nil {
				return sparkv1.NewStageError(err)
			}
			if err = ctx.StageResult("stage-2").Bind(&stg2Res); err != nil {
				return sparkv1.NewStageError(err)
			}

			// write the output of the spark
			if err = ctx.Output(
				sparkv1.NewVar("message", "", fmt.Sprintf("%s %s", stg1Res, stg2Res)),
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
