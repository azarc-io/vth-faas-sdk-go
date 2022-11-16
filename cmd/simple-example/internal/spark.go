package spark

import (
	"fmt"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

type Spark struct {
}

func (s Spark) BuildChain(b sdk_v1.Builder) sdk_v1.ChainNodeFinalizer {
	return b.NewChain("chain-1").
		Stage("stage-1", func(_ sdk_v1.StageContext) (any, sdk_v1.StageError) {
			return "hello", nil
		}).
		Stage("stage-2", func(_ sdk_v1.StageContext) (any, sdk_v1.StageError) {
			return "world", nil
		}).
		Complete(func(ctx sdk_v1.CompleteContext) sdk_v1.StageError {
			var (
				stg1Res, stg2Res string
				err              error
			)

			// get the result of the 2 stages
			if err = ctx.StageResult("stage-1").Bind(&stg1Res); err != nil {
				return sdk_v1.NewStageError(err)
			}
			if err = ctx.StageResult("stage-2").Bind(&stg2Res); err != nil {
				return sdk_v1.NewStageError(err)
			}

			// write the output of the spark
			if err = ctx.Output(
				sdk_v1.NewVar("message", "", fmt.Sprintf("%s %s", stg1Res, stg2Res)),
			); err != nil {
				return sdk_v1.NewStageError(err)
			}

			return nil
		})
}

func NewSpark() *Spark {
	return &Spark{}
}
