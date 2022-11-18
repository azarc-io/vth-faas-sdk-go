package spark_v1

/************************************************************************/
// STAGE OPTIONS
/************************************************************************/

type stageOptionParams struct {
	stageName string
	sph       StageProgressHandler
	vh        IOHandler
	ctx       SparkContext
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) StageProgressHandler() StageProgressHandler {
	return s.sph
}

func (s stageOptionParams) IOHandler() IOHandler {
	return s.vh
}

func (s stageOptionParams) Context() Context {
	return s.ctx
}

func newStageOptionParams(ctx SparkContext, stageName string) StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		sph:       ctx.StageProgressHandler(),
		vh:        ctx.IOHandler(),
		ctx:       ctx,
	}
}

func WithStageStatus(stageName string, status StageStatus) StageOption {
	return func(sop StageOptionParams) StageError {
		stageStatus, err := sop.StageProgressHandler().Get(sop.Context().JobKey(), stageName)
		if err != nil {
			return NewStageError(err, WithErrorType(ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED))
		}
		if *stageStatus != status {
			return NewStageError(newErrConditionalStageSkipped(stageName), WithErrorType(ErrorType_ERROR_TYPE_SKIP))
		}
		return nil
	}
}

/************************************************************************/
// SPARK OPTIONS
/************************************************************************/

type sparkOpts struct {
	variableHandler      IOHandler
	stageProgressHandler StageProgressHandler
	log                  Logger
	delegateStage        DelegateStageDefinitionFn
	delegateComplete     DelegateCompleteDefinitionFn
}

type Option = func(je *sparkOpts) *sparkOpts

func WithIOHandler(vh IOHandler) Option {
	return func(jw *sparkOpts) *sparkOpts {
		jw.variableHandler = vh
		return jw
	}
}

func WithStageProgressHandler(sph StageProgressHandler) Option {
	return func(jw *sparkOpts) *sparkOpts {
		jw.stageProgressHandler = sph
		return jw
	}
}

func WithLog(log Logger) Option {
	return func(jw *sparkOpts) *sparkOpts {
		jw.log = log
		return jw
	}
}

// WithDelegateStage delegates execution of all stages
// TODO support delegating single stage by name
func WithDelegateStage(delegate DelegateStageDefinitionFn) Option {
	return func(jw *sparkOpts) *sparkOpts {
		jw.delegateStage = delegate
		return jw
	}
}

// WithDelegateCompletion delegates execution of all completion stages
// TODO support delegating single completion stage by name
func WithDelegateCompletion(delegate DelegateCompleteDefinitionFn) Option {
	return func(jw *sparkOpts) *sparkOpts {
		jw.delegateComplete = delegate
		return jw
	}
}
