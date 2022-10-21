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

func (s stage) ApplyStageOptionsParams(ctx sdk_v1.SparkContext, stageName string) sdk_v1.StageError {
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
	vh        sdk_v1.VariableHandler
	ctx       sdk_v1.SparkContext
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) StageProgressHandler() sdk_v1.StageProgressHandler {
	return s.sph
}

func (s stageOptionParams) VariableHandler() sdk_v1.VariableHandler {
	return s.vh
}

func (s stageOptionParams) Context() sdk_v1.Context {
	return s.ctx
}

func newStageOptionParams(ctx sdk_v1.SparkContext, stageName string) sdk_v1.StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		sph:       ctx.StageProgressHandler(),
		vh:        ctx.VariableHandler(),
		ctx:       ctx,
	}
}

func WithStageStatus(stageName string, status sdk_v1.StageStatus) sdk_v1.StageOption {
	return func(sop sdk_v1.StageOptionParams) sdk_v1.StageError {
		stageStatus, err := sop.StageProgressHandler().Get(sop.Context().JobKey(), stageName)
		if err != nil {
			return sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_Failed))
		}
		if *stageStatus != status {
			return sdk_errors.NewStageError(fmt.Errorf("conditional stage execution skipped this stage"), sdk_errors.WithErrorType(sdk_v1.ErrorType_Skip))
		}
		return nil
	}
}
