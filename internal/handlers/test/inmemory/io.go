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

package inmemory

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"testing"
)

type inMemoryIOHandler struct {
	variables map[string]*handlers.Variable
	t         *testing.T
}

func NewIOHandler(t *testing.T) sdk_v1.IOHandler {
	i := &inMemoryIOHandler{t: t, variables: map[string]*handlers.Variable{}}
	return i
}

func (i *inMemoryIOHandler) Inputs(jobKey string, names ...string) *sdk_v1.Inputs {
	var (
		vars []*sdk_v1.Variable
		err  error
	)
	for _, n := range names {
		key := i.key(jobKey, n)
		if v, ok := i.variables[key]; ok {
			var va *sdk_v1.Variable
			va, err = sdk_v1.NewVariable(v.Name, v.MimeType, v.Value)
			vars = append(vars, va)
		}
	}
	if len(vars) == 0 {
		i.t.Fatalf("no variables found for the params: ")
	}
	return sdk_v1.NewInputs(err, vars...)
}

func (i *inMemoryIOHandler) Input(jobKey, name string) *sdk_v1.Input {
	inputs := i.Inputs(jobKey, name)
	return inputs.Get(name)
}

func (i *inMemoryIOHandler) Output(jobKey string, variables ...*handlers.Variable) error {
	for _, v := range variables {
		i.variables[i.key(jobKey, v.Name)] = v
	}
	return nil
}

func (i *inMemoryIOHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}
