package job_test

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/job"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"testing"
)

func TestBuilder(t *testing.T) {
	_, err := job.NewChain(job.NodeBuilder().
		Stage("stage10", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
		Stage("stage2", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
		Stage("stage3", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
		Complete("stage4", func(context api.CompleteContext) api.StageError { return nil }).
		Canceled(
			job.NodeBuilder().
				Stage("stage5", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
				Canceled(
					job.NodeBuilder().
						Stage("stage6", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
						Compensate(
							job.NodeBuilder().Stage("stage7", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
								Canceled(
									job.NodeBuilder().Stage("stage8", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
										Compensate(
											job.NodeBuilder().Build(), //Stage("stage9", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
										).Build(),
								).Build(),
						).Build()).
				Build()).
		Compensate(job.NodeBuilder().
			Stage("stage10", func(context api.StageContext) (any, api.StageError) { return nil, nil }).
			Compensate(
				job.NodeBuilder().Complete("stage8", func(context api.CompleteContext) api.StageError { return nil }).Build(),
			).Build()). //.Stage("stage11", func(context api.StageContext) (any, api.StageError) { return nil, nil })
		Build()).
		Build()
	print(err.Error())
}
