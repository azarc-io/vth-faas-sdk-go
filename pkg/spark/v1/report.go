package spark_v1

import "fmt"

type (
	ChainReport struct {
		Errors   []error
		StageMap map[string]ChainReportStage
		NodeMap  map[string]ChainReportNode
	}

	ChainReportNode struct {
		Name          string
		CanCompensate bool
		CanCancel     bool
		node          *node
	}

	ChainReportStage struct {
		Name  string
		Crumb string
	}
)

// generateReportForChain generates a map of what the chain looks like
func generateReportForChain(n *chain) ChainReport {
	r := ChainReport{
		StageMap: map[string]ChainReportStage{},
		NodeMap:  map[string]ChainReportNode{},
	}

	generateReportForChainRecursively(&r, n.rootNode)

	return r
}

func generateReportForChainRecursively(r *ChainReport, n *node) {
	// must not have an empty name
	if n.name == "" {
		r.Errors = append(r.Errors, fmt.Errorf("chain name can not be empty [at]: %s", n.breadcrumb))
	}

	// chain node name must be unique
	if _, ok := r.NodeMap[n.name]; ok {
		r.Errors = append(r.Errors, fmt.Errorf("duplicate chain names are not permitted [name]: %s [at]: %s",
			n.name, n.breadcrumb))
	} else {
		r.NodeMap[n.name] = ChainReportNode{
			Name:          n.name,
			CanCompensate: n.HasCompensationStage(),
			CanCancel:     n.HasCancellationStage(),
			node:          n,
		}
	}

	// first flat map all the stages and capture validation errors
	for _, s := range n.stages {
		if s.name == "" {
			r.Errors = append(r.Errors, fmt.Errorf("stage name can not be empty [at]: %s", n.breadcrumb))
			continue
		}

		if _, ok := r.StageMap[s.name]; ok {
			r.Errors = append(r.Errors, fmt.Errorf("duplicate stage names are not permitted [chain]: %s [at]: %s",
				s.name, n.breadcrumb))
		} else {
			r.StageMap[s.name] = ChainReportStage{
				Name:  s.name,
				Crumb: n.breadcrumb,
			}
		}
	}

	if n.HasCompensationStage() {
		generateReportForChainRecursively(r, n.compensate)
	}

	if n.HasCancellationStage() {
		generateReportForChainRecursively(r, n.cancel)
	}
}
