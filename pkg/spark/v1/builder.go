package sparkv1

import "fmt"

/************************************************************************/
// TYPES
/************************************************************************/

type chainBuilder struct {
	rootNode *chainNode
	current  *chainNode
	prev     *chainNode
}

// chainNode wraps the SparkChain and the Node for easier access to both
type chainNode struct {
	// BuilderChain This is the builder that was used to build this SparkChain, this enables
	// the builder to return the node as a child of a Stage or as the root of the SparkChain
	BuilderChain
	// node holds the node information until the SparkChain is built
	node *Node
}

/************************************************************************/
// BUILDER API
/************************************************************************/

// Stage adds a Stage to the current SparkChain Node, this could be at any depth of the SparkChain
func (c *chainBuilder) Stage(name string, stageDefinitionFn StageDefinitionFn, options ...StageOption) ChainStageAny {
	s := &Stage{
		Node: c.current.node,
		Name: name,
		cb:   stageDefinitionFn,
		so:   options,
	}
	c.current.node.Stages = append(c.current.node.Stages, s)
	return c
}

// Compensate registers a SparkChain Node at depth-1 in the SparkChain, compensation is always on the parent
// so this function looks at the previous Node in the SparkChain which is always the parent
func (c *chainBuilder) Compensate(newNode Chain) ChainCancelledOrComplete {
	n := newNode.build() // this causes the SparkChain to move from depth to depth-1
	n.nodeType = compensateNodeType

	n.appendBreadcrumb(compensateNodeType)

	c.current.node.Compensate = n

	return c
}

// Cancelled registers a SparkChain Node at depth-1 in the SparkChain, compensation is always on the parent
// so this function looks at the previous Node in the SparkChain which is always the parent
func (c *chainBuilder) Cancelled(newNode Chain) ChainComplete {
	n := newNode.build() // this causes the SparkChain to move from depth to depth-1
	n.nodeType = cancelNodeType

	n.appendBreadcrumb(cancelNodeType)

	c.current.node.Cancel = n

	return c
}

// Complete returns a finalizer that can be used to build the Node SparkChain
func (c *chainBuilder) Complete(completeDefinitionFn CompleteDefinitionFn, options ...StageOption) Chain {
	name := fmt.Sprintf("%s_complete", c.current.node.Name)
	c.current.node.Complete = &CompleteStage{
		Node: c.rootNode.node,
		Name: name,
		cb:   completeDefinitionFn,
		so:   options,
	}

	return c
}

/************************************************************************/
// HELPERS
/************************************************************************/

// createResumeOnRetryStagesMap maps Stages that can be retried
func (c *chainBuilder) createResumeOnRetryStagesMap(newChain *SparkChain) {
	stages, completeStages := c.getStages([]*Stage{}, []*CompleteStage{}, newChain.RootNode)
	for _, stg := range stages {
		newChain.stagesMap[stg.Name] = stg
	}
	for _, cStg := range completeStages {
		newChain.completeMap[cStg.Name] = cStg
	}
}

// getStages returns all Stages + completion Stages
func (c *chainBuilder) getStages(stages []*Stage, completeStages []*CompleteStage, nodes ...*Node) ([]*Stage, []*CompleteStage) {
	var nextNodes []*Node
	for _, n := range nodes {
		completeStages = appendIfNotNil(completeStages, n.Complete)
		stages = appendIfNotNil(stages, n.Stages...)
		nextNodes = appendIfNotNil(nextNodes, n.Compensate, n.Cancel)
	}
	return stages, completeStages
}

/************************************************************************/
// FINALIZERS
/************************************************************************/

// buildChain creates a SparkChain that can be executed
// - Maps the SparkChain
// - Decorates it with breadcrumbs
func (c *chainBuilder) BuildChain() *SparkChain {
	newChain := &SparkChain{
		RootNode:    c.rootNode.node,
		stagesMap:   map[string]*Stage{},
		completeMap: map[string]*CompleteStage{},
	}
	c.createResumeOnRetryStagesMap(newChain)
	addBreadcrumb(newChain.RootNode)
	return newChain
}

// Build iterates over the entire SparkChain and performs the following operations
// - Decorate the SparkChain Node with breadcrumbs
// - Move the pointer back up the SparkChain (depth-1)
func (c *chainBuilder) build() *Node {
	ret := c.current

	addBreadcrumb(c.rootNode.node)

	// this is a finalizer so switch current back to previous Node in order to push the pointer back up the tree
	if c.prev != nil {
		c.current = c.prev // move the pointer up one
		c.prev = nil       // clear previous because it should only be set when the depth of the SparkChain changes
	}

	return ret.node
}

/************************************************************************/
// FACTORIES
/************************************************************************/

// NewChain creates a new SparkChain with the following rules
// - if this is the first SparkChain the builder sees then it's marked as the root SparkChain
// - if a root SparkChain already exists then the new SparkChain is returned but not stored because it will be a
// compensation or a cancellation SparkChain
func (c *chainBuilder) NewChain(name string) BuilderChain {
	n := &chainNode{
		BuilderChain: c,
		node: &Node{
			Name: name,
		},
	}

	// The first time NewChain is called will store the chainNode as the root Node
	// any future chains that are created are children of the root Node as such
	// simply return the chainNode
	if c.rootNode == nil {
		n.node.nodeType = rootNodeType
		c.rootNode = n
		n.node.appendBreadcrumb(rootNodeType)
	}

	// this holds the next and previous SparkChain that is being built, it could be the root or a nested SparkChain
	c.prev = c.current
	c.current = n

	return n
}

// NewBuilder main entry point to the builder
func NewBuilder() Builder {
	return &chainBuilder{}
}
