package spark_v1

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
)

/************************************************************************/
// FACTORIES
/************************************************************************/

func newSetStageResultReq(jobKey, name string, data interface{}) (*sparkv1.SetStageResultRequest, error) {
	b, err := sparkv1.MarshalBinary(data)

	return &sparkv1.SetStageResultRequest{
		Key:  jobKey,
		Name: name,
		Data: b,
	}, err
}

func newVariable(name, mimeType string, value interface{}) (*sparkv1.Variable, error) {
	pbValue, err := sparkv1.MarshalBinary(value)
	if err != nil {
		return nil, fmt.Errorf("error creating variable named '%s': %w", name, err)
	}
	return &sparkv1.Variable{
		Data:     pbValue,
		MimeType: mimeType,
	}, nil
}

func newStageResultReq(jobKey, stageName string) *sparkv1.GetStageResultRequest {
	return &sparkv1.GetStageResultRequest{
		Name: stageName,
		Key:  jobKey,
	}
}

func newSetStageStatusReq(jobKey, stageName string, status sparkv1.StageStatus, err ...*sparkv1.Error) *sparkv1.SetStageStatusRequest {
	sssr := &sparkv1.SetStageStatusRequest{
		Name:   stageName,
		Key:    jobKey,
		Status: status,
	}
	if len(err) > 0 {
		sssr.Err = err[0]
	}
	return sssr
}

func newGetVariablesRequest(jobKey string, names ...string) *sparkv1.GetInputsRequest {
	vr := &sparkv1.GetInputsRequest{
		Key: jobKey,
	}
	vr.Names = append(vr.Names, names...)
	return vr
}

func newSetVariablesRequest(jobKey string, variables ...*Var) (*sparkv1.SetOutputsRequest, error) {
	m := map[string]*sparkv1.Variable{}
	for _, v := range variables {
		variable, err := newVariable(v.Name, v.MimeType, v.Value)
		if err != nil {
			return nil, err
		}
		m[v.Name] = variable
	}
	return &sparkv1.SetOutputsRequest{Key: jobKey, Variables: m}, nil
}

func newGetStageStatusReq(jobKey, stageName string) *sparkv1.GetStageStatusRequest {
	return &sparkv1.GetStageStatusRequest{Key: jobKey, Name: stageName}
}

/************************************************************************/
// INPUT
/************************************************************************/

type input struct {
	variable *sparkv1.Variable
	err      error
}

func newInput(variable *sparkv1.Variable, err error) *input {
	return &input{variable: variable, err: err}
}

func (i *input) String() string {
	b, _ := i.Raw()
	return string(b)
}

func (i *input) Raw() ([]byte, error) {
	if i.err != nil {
		return nil, i.err
	}

	return sparkv1.ConvertBytes(i.variable.Data, i.variable.MimeType)
}

func (i *input) Bind(a interface{}) error {
	if i.err != nil {
		return i.err
	}

	if err := sparkv1.UnmarshalBinaryTo(i.variable.Data, a, ""); err != nil {
		return err
	}

	return nil
}

/************************************************************************/
// BATCH INPUTS
/************************************************************************/

type inputs struct {
	vars map[string]*sparkv1.Variable
	err  error
}

func newInputs(err error, vars map[string]*sparkv1.Variable) Inputs {
	return &inputs{vars: vars, err: err}
}

func (v inputs) Get(name string) Bindable {
	found, ok := v.vars[name]
	if ok {
		return newInput(found, v.err)
	}
	err := v.err
	if err == nil {
		err = ErrInputVariableNotFound
	}
	return newInput(nil, v.err)
}

func (v inputs) Error() error {
	return v.err
}

/************************************************************************/
// STAGE RESULT
/************************************************************************/

type result struct {
	result *sparkv1.GetStageResultResponse
	err    error
}

func newResult(err error, r *sparkv1.GetStageResultResponse) Bindable {
	return &result{
		result: r,
		err:    err,
	}
}

func (r *result) Raw() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}

	return sparkv1.ConvertBytes(r.result.Data, "")
}

func (r *result) Bind(a interface{}) error {
	if r.err != nil {
		return r.err
	}

	return sparkv1.UnmarshalBinaryTo(r.result.Data, a, "")
}

func (r *result) String() string {
	b, _ := r.Raw()
	return string(b)
}
