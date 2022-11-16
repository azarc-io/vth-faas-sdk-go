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
			// get the result of the 2 stages
			stage1Result, err := ctx.StageResult("stage-1").Raw()
			if err != nil {
				return sdk_v1.NewStageError(err)
			}
			stage2Result, err := ctx.StageResult("stage-2").Raw()
			if err != nil {
				return sdk_v1.NewStageError(err)
			}

			// write the output of the spark
			err = ctx.Output(
				sdk_v1.NewVar("message", "", fmt.Sprintf("%s %s",
					string(stage1Result), string(stage2Result))),
			)
			if err != nil {
				return sdk_v1.NewStageError(err)
			}

			return nil
		})
}

func NewSpark() *Spark {
	return &Spark{}
}
