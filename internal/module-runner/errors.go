package module_runner

import "errors"

var (
	ErrStageResultNotFound            = errors.New("stage result not found")
	ErrChainDoesNotHaveACompleteStage = errors.New("no complete stage found in spark chain")
)
