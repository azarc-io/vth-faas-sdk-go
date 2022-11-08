package spark

import (
	"errors"
	"fmt"

	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"

	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

type nodeType string

var (
	nodeTypeRoot       = nodeType("root")
	compensateNodeType = nodeType("compensate")
	cancelNodeType     = nodeType("canceled")
)

type Chain struct {
	rootNode    *Node
	stagesMap   map[string]*stage
	completeMap map[string]*completeStage
}

var ErrStageNotFoundInNodeChain = errors.New("stage not found in the node chain")

func newErrStageNotFoundInNodeChain(stage string) error {
	return fmt.Errorf("%w: %s", ErrStageNotFoundInNodeChain, stage)
}

func (c *Chain) getNodeToResume(lastActiveStage *v1.LastActiveStage) (*Node, error) {
	if lastActiveStage == nil {
		return c.rootNode, nil
	}
	if s, ok := c.stagesMap[lastActiveStage.Name]; ok {
		return s.node, nil
	}
	if s, ok := c.completeMap[lastActiveStage.Name]; ok {
		return s.node, nil
	}
	return nil, newErrStageNotFoundInNodeChain(lastActiveStage.Name)
}

type Node struct {
	stages     []*stage
	complete   *completeStage
	cancel     *Node
	compensate *Node
	nodeType   nodeType
	breadcrumb string
}

func (n *Node) appendBreadcrumb(nodeType nodeType, breadcrumb ...string) {
	if n == nil {
		return
	}
	n.nodeType = nodeType
	if breadcrumb == nil {
		n.breadcrumb = string(nodeType)
		return
	}
	n.breadcrumb = fmt.Sprintf("%s > %s", breadcrumb[0], nodeType)
}

type stage struct {
	node *Node
	name string
	so   []v1.StageOption
	cb   v1.StageDefinitionFn
}

type completeStage struct {
	node *Node
	name string
	so   []v1.StageOption
	cb   v1.CompleteDefinitionFn
}

func (s stage) ApplyConditionalExecutionOptions(ctx v1.SparkContext, stageName string) v1.StageError {
	params := newStageOptionParams(ctx, stageName)
	for _, stageOptions := range s.so {
		if err := stageOptions(params); err != nil {
			return err
		}
	}
	return nil
}

type stageOptionParams struct {
	stageName string
	sph       v1.StageProgressHandler
	vh        v1.IOHandler
	ctx       v1.SparkContext
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) StageProgressHandler() v1.StageProgressHandler {
	return s.sph
}

func (s stageOptionParams) IOHandler() v1.IOHandler {
	return s.vh
}

func (s stageOptionParams) Context() v1.Context {
	return s.ctx
}

func newStageOptionParams(ctx v1.SparkContext, stageName string) v1.StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		sph:       ctx.StageProgressHandler(),
		vh:        ctx.IOHandler(),
		ctx:       ctx,
	}
}

var ErrConditionalStageSkipped = errors.New("conditional stage execution")

func newErrConditionalStageSkipped(stageName string) error {
	return fmt.Errorf("%w: stage '%s' skipped", ErrConditionalStageSkipped, stageName)
}

func WithStageStatus(stageName string, status v1.StageStatus) v1.StageOption {
	return func(sop v1.StageOptionParams) v1.StageError {
		stageStatus, err := sop.StageProgressHandler().Get(sop.Context().JobKey(), stageName)
		if err != nil {
			return sdk_errors.NewStageError(err, sdk_errors.WithErrorType(v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED))
		}
		if *stageStatus != status {
			return sdk_errors.NewStageError(newErrConditionalStageSkipped(stageName), sdk_errors.WithErrorType(v1.ErrorType_ERROR_TYPE_SKIP))
		}
		return nil
	}
}
