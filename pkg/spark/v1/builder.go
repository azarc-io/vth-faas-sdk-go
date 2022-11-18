package spark_v1

import "fmt"

/************************************************************************/
// TYPES
/************************************************************************/

type chainBuilder struct {
	rootNode *chainNode
	current  *chainNode
	prev     *chainNode
}

// chainNode wraps the chain and the node for easier access to both
type chainNode struct {
	// BuilderChain This is the builder that was used to build this chain, this enables
	// the builder to return the node as a child of a stage or as the root of the chain
	BuilderChain
	// node holds the node information until the chain is built
	node *node
}

/************************************************************************/
// BUILDER API
/************************************************************************/

// Stage adds a stage to the current chain node, this could be at any depth of the chain
func (c *chainBuilder) Stage(name string, stageDefinitionFn StageDefinitionFn, options ...StageOption) ChainStageAny {
	s := &stage{
		node: c.current.node,
		name: name,
		cb:   stageDefinitionFn,
		so:   options,
	}
	c.current.node.stages = append(c.current.node.stages, s)
	return c
}

// Compensate registers a chain node at depth-1 in the chain, compensation is always on the parent
// so this function looks at the previous node in the chain which is always the parent
func (c *chainBuilder) Compensate(newNode Chain) ChainCancelledOrComplete {
	n := newNode.build() // this causes the chain to move from depth to depth-1
	n.nodeType = compensateNodeType

	n.appendBreadcrumb(compensateNodeType)

	c.current.node.compensate = n

	return c
}

// Cancelled registers a chain node at depth-1 in the chain, compensation is always on the parent
// so this function looks at the previous node in the chain which is always the parent
func (c *chainBuilder) Cancelled(newNode Chain) ChainComplete {
	n := newNode.build() // this causes the chain to move from depth to depth-1
	n.nodeType = cancelNodeType

	n.appendBreadcrumb(cancelNodeType)

	c.current.node.cancel = n

	return c
}

// Complete returns a finalizer that can be used to build the node chain
func (c *chainBuilder) Complete(completeDefinitionFn CompleteDefinitionFn, options ...StageOption) Chain {
	name := fmt.Sprintf("%s_complete", c.current.node.name)
	c.current.node.complete = &completeStage{
		node: c.rootNode.node,
		name: name,
		cb:   completeDefinitionFn,
		so:   options,
	}

	return c
}

/************************************************************************/
// HELPERS
/************************************************************************/

// createResumeOnRetryStagesMap maps stages that can be retried
func (c *chainBuilder) createResumeOnRetryStagesMap(newChain *chain) {
	stages, completeStages := c.getStages([]*stage{}, []*completeStage{}, newChain.rootNode)
	for _, stg := range stages {
		newChain.stagesMap[stg.name] = stg
	}
	for _, cStg := range completeStages {
		newChain.completeMap[cStg.name] = cStg
	}
}

// getStages returns all stages + completion stages
func (c *chainBuilder) getStages(stages []*stage, completeStages []*completeStage, nodes ...*node) ([]*stage, []*completeStage) {
	var nextNodes []*node
	for _, n := range nodes {
		completeStages = appendIfNotNil(completeStages, n.complete)
		stages = appendIfNotNil(stages, n.stages...)
		nextNodes = appendIfNotNil(nextNodes, n.compensate, n.cancel)
	}
	return stages, completeStages
}

/************************************************************************/
// FINALIZERS
/************************************************************************/

// buildChain creates a chain that can be executed
// - Maps the chain
// - Decorates it with breadcrumbs
func (c *chainBuilder) buildChain() *chain {
	newChain := &chain{
		rootNode:    c.rootNode.node,
		stagesMap:   map[string]*stage{},
		completeMap: map[string]*completeStage{},
	}
	c.createResumeOnRetryStagesMap(newChain)
	addBreadcrumb(newChain.rootNode)
	return newChain
}

// Build iterates over the entire chain and performs the following operations
// - Decorate the chain node with breadcrumbs
// - Move the pointer back up the chain (depth-1)
func (c *chainBuilder) build() *node {
	ret := c.current

	addBreadcrumb(c.rootNode.node)

	// this is a finalizer so switch current back to previous node in order to push the pointer back up the tree
	if c.prev != nil {
		c.current = c.prev // move the pointer up one
		c.prev = nil       // clear previous because it should only be set when the depth of the chain changes
	}

	return ret.node
}

/************************************************************************/
// FACTORIES
/************************************************************************/

// NewChain creates a new chain with the following rules
// - if this is the first chain the builder sees then it's marked as the root chain
// - if a root chain already exists then the new chain is returned but not stored because it will be a
// compensation or a cancellation chain
func (c *chainBuilder) NewChain(name string) BuilderChain {
	n := &chainNode{
		BuilderChain: c,
		node: &node{
			name: name,
		},
	}

	// The first time NewChain is called will store the chainNode as the root node
	// any future chains that are created are children of the root node as such
	// simply return the chainNode
	if c.rootNode == nil {
		n.node.nodeType = rootNodeType
		c.rootNode = n
		n.node.appendBreadcrumb(rootNodeType)
	}

	// this holds the next and previous chain that is being built, it could be the root or a nested chain
	c.prev = c.current
	c.current = n

	return n
}

// newBuilder main entry point to the builder
func newBuilder() Builder {
	return &chainBuilder{}
}
