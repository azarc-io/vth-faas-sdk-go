package job

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
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
	initFunction func()
}

type node struct {
	stages     []*stage
	complete   *stage
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
	cb   api.StageDefinitionFn
	so   []api.StageOption
}

func (s stage) ApplyStageOptionsParams(ctx api.JobContext, stageName string) api.StageError {
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
	sph       api.StageProgressHandler
	vh        api.VariableHandler
	ctx       api.JobContext
}

func (s stageOptionParams) StageName() string {
	return s.stageName
}

func (s stageOptionParams) StageProgressHandler() api.StageProgressHandler {
	return s.sph
}

func (s stageOptionParams) VariableHandler() api.VariableHandler {
	return s.vh
}

func (s stageOptionParams) Context() api.Context {
	return s.ctx
}

func newStageOptionParams(ctx api.JobContext, stageName string) api.StageOptionParams {
	return stageOptionParams{
		stageName: stageName,
		sph:       ctx.StageProgressHandler(),
		vh:        ctx.VariableHandler(),
		ctx:       ctx,
	}
}

func WithStageStatus(stageName string, status sdk_v1.StageStatus) api.StageOption {
	return func(sop api.StageOptionParams) api.StageError {
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
