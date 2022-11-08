package inmemory

import (
	"fmt"
	"testing"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/models"

	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

type inMemoryIOHandler struct {
	variables map[string]*models.Variable
	t         *testing.T
}

func NewIOHandler(t *testing.T) v1.IOHandler {
	i := &inMemoryIOHandler{t: t, variables: map[string]*models.Variable{}}
	return i
}

func (i *inMemoryIOHandler) Inputs(jobKey string, names ...string) *v1.Inputs {
	var (
		vars []*v1.Variable
		err  error
	)
	for _, n := range names {
		key := i.key(jobKey, n)
		if v, ok := i.variables[key]; ok {
			var va *v1.Variable
			va, err = v1.NewVariable(v.Name, v.MimeType, v.Value)
			vars = append(vars, va)
		}
	}
	if len(vars) == 0 {
		i.t.Fatalf("no variables found for the params: ")
	}
	return v1.NewInputs(err, vars...)
}

func (i *inMemoryIOHandler) Input(jobKey, name string) *v1.Input {
	inputs := i.Inputs(jobKey, name)
	return inputs.Get(name)
}

func (i *inMemoryIOHandler) Output(jobKey string, variables ...*models.Variable) error {
	for _, v := range variables {
		i.variables[i.key(jobKey, v.Name)] = v
	}
	return nil
}

func (i *inMemoryIOHandler) key(jobKey, name string) string {
	return fmt.Sprintf("%s_%s", jobKey, name)
}
