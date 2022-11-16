package spark_v1

import (
	"errors"
	"fmt"
)

var ErrValidationErr = errors.New("chain validation error : ")

func appendValidationErrorMessage(err error, message string) error {
	return fmt.Errorf("%w\n\t%s", err, message)
}

type validateFn struct {
	fn   func(n *node) string
	errs []string
}

func (v *validateFn) exec(n *node) {
	if err := v.fn(n); err != "" {
		v.errs = append(v.errs, err)
	}
}

func (c *ChainBuilder) validate(fns []*validateFn, nodes ...*node) error {
	var nextNodes []*node
	for _, n := range nodes {
		for _, fn := range fns {
			fn.exec(n)
		}
		nextNodes = appendIfNotNil(nextNodes, n.compensate, n.cancel)
	}
	if len(nextNodes) > 0 {
		return c.validate(fns, nextNodes...)
	}
	var errs []string
	for _, fn := range fns {
		if len(fn.errs) > 0 {
			errs = append(errs, fn.errs...)
		}
	}
	return aggregateValidationError(errs)
}

func aggregateValidationError(errs []string) error {
	if len(errs) < 1 {
		return nil
	}
	e := ErrValidationErr
	for _, err := range errs {
		e = appendValidationErrorMessage(e, err)
	}
	return e
}

var atLeastOneStagePerNodeValidator = &validateFn{
	fn: func(n *node) string {
		if len(n.stages) < 1 {
			return fmt.Sprintf("no stage defined for node: %s", n.breadcrumb)
		}
		return ""
	},
}

var stageNamesMustNoBeEmpty = &validateFn{
	fn: func(n *node) string {
		var stagesFromNodes []string
		for _, stg := range n.stages {
			stagesFromNodes = append(stagesFromNodes, stg.name)
		}
		if n.complete != nil {
			stagesFromNodes = append(stagesFromNodes, n.complete.name)
		}
		for _, name := range stagesFromNodes {
			if name == "" {
				return fmt.Sprintf("stage with empty name <\"\"> found at '%s'", n.breadcrumb)
			}
		}
		return ""
	},
}

var uniqueStageNamesValidator = func() *validateFn {
	stageNames := map[string]string{}
	return &validateFn{
		fn: func(n *node) string {
			var stagesFromNodes []string
			for _, stg := range n.stages {
				stagesFromNodes = append(stagesFromNodes, stg.name)
			}
			if n.complete != nil {
				stagesFromNodes = append(stagesFromNodes, n.complete.name)
			}
			for _, stageName := range stagesFromNodes {
				if stageName == "" {
					return ""
				}
				if bc, ok := stageNames[stageName]; ok {
					return fmt.Sprintf("unique stage name restriction violated: a stage or complete stage in '%s' and '%s' have the same name: '%s'", bc, n.breadcrumb, stageName)
				}
				stageNames[stageName] = n.breadcrumb
			}
			return ""
		},
	}
}
