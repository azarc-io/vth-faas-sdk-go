package sdk_v1

import (
	"errors"
	"fmt"
)

/************************************************************************/
// ERRORS
/************************************************************************/

var ErrStageNotFoundInNodeChain = errors.New("stage not found in the node chain")

func newErrStageNotFoundInNodeChain(stage string) error {
	return fmt.Errorf("%w: %s", ErrStageNotFoundInNodeChain, stage)
}

var ErrConditionalStageSkipped = errors.New("conditional stage execution")

func newErrConditionalStageSkipped(stageName string) error {
	return fmt.Errorf("%w: stage '%s' skipped", ErrConditionalStageSkipped, stageName)
}
