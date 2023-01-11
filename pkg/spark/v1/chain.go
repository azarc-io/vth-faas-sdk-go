package sparkv1

import (
	"fmt"
)

var (
	rootNodeType       = nodeType("root")
	compensateNodeType = nodeType("Compensate")
	cancelNodeType     = nodeType("canceled")
)

/************************************************************************/
// TYPES
/************************************************************************/

// nodeType describes the type of SparkChain Node
type nodeType string

// SparkChain represents the entire SparkChain
// the RootNode is the main entry point of the entire SparkChain
// it holds its own children as a tree below the RootNode
type SparkChain struct {
	RootNode    *Node
	stagesMap   map[string]*Stage
	completeMap map[string]*CompleteStage
}

func (sc *SparkChain) GetStageFunc(name string) StageDefinitionFn {
	return sc.stagesMap[name].cb
}

func (sc *SparkChain) GetStageCompleteFunc(name string) CompleteDefinitionFn {
	return sc.completeMap[name].cb
}

// Node wraps all the Stages of a single SparkChain
// these are represented as one or more Stages but only one of each
// - cancellation
// - compensation
// - completion (finalizer)
type Node struct {
	Stages     []*Stage
	Complete   *CompleteStage
	Cancel     *Node
	Compensate *Node
	Name       string
	nodeType   nodeType
	breadcrumb string
}

type CompleteStage struct {
	Node *Node
	Name string
	so   []StageOption
	cb   CompleteDefinitionFn
}

type Stage struct {
	Node *Node
	Name string
	so   []StageOption
	cb   StageDefinitionFn
}

/************************************************************************/
// CHAIN NODE HELPERS
/************************************************************************/

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

/************************************************************************/
// CHAIN NODE STAGE HELPERS
/************************************************************************/

func (s *Stage) ApplyConditionalExecutionOptions(ctx Context, stageName string) StageError {
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

func (n *Node) ChainName() string {
	return n.Name
}

func (n *Node) CompletionName() string {
	return n.Complete.Name
}

func (n *Node) CountOfStages() int {
	return len(n.Stages)
}

func (n *Node) HasCompletionStage() bool {
	return n.Complete != nil
}

func (n *Node) HasCompensationStage() bool {
	return n.Compensate != nil
}

func (n *Node) HasCancellationStage() bool {
	return n.Cancel != nil
}

func (n *Node) IsRoot() bool {
	return n.nodeType == rootNodeType
}

func (n *Node) IsCompensate() bool {
	return n.nodeType == compensateNodeType
}

func (n *Node) IsCancel() bool {
	return n.nodeType == cancelNodeType
}
