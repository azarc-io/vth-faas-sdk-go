package spark

import (
	"context"
	"fmt"
	spark_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

type Spark struct {
	shutdown context.CancelFunc
}

func (s Spark) BuildChain(b spark_v1.Builder) spark_v1.Chain {
	return b.NewChain("chain-1").
		Stage("stage-1", func(ctx spark_v1.StageContext) (any, spark_v1.StageError) {
			var name string
			if err := ctx.Input("name").Bind(&name); err != nil {
				return nil, spark_v1.NewStageError(err)
			}
			return fmt.Sprintf("Hello, %s", name), nil
		}).
		Complete(func(ctx spark_v1.CompleteContext) spark_v1.StageError {
			var (
				message string
				err     error
			)

			// fetch the output of stage-1
			if err = ctx.StageResult("stage-1").Bind(&message); err != nil {
				return spark_v1.NewStageError(err)
			}

			// write the output of the spark
			if err = ctx.Output(
				spark_v1.NewVar("message", "", message),
			); err != nil {
				return spark_v1.NewStageError(err)
			}

			return nil
		})
}

// NewSpark creates a Spark
// shutdown is a reference to the context cancel function, it can be used to gracefully stop the worker if needed
func NewSpark(shutdown context.CancelFunc) spark_v1.Spark {
	return &Spark{shutdown: shutdown}
}
