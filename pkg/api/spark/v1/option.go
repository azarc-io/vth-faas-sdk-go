package sdk_v1

/************************************************************************/
// OPTION TYPES
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

/************************************************************************/
// STAGE OPTIONS
/************************************************************************/

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
