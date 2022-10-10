package inmemory

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"testing"
)

type inMemoryVariableHandler struct {
	variables map[string]*sdk_v1.Variable
	t         *testing.T
}

func NewMockVariableHandler(t *testing.T) api.VariableHandler {
	return inMemoryVariableHandler{t: t, variables: make(map[string]*sdk_v1.Variable)}
}

func (i inMemoryVariableHandler) Get(name, stage, jobKey string) (*sdk_v1.Variable, error) {
	if variable, ok := i.variables[i.key(name, stage, jobKey)]; ok {
		return variable, nil
	}
	i.t.Fatalf("variable not found for params >> name: %s, jobKey: %s, stage: %s", name, jobKey, stage)
	return nil, nil
}

func (i inMemoryVariableHandler) Set(req *sdk_v1.SetVariableRequest) error {
	i.variables[i.key(req.Name, req.Stage, req.JobKey)] = req.Variable
	return nil
}

func (i inMemoryVariableHandler) key(name, stage, jobKey string) string {
	return fmt.Sprintf("%s_%s_%s", name, stage, jobKey)
}
