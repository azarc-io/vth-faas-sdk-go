package spark_v1

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/spark/v1"
)

var (
	rootNodeType       = nodeType("root")
	compensateNodeType = nodeType("compensate")
	cancelNodeType     = nodeType("canceled")
)

/************************************************************************/
// TYPES
/************************************************************************/

// nodeType describes the type of chain node
type nodeType string

// chain represents the entire chain
// the rootNode is the main entry point of the entire chain
// it holds its own children as a tree below the rootNode
type chain struct {
	rootNode    *node
	stagesMap   map[string]*stage
	completeMap map[string]*completeStage
}

// node wraps all the stages of a single chain
// these are represented as one or more stages but only one of each
// - cancellation
// - compensation
// - completion (finalizer)
type node struct {
	stages     []*stage
	complete   *completeStage
	cancel     *node
	compensate *node
	nodeType   nodeType
	breadcrumb string
	name       string
}

type completeStage struct {
	node *node
	name string
	so   []StageOption
	cb   CompleteDefinitionFn
}

type stage struct {
	node *node
	name string
	so   []StageOption
	cb   StageDefinitionFn
}

/************************************************************************/
// CHAIN HELPERS
/************************************************************************/

func (c *chain) getNodeToResume(lastActiveStage *sparkv1.LastActiveStage) (*node, error) {
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

/************************************************************************/
// CHAIN NODE HELPERS
/************************************************************************/

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

/************************************************************************/
// CHAIN NODE STAGE HELPERS
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
// CHAIN NODE ACCESSORS
/************************************************************************/

func (n *node) ChainName() string {
	return n.name
}

func (n *node) CompletionName() string {
	return n.complete.name
}

func (n *node) CountOfStages() int {
	return len(n.stages)
}

func (n *node) HasCompletionStage() bool {
	return n.complete != nil
}

func (n *node) HasCompensationStage() bool {
	return n.compensate != nil
}

func (n *node) HasCancellationStage() bool {
	return n.cancel != nil
}

func (n *node) IsRoot() bool {
	return n.nodeType == rootNodeType
}

func (n *node) IsCompensate() bool {
	return n.nodeType == compensateNodeType
}

func (n *node) IsCancel() bool {
	return n.nodeType == cancelNodeType
}
