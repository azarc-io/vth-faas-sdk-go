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

package grpc

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type VariableHandler struct {
	client sdk_v1.ManagerServiceClient
}

func NewIOHandler(client sdk_v1.ManagerServiceClient) sdk_v1.IOHandler {
	return VariableHandler{client}
}

func (g VariableHandler) Inputs(jobKey string, names ...string) *sdk_v1.Inputs {
	variables, err := g.client.GetVariables(context.Background(), sdk_v1.NewGetVariablesRequest(jobKey, names...))
	if err != nil {
		return sdk_v1.NewInputs(err)
	}
	var vars []*sdk_v1.Variable
	for _, v := range variables.Variables {
		vars = append(vars, v)
	}
	return sdk_v1.NewInputs(err, vars...)
}

func (g VariableHandler) Input(jobKey, name string) *sdk_v1.Input {
	return g.Inputs(jobKey, name).Get(name)
}

func (g VariableHandler) Output(jobKey string, variables ...*handlers.Variable) error {
	request, err := sdk_v1.NewSetVariablesRequest(jobKey, variables...)
	if err != nil {
		return err
	}
	_, err = g.client.SetVariables(context.Background(), request)
	return err
}
