package sparkv1

import (
	"errors"
)

var CompleteSuccess = func(ctx CompleteContext) StageError {
	return nil
}

var CompleteError = func(ctx CompleteContext) StageError {
	return NewStageErrorWithCode(errorCodeInternal, errors.New("Complete failed"))
}

func appendIfNotNil[T any](array []*T, items ...*T) []*T {
	for _, item := range items {
		if item != nil {
			array = append(array, item)
		}
	}
	return array
}

func addBreadcrumb(nodes ...*Node) {
	var nextNodes []*Node
	for _, n := range nodes {
		n.Cancel.appendBreadcrumb(cancelNodeType, n.breadcrumb)
		n.Compensate.appendBreadcrumb(compensateNodeType, n.breadcrumb)
		nextNodes = appendIfNotNil(nextNodes, n.Compensate, n.Cancel)
	}
	if len(nextNodes) > 0 {
		addBreadcrumb(nextNodes...)
	}
}
