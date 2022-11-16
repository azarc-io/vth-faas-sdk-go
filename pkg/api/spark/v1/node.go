package sdk_v1

import (
	"fmt"
)

/************************************************************************/
// TYPES
/************************************************************************/

type nodeType string

var (
	rootNodeType       = nodeType("root")
	compensateNodeType = nodeType("compensate")
	cancelNodeType     = nodeType("canceled")
)

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

/************************************************************************/
// HELPERS
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
