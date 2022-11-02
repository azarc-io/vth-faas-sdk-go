// Copyright 2020-2022 Azarc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package spark

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
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
	if err := c.validate([]*validateFn{stageNamesMustNoBeEmpty, atLeastOneStagePerNodeValidator, uniqueStageNamesValidator()}, newChain.rootNode); err != nil {
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

func (sb *nodeBuilder) Stage(name string, stageDefinitionFn sdk_v1.StageDefinitionFn, options ...sdk_v1.StageOption) *nodeBuilder {
	s := &stage{
		node: sb.node,
		name: name,
		cb:   stageDefinitionFn,
		so:   options,
	}
	sb.node.stages = append(sb.node.stages, s)
	return sb
}

func (sb *nodeBuilder) Complete(name string, completeDefinitionFn sdk_v1.CompleteDefinitionFn, options ...sdk_v1.StageOption) *nodeBuilder {
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
