package spark

import (
	"context"
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

type Spark struct {
	shutdown context.CancelFunc
}

func (s Spark) Init(ctx sparkv1.InitContext) error {
	return nil
}

func (s Spark) Stop() {

}

func (s Spark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	return b.NewChain("chain-1").
		Stage("stage-1", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			var name string
			if err := ctx.Input("name").Bind(&name); err != nil {
				return nil, sparkv1.NewStageError(err)
			}
			return fmt.Sprintf("Hello, %s", name), nil
		}).
		Complete(func(ctx sparkv1.CompleteContext) sparkv1.StageError {
			var (
				message string
				err     error
			)

			// fetch the output of stage-1
			if err = ctx.StageResult("stage-1").Bind(&message); err != nil {
				return sparkv1.NewStageError(err)
			}

			// write the output of the spark
			if err = ctx.Output(
				sparkv1.NewVar("message", "", message),
			); err != nil {
				return sparkv1.NewStageError(err)
			}

			return nil
		})
}

// NewSpark creates a Spark
// shutdown is a reference to the context cancel function, it can be used to gracefully stop the worker if needed
func NewSpark(shutdown context.CancelFunc) sparkv1.Spark {
	return &Spark{shutdown: shutdown}
}
