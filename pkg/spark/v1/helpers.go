package spark_v1

import (
	"encoding/json"
	"errors"
)

var CompleteSuccess = func(ctx CompleteContext) StageError {
	return nil
}

var CompleteError = func(ctx CompleteContext) StageError {
	return NewStageError(errors.New("complete failed"))
}

func appendIfNotNil[T any](array []*T, items ...*T) []*T {
	for _, item := range items {
		if item != nil {
			array = append(array, item)
		}
	}
	return array
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
	}
}

func toMap(data any) (map[string]any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m = map[string]any{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}