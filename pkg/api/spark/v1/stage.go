package sdk_v1

/************************************************************************/
// TYPES
/************************************************************************/

type stage struct {
	node *node
	name string
	so   []StageOption
	cb   StageDefinitionFn
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (s *stage) ApplyConditionalExecutionOptions(ctx SparkContext, stageName string) StageError {
	params := newStageOptionParams(ctx, stageName)
	for _, stageOptions := range s.so {
		if err := stageOptions(params); err != nil {
			return err
		}
	}
	return nil
}

/************************************************************************/
// CONTEXT
/************************************************************************/

type stageContext struct {
	Context
	jobContext SparkContext
}

func NewCompleteContext(ctx SparkContext) CompleteContext {
	return stageContext{jobContext: ctx, Context: ctx}
}

func NewStageContext(ctx SparkContext) StageContext {
	return stageContext{jobContext: ctx, Context: ctx}
}

func (sc stageContext) Inputs(names ...string) *Inputs {
	return sc.jobContext.IOHandler().Inputs(sc.JobKey(), names...)
}

func (sc stageContext) Input(name string) *Input {
	return sc.jobContext.IOHandler().Input(sc.JobKey(), name)
}

func (sc stageContext) StageResult(name string) *Result {
	return sc.jobContext.StageProgressHandler().GetResult(sc.JobKey(), name)
}

func (sc stageContext) Output(variables ...*Var) error {
	return sc.jobContext.IOHandler().Output(sc.JobKey(), variables...)
}

func (sc stageContext) Log() Logger {
	return sc.jobContext.Log()
}
