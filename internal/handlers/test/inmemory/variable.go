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

func NewVariableHandler(t *testing.T) api.VariableHandler {
	return inMemoryVariableHandler{t: t, variables: make(map[string]*sdk_v1.Variable)}
}

func (i inMemoryVariableHandler) Get(jobKey, stage string, names ...string) ([]*sdk_v1.Variable, error) {
	var vars []*sdk_v1.Variable
	for _, n := range names {
		vars = append(vars, i.variables[i.key(n, stage, jobKey)])
	}
	if len(vars) == 0 {
		i.t.Fatalf("no variables found for the params: ")
	}
	return vars, nil
}

func (i inMemoryVariableHandler) Set(jobKey, stage string, variables ...*sdk_v1.Variable) error {
	for _, v := range variables {
		i.variables[i.key(v.Name, stage, jobKey)] = v
	}
	return nil
}

func (i inMemoryVariableHandler) key(name, stage, jobKey string) string {
	return fmt.Sprintf("%s_%s_%s", name, stage, jobKey)
}
