package spark

import (
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

type ChainBuilder struct {
	rootNode *Node
}

func NewChain(node *Node) *ChainBuilder {
	c := &ChainBuilder{}
	c.rootNode = node
	c.rootNode.appendBreadcrumb(nodeTypeRoot)
	return c
}

func (c *ChainBuilder) Build() (*Chain, error) {
	newChain := &Chain{
		rootNode:    c.rootNode,
		stagesMap:   map[string]*stage{},
		completeMap: map[string]*completeStage{},
	}
	c.createResumeOnRetryStagesMap(newChain)
	addBreadcrumb(newChain.rootNode)
	if err := c.validate([]*validateFn{stageNamesMustNoBeEmpty, atLeastOneStagePerNodeValidator, uniqueStageNamesValidator()}, newChain.rootNode); err != nil {
		return nil, err
	}
	return newChain, nil
}

func (c *ChainBuilder) createResumeOnRetryStagesMap(newChain *Chain) {
	stages, completeStages := c.getStages([]*stage{}, []*completeStage{}, newChain.rootNode)
	for _, stg := range stages {
		newChain.stagesMap[stg.name] = stg
	}
	for _, cStg := range completeStages {
		newChain.completeMap[cStg.name] = cStg
	}
}

func (c *ChainBuilder) getStages(stages []*stage, completeStages []*completeStage, nodes ...*Node) ([]*stage, []*completeStage) {
	var nextNodes []*Node
	for _, n := range nodes {
		completeStages = appendIfNotNil(completeStages, n.complete)
		stages = appendIfNotNil(stages, n.stages...)
		nextNodes = appendIfNotNil(nextNodes, n.compensate, n.cancel)
	}
	if len(nextNodes) > 0 {
		return c.getStages(stages, completeStages, nextNodes...)
	}
	return stages, completeStages
}

type NodeBuilder struct {
	node *Node
}

func NewNode() *NodeBuilder {
	return &NodeBuilder{node: &Node{}}
}

func (sb *NodeBuilder) Stage(name string, stageDefinitionFn sdk_v1.StageDefinitionFn, options ...sdk_v1.StageOption) *NodeBuilder {
	s := &stage{
		node: sb.node,
		name: name,
		cb:   stageDefinitionFn,
		so:   options,
	}
	sb.node.stages = append(sb.node.stages, s)
	return sb
}

func (sb *NodeBuilder) Complete(name string, completeDefinitionFn sdk_v1.CompleteDefinitionFn, options ...sdk_v1.StageOption) *NodeBuilder {
	sb.node.complete = &completeStage{
		node: sb.node,
		name: name,
		cb:   completeDefinitionFn,
		so:   options,
	}
	return sb
}

func (sb *NodeBuilder) Cancelled(newNode *Node) *NodeBuilder {
	sb.node.cancel = newNode
	return sb
}

func (sb *NodeBuilder) Compensate(newNode *Node) *NodeBuilder {
	sb.node.compensate = newNode
	return sb
}

func (sb *NodeBuilder) Build() *Node {
	return sb.node
}

func addBreadcrumb(nodes ...*Node) {
	var nextNodes []*Node
	for _, n := range nodes {
		n.cancel.appendBreadcrumb(cancelNodeType, n.breadcrumb)
		n.compensate.appendBreadcrumb(compensateNodeType, n.breadcrumb)
		nextNodes = appendIfNotNil(nextNodes, n.compensate, n.cancel)
	}
	if len(nextNodes) > 0 {
		addBreadcrumb(nextNodes...)
	}
}
