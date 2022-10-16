package job

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
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
}
