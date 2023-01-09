package sparkv1

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
		node          *Node
	}

	ChainReportStage struct {
		Name  string
		Crumb string
	}
)

// generateReportForChain generates a map of what the SparkChain looks like
func generateReportForChain(n *SparkChain) ChainReport {
	r := ChainReport{
		StageMap: map[string]ChainReportStage{},
		NodeMap:  map[string]ChainReportNode{},
	}

	generateReportForChainRecursively(&r, n.RootNode)

	return r
}

func generateReportForChainRecursively(r *ChainReport, n *Node) {
	// must not have an empty Name
	if n.Name == "" {
		r.Errors = append(r.Errors, fmt.Errorf("SparkChain Name can not be empty [at]: %s", n.breadcrumb))
	}

	// SparkChain Node Name must be unique
	if _, ok := r.NodeMap[n.Name]; ok {
		r.Errors = append(r.Errors, fmt.Errorf("duplicate SparkChain names are not permitted [Name]: %s [at]: %s",
			n.Name, n.breadcrumb))
	} else {
		r.NodeMap[n.Name] = ChainReportNode{
			Name:          n.Name,
			CanCompensate: n.HasCompensationStage(),
			CanCancel:     n.HasCancellationStage(),
			node:          n,
		}
	}

	// first flat map all the Stages and capture validation errors
	for _, s := range n.Stages {
		if s.Name == "" {
			r.Errors = append(r.Errors, fmt.Errorf("Stage Name can not be empty [at]: %s", n.breadcrumb))
			continue
		}

		if _, ok := r.StageMap[s.Name]; ok {
			r.Errors = append(r.Errors, fmt.Errorf("duplicate Stage names are not permitted [SparkChain]: %s [at]: %s",
				s.Name, n.breadcrumb))
		} else {
			r.StageMap[s.Name] = ChainReportStage{
				Name:  s.Name,
				Crumb: n.breadcrumb,
			}
		}
	}

	if n.HasCompensationStage() {
		generateReportForChainRecursively(r, n.Compensate)
	}

	if n.HasCancellationStage() {
		generateReportForChainRecursively(r, n.Cancel)
	}
}
