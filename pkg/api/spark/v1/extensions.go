package sdk_v1

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/models"

	"github.com/samber/lo"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	"google.golang.org/protobuf/types/known/structpb"
)

var ErrInputVariableNotFound = errors.New("input variable not found")

func (x *Variable) Raw() ([]byte, error) {
	return x.Value.MarshalJSON()
}

func (x *Variable) Bind(a any) error {
	return serdesMap[x.MimeType].unmarshal(x.Value, a)
}

func (x *StageResult) Raw() ([]byte, error) {
	return x.Data.MarshalJSON()
}

func (x *StageResult) Bind(a any) error {
	return serdesMap[api.MimeTypeJSON].unmarshal(x.Data, a)
}

func NewSetStageResultReq(jobKey, name string, data any) (*SetStageResultRequest, error) {
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
	return &SetStageResultRequest{
		JobKey: jobKey,
		Name:   name,
		Result: &StageResult{Data: pbValue},
	}, nil
}

func toMap(data any) (map[string]any, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	var m = map[string]any{}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func NewVariable(name, mimeType string, value any) (*Variable, error) {
	pbValue, err := serdesMap[mimeType].marshal(value)
	if err != nil {
		return nil, fmt.Errorf("error creating variable named '%s': %w", name, err)
	}
	return &Variable{
		Name:     name,
		Value:    pbValue,
		MimeType: mimeType,
	}, nil
}

func NewSetJobStatusReq(key string, status JobStatus, err ...*Error) *SetJobStatusRequest {
	req := &SetJobStatusRequest{Key: key, Status: status}
	if len(err) > 0 {
		req.Err = err[0]
	}
	return req
}

func NewStageResultReq(jobKey, stageName string) *GetStageResultRequest {
	return &GetStageResultRequest{
		Name:   stageName,
		JobKey: jobKey,
	}
}

func NewSetStageStatusReq(jobKey, stageName string, status StageStatus, err ...*Error) *SetStageStatusRequest {
	sssr := &SetStageStatusRequest{
		Name:   stageName,
		JobKey: jobKey,
		Status: status,
	}
	if len(err) > 0 {
		sssr.Err = err[0]
	}
	return sssr
}

func NewGetVariablesRequest(jobKey string, names ...string) *GetVariablesRequest {
	vr := &GetVariablesRequest{
		JobKey: jobKey,
	}
	vr.Name = append(vr.Name, names...)
	return vr
}

func NewSetVariablesRequest(jobKey string, variables ...*models.Variable) (*SetVariablesRequest, error) {
	m := map[string]*Variable{}
	for _, v := range variables {
		variable, err := NewVariable(v.Name, v.MimeType, v.Value)
		if err != nil {
			return nil, err
		}
		m[v.Name] = variable
	}
	return &SetVariablesRequest{JobKey: jobKey, Variables: m}, nil
}

func NewGetStageStatusReq(jobKey, stageName string) *GetStageStatusRequest {
	return &GetStageStatusRequest{JobKey: jobKey, Name: stageName}
}

type serdes struct {
	unmarshal func(value *structpb.Value, a any) error
	marshal   func(a any) (*structpb.Value, error)
}

var serdesMap = map[string]serdes{
	api.MimeTypeJSON: {
		unmarshal: func(value *structpb.Value, a any) error {
			data, err := value.MarshalJSON()
			if err != nil {
				return err
			}
			return json.Unmarshal(data, a)
		},
		marshal: func(a any) (*structpb.Value, error) {
			value, err := structpb.NewValue(a)
			if err != nil {
				b, err := json.Marshal(a)
				if err != nil {
					return nil, err
				}
				v := map[string]any{}
				err = json.Unmarshal(b, &v)
				if err != nil {
					return nil, err
				}
				return structpb.NewValue(v)
			}
			return value, nil
		},
	},
}

type Input struct {
	variable *Variable
	err      error
}

func (i *Input) Raw() ([]byte, error) {
	if i.err != nil {
		return nil, i.err
	}
	return i.variable.Value.MarshalJSON()
}

func (i *Input) Bind(a any) error {
	if i.err != nil {
		return i.err
	}
	return serdesMap[i.variable.MimeType].unmarshal(i.variable.Value, a)
}

type Inputs struct {
	vars []*Variable
	err  error
}

func NewInputs(err error, vars ...*Variable) *Inputs {
	return &Inputs{vars: vars, err: err}
}

func (v Inputs) Get(name string) *Input {
	found, ok := lo.Find(v.vars, func(variable *Variable) bool {
		return variable.Name == name
	})
	if ok {
		return &Input{found, v.err}
	}
	err := v.err
	if err == nil {
		err = ErrInputVariableNotFound
	}
	return &Input{nil, err}
}

func (v Inputs) Error() error {
	return v.err
}

type Result struct {
	result *StageResult
	err    error
}

func NewResult(err error, result *StageResult) *Result {
	return &Result{
		result: result,
		err:    err,
	}
}

func (r *Result) Raw() ([]byte, error) {
	if r.err != nil {
		return nil, r.err
	}
	return r.result.Raw()
}

func (r *Result) Bind(a any) error {
	if r.err != nil {
		return r.err
	}
	return r.result.Bind(a)
}
