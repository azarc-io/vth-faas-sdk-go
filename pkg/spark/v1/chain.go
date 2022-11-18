package spark_v1

/************************************************************************/
// TYPES
/************************************************************************/

// chain represents the entire chain
// the rootNode is the main entry point of the entire chain
// it holds its own children as a tree below the rootNode
type chain struct {
	rootNode    *node
	stagesMap   map[string]*stage
	completeMap map[string]*completeStage
}

/************************************************************************/
// HELPERS
/************************************************************************/

func (c *chain) getNodeToResume(lastActiveStage *LastActiveStage) (*node, error) {
	if lastActiveStage == nil {
		return c.rootNode, nil
	}
	if s, ok := c.stagesMap[lastActiveStage.Name]; ok {
		return s.node, nil
	}
	if s, ok := c.completeMap[lastActiveStage.Name]; ok {
		return s.node, nil
	}
	return nil, newErrStageNotFoundInNodeChain(lastActiveStage.Name)
}
