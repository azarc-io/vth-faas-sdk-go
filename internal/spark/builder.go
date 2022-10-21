package spark

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
)

type chainBuilder struct {
	rootNode *node
}

func NewChain(node *node) *chainBuilder {
	c := &chainBuilder{}
	c.rootNode = node
	c.rootNode.appendBreadcrumb(nodeTypeRoot)
	return c
}

func (c *chainBuilder) Build() (*Chain, error) {
	newChain := &Chain{
		rootNode:    c.rootNode,
		stagesMap:   map[string]*stage{},
		completeMap: map[string]*completeStage{},
	}
	c.createResumeOnRetryStagesMap(newChain)
	addBreadcrumb(newChain.rootNode)
	if err := c.validate([]*validateFn{atLeastOneStagePerNodeValidator, uniqueStageNamesValidator()}, newChain.rootNode); err != nil {
		return nil, err
	}
	return newChain, nil
}

func (c *chainBuilder) createResumeOnRetryStagesMap(newChain *Chain) {
	stages, completeStages := c.getStages([]*stage{}, []*completeStage{}, newChain.rootNode)
	for _, stg := range stages {
		newChain.stagesMap[stg.name] = stg
	}
	for _, cStg := range completeStages {
		newChain.completeMap[cStg.name] = cStg
	}
}

func (c *chainBuilder) getStages(stages []*stage, completeStages []*completeStage, nodes ...*node) ([]*stage, []*completeStage) {
	var nextNodes []*node
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

type nodeBuilder struct {
	node *node
}

func NewNode() *nodeBuilder {
	return &nodeBuilder{node: &node{}}
}

func (sb *nodeBuilder) Stage(name string, stageDefinitionFn api.StageDefinitionFn, options ...api.StageOption) *nodeBuilder {
	s := &stage{
		node: sb.node,
		name: name,
		cb:   stageDefinitionFn,
		so:   options,
	}
	sb.node.stages = append(sb.node.stages, s)
	return sb
}

func (sb *nodeBuilder) Complete(name string, completeDefinitionFn api.CompleteDefinitionFn, options ...api.StageOption) *nodeBuilder {
	sb.node.complete = &completeStage{
		node: sb.node,
		name: name,
		cb:   completeDefinitionFn,
		so:   options,
	}
	return sb
}

func (sb *nodeBuilder) Canceled(newNode *node) *nodeBuilder {
	sb.node.cancel = newNode
	return sb
}

func (sb *nodeBuilder) Compensate(newNode *node) *nodeBuilder {
	sb.node.compensate = newNode
	return sb
}

func (sb *nodeBuilder) Build() *node {
	return sb.node
}

func addBreadcrumb(nodes ...*node) {
	var nextNodes []*node
	for _, n := range nodes {
		n.cancel.appendBreadcrumb(cancelNodeType, n.breadcrumb)
		n.compensate.appendBreadcrumb(compensateNodeType, n.breadcrumb)
		nextNodes = appendIfNotNil(nextNodes, n.compensate, n.cancel)
	}
	if len(nextNodes) > 0 {
		addBreadcrumb(nextNodes...)
		return
	}
	return
}
