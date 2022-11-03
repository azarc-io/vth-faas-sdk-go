package spark

import (
	"fmt"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

type nodeType string

var (
	nodeTypeRoot       = nodeType("root")
	compensateNodeType = nodeType("compensate")
	cancelNodeType     = nodeType("canceled")
)

type Chain struct {
	rootNode     *node
	stagesMap    map[string]*stage
	completeMap  map[string]*completeStage
	initFunction func()
}

func (c *Chain) getNodeToResume(lastActiveStage *sdk_v1.LastActiveStage) (*node, error) {
	if lastActiveStage == nil {
		return c.rootNode, nil
	}
	if s, ok := c.stagesMap[lastActiveStage.Name]; ok {
		return s.node, nil
	}
	if s, ok := c.completeMap[lastActiveStage.Name]; ok {
		return s.node, nil
	}
	return nil, fmt.Errorf("stage '%s' not found in the node chain", lastActiveStage.Name)
}

type node struct {
	stages     []*stage
	complete   *completeStage
	cancel     *node
	compensate *node
	nodeType   nodeType
	breadcrumb string
}

func (n *node) appendBreadcrumb(nodeType nodeType, breadcrumb ...string) {
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
	node *node
	name string
	so   []sdk_v1.StageOption
	cb   sdk_v1.StageDefinitionFn
}

type completeStage struct {
	node *node
	name string
	so   []sdk_v1.StageOption
	cb   sdk_v1.CompleteDefinitionFn
}

func (s stage) ApplyConditionalExecutionOptions(ctx sdk_v1.SparkContext, stageName string) sdk_v1.StageError {
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
	sph       sdk_v1.StageProgressHandler
	vh        sdk_v1.IOHandler
	ctx       sdk_v1.SparkContext
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) StageProgressHandler() sdk_v1.StageProgressHandler {
	return s.sph
}

func (s stageOptionParams) IOHandler() sdk_v1.IOHandler {
	return s.vh
}

func (s stageOptionParams) Context() sdk_v1.Context {
	return s.ctx
}

func newStageOptionParams(ctx sdk_v1.SparkContext, stageName string) sdk_v1.StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		sph:       ctx.StageProgressHandler(),
		vh:        ctx.IOHandler(),
		ctx:       ctx,
	}
}

func WithStageStatus(stageName string, status sdk_v1.StageStatus) sdk_v1.StageOption {
	return func(sop sdk_v1.StageOptionParams) sdk_v1.StageError {
		stageStatus, err := sop.StageProgressHandler().Get(sop.Context().JobKey(), stageName)
		if err != nil {
			return sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED))
		}
		if *stageStatus != status {
			return sdk_errors.NewStageError(fmt.Errorf("conditional stage execution: stage '%s' skipped", stageName), sdk_errors.WithErrorType(sdk_v1.ErrorType_ERROR_TYPE_SKIP))
		}
		return nil
	}
}
