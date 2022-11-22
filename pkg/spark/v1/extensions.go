package spark_v1

import (
	"fmt"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/spark/v1"
	"reflect"

	"github.com/samber/lo"

	"google.golang.org/protobuf/types/known/structpb"
)

/************************************************************************/
// FACTORIES
/************************************************************************/

func newSetStageResultReq(jobKey, name string, data any) (*sparkv1.SetStageResultRequest, error) {
	pbValue, err := structpb.NewValue(data)
	if err != nil {
		switch reflect.TypeOf(data).Kind() { //nolint:exhaustive
		case reflect.Slice, reflect.Array:
			arr := reflect.ValueOf(data)
			var anyArr []any
			for i := 0; i < arr.Len(); i++ {
				m, err := toMap(arr.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				anyArr = append(anyArr, m)
			}
			pbValue, err = structpb.NewValue(anyArr)
			if err != nil {
				return nil, err
			}
		default:
			m, err := toMap(data)
			if err != nil {
				return nil, err
			}
			pbValue, err = structpb.NewValue(m)
			if err != nil {
				return nil, err
			}
		}
	}
	return &sparkv1.SetStageResultRequest{
		JobKey: jobKey,
		Name:   name,
		Result: &sparkv1.StageResult{Data: pbValue},
	}, nil
}

func newVariable(name, mimeType string, value any) (*sparkv1.Variable, error) {
	pbValue, err := sparkv1.SerdesMap[mimeType].Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("error creating variable named '%s': %w", name, err)
	}
	return &sparkv1.Variable{
		Name:     name,
		Value:    pbValue,
		MimeType: mimeType,
	}, nil
}

func newSetJobStatusReq(key string, status sparkv1.JobStatus, err ...*sparkv1.Error) *sparkv1.SetJobStatusRequest {
	req := &sparkv1.SetJobStatusRequest{Key: key, Status: status}
	if len(err) > 0 {
		req.Err = err[0]
	}
	return req
}

func newStageResultReq(jobKey, stageName string) *sparkv1.GetStageResultRequest {
	return &sparkv1.GetStageResultRequest{
		Name:   stageName,
		JobKey: jobKey,
	}
}

func newSetStageStatusReq(jobKey, stageName string, status sparkv1.StageStatus, err ...*sparkv1.Error) *sparkv1.SetStageStatusRequest {
	sssr := &sparkv1.SetStageStatusRequest{
		Name:   stageName,
		JobKey: jobKey,
		Status: status,
	}
	if len(err) > 0 {
		sssr.Err = err[0]
	}
	return sssr
}

func newGetVariablesRequest(jobKey string, names ...string) *sparkv1.GetVariablesRequest {
	vr := &sparkv1.GetVariablesRequest{
		JobKey: jobKey,
	}
	vr.Name = append(vr.Name, names...)
	return vr
}

func newSetVariablesRequest(jobKey string, variables ...*Var) (*sparkv1.SetVariablesRequest, error) {
	m := map[string]*sparkv1.Variable{}
	for _, v := range variables {
		variable, err := newVariable(v.Name, v.MimeType, v.Value)
		if err != nil {
			return nil, err
		}
		m[v.Name] = variable
	}
	return &sparkv1.SetVariablesRequest{JobKey: jobKey, Variables: m}, nil
}

func newGetStageStatusReq(jobKey, stageName string) *sparkv1.GetStageStatusRequest {
	return &sparkv1.GetStageStatusRequest{JobKey: jobKey, Name: stageName}
}

/************************************************************************/
// INPUT
/************************************************************************/

type input struct {
	variable *sparkv1.Variable
	err      error
}

func (i *input) String() string {
	return i.variable.Value.GetStringValue()
}

func (i *input) Raw() ([]byte, error) {
	if i.err != nil {
		return nil, i.err
	}

	return sparkv1.GetRawFromPb(i.variable.Value)
}

func (i *input) Bind(a any) error {
	if i.err != nil {
		return i.err
	}
	return sparkv1.SerdesMap[i.variable.MimeType].Unmarshal(i.variable.Value, a)
}

/************************************************************************/
// BATCH INPUTS
/************************************************************************/

type inputs struct {
	vars []*sparkv1.Variable
	err  error
}

func newInputs(err error, vars ...*sparkv1.Variable) Inputs {
	return &inputs{vars: vars, err: err}
}

func (v inputs) Get(name string) Bindable {
	found, ok := lo.Find(v.vars, func(variable *sparkv1.Variable) bool {
		return variable.Name == name
	})
	if ok {
		return &input{found, v.err}
	}
	err := v.err
	if err == nil {
		err = ErrInputVariableNotFound
	}
	return &input{nil, err}
}

func (v inputs) Error() error {
	return v.err
}

/************************************************************************/
// STAGE RESULT
/************************************************************************/

type result struct {
	result *sparkv1.StageResult
	err    error
}

func newResult(err error, r *sparkv1.StageResult) Bindable {
	return &result{
		result: r,
		err:    err,
	}
}

func (r *result) Raw() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}

	if r.result.GetData() != nil {
		return sparkv1.GetRawFromPb(r.result.GetData())
	}

	return r.result.Raw()
}

func (r *result) Bind(a any) error {
	if r.err != nil {
		return r.err
	}
	return r.result.Bind(a)
}

func (r *result) String() string {
	return r.result.GetData().GetStringValue()
}
