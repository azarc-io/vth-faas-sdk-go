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
	"errors"
	"fmt"
)

type validateFn struct {
	fn   func(n *node) error
	errs []error
}

func (v *validateFn) exec(n *node) {
	if err := v.fn(n); err != nil {
		v.errs = append(v.errs, err)
	}
}

func (c *chainBuilder) validate(fns []*validateFn, nodes ...*node) error {
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
	var errs []error
	for _, fn := range fns {
		if len(fn.errs) > 0 {
			errs = append(errs, fn.errs...)
		}
	}
	return aggregateValidationError(errs)
}

func aggregateValidationError(errs []error) error {
	if len(errs) < 1 {
		return nil
	}
	e := errors.New("chain validation error : ")
	for _, err := range errs {
		e = fmt.Errorf("%v\n\t%w", e, err)
	}
	return e
}

var atLeastOneStagePerNodeValidator = &validateFn{
	fn: func(n *node) error {
		if len(n.stages) < 1 {
			return fmt.Errorf("no stage defined for node: %s", n.breadcrumb)
		}
		return nil
	},
}

var stageNamesMustNoBeEmpty = &validateFn{
	fn: func(n *node) error {
		var stagesFromNodes []string
		for _, stg := range n.stages {
			stagesFromNodes = append(stagesFromNodes, stg.name)
		}
		if n.complete != nil {
			stagesFromNodes = append(stagesFromNodes, n.complete.name)
		}
		for _, name := range stagesFromNodes {
			if name == "" {
				return fmt.Errorf("stage with empty name <\"\"> found at '%s'", n.breadcrumb)
			}
		}
		return nil
	},
}

var uniqueStageNamesValidator = func() *validateFn {
	stageNames := map[string]string{}
	return &validateFn{
		fn: func(n *node) error {
			var stagesFromNodes []string
			for _, stg := range n.stages {
				stagesFromNodes = append(stagesFromNodes, stg.name)
			}
			if n.complete != nil {
				stagesFromNodes = append(stagesFromNodes, n.complete.name)
			}
			for _, stageName := range stagesFromNodes {
				if stageName == "" {
					return nil
				}
				if bc, ok := stageNames[stageName]; ok {
					return fmt.Errorf("unique stage name restriction violated: a stage or complete stage in '%s' and '%s' have the same name: '%s'", bc, n.breadcrumb, stageName)
				} else {
					stageNames[stageName] = n.breadcrumb
				}
			}
			return nil
		},
	}
}
