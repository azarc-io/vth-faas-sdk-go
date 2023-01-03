package spark_v1

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"testing"
)

type inMemoryIOHandler struct {
	variables map[string]*Var
	t         *testing.T
}

func NewInMemoryIOHandler(t *testing.T) TestIOHandler {
	i := &inMemoryIOHandler{t: t, variables: map[string]*Var{}}
	return i
}

func (i *inMemoryIOHandler) Inputs(ctx SparkContext, names ...string) Inputs {
	var (
		vars = map[string]*sparkv1.Variable{}
		err  error
	)
	for _, n := range names {
		key := i.key(ctx, n)
		if v, ok := i.variables[key]; ok {
			var va *sparkv1.Variable
			va, err = newVariable(v.Name, v.MimeType, v.Value)
			vars[v.Name] = va
		}
	}
	if len(vars) == 0 {
		i.t.Fatal("no variables found for the params: ")
	}
	return newInputs(err, vars)
}

func (i *inMemoryIOHandler) Input(ctx SparkContext, name string) Input {
	inputs := i.Inputs(ctx, name)
	return inputs.Get(name)
}

func (i *inMemoryIOHandler) Output(ctx SparkContext, variables ...*Var) error {
	for _, v := range variables {
		i.variables[i.key(ctx, v.Name)] = v
	}
	return nil
}

func (i *inMemoryIOHandler) SetVar(ctx SparkContext, v *Var) {
	i.variables[i.key(ctx, v.Name)] = v
}

func (i *inMemoryIOHandler) key(ctx SparkContext, name string) string {
	return fmt.Sprintf("%s_%s", ctx.JobKey(), name)
}
