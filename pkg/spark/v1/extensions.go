package spark_v1

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/samber/lo"

	"google.golang.org/protobuf/types/known/structpb"
)

/************************************************************************/
// MARSHALLING
/************************************************************************/

type serdes struct {
	unmarshal func(value *structpb.Value, a any) error
	marshal   func(a any) (*structpb.Value, error)
}

var serdesMap = map[string]serdes{
	MimeTypeJSON: {
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
	NoMimeType: {
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

func getRawFromPb(data *structpb.Value) ([]byte, error) {
	switch data.Kind.(type) {
	case *structpb.Value_NullValue:
		return nil, nil
	case *structpb.Value_StringValue:
		return []byte(data.GetStringValue()), nil
	}

	return data.MarshalJSON()
}

/************************************************************************/
// VARIABLE EXTENSIONS
/************************************************************************/

func (x *Variable) Raw() ([]byte, error) {
	return x.Value.MarshalJSON()
}

func (x *Variable) Bind(a any) error {
	return serdesMap[x.MimeType].unmarshal(x.Value, a)
}

/************************************************************************/
// STAGE RESULT EXTENSIONS
/************************************************************************/

func (x *StageResult) Raw() ([]byte, error) {
	return x.Data.MarshalJSON()
}

func (x *StageResult) Bind(a any) error {
	return serdesMap[MimeTypeJSON].unmarshal(x.Data, a)
}

/************************************************************************/
// FACTORIES
/************************************************************************/

func newSetStageResultReq(jobKey, name string, data any) (*SetStageResultRequest, error) {
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

func newVariable(name, mimeType string, value any) (*Variable, error) {
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

func newSetJobStatusReq(key string, status JobStatus, err ...*Error) *SetJobStatusRequest {
	req := &SetJobStatusRequest{Key: key, Status: status}
	if len(err) > 0 {
		req.Err = err[0]
	}
	return req
}

func newStageResultReq(jobKey, stageName string) *GetStageResultRequest {
	return &GetStageResultRequest{
		Name:   stageName,
		JobKey: jobKey,
	}
}

func newSetStageStatusReq(jobKey, stageName string, status StageStatus, err ...*Error) *SetStageStatusRequest {
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

func newGetVariablesRequest(jobKey string, names ...string) *GetVariablesRequest {
	vr := &GetVariablesRequest{
		JobKey: jobKey,
	}
	vr.Name = append(vr.Name, names...)
	return vr
}

func newSetVariablesRequest(jobKey string, variables ...*Var) (*SetVariablesRequest, error) {
	m := map[string]*Variable{}
	for _, v := range variables {
		variable, err := newVariable(v.Name, v.MimeType, v.Value)
		if err != nil {
			return nil, err
		}
		m[v.Name] = variable
	}
	return &SetVariablesRequest{JobKey: jobKey, Variables: m}, nil
}

func newGetStageStatusReq(jobKey, stageName string) *GetStageStatusRequest {
	return &GetStageStatusRequest{JobKey: jobKey, Name: stageName}
}

/************************************************************************/
// INPUT
/************************************************************************/

type input struct {
	variable *Variable
	err      error
}

func (i *input) String() string {
	return i.variable.Value.GetStringValue()
}

func (i *input) Raw() ([]byte, error) {
	if i.err != nil {
		return nil, i.err
	}

	return getRawFromPb(i.variable.Value)
}

func (i *input) Bind(a any) error {
	if i.err != nil {
		return i.err
	}
	return serdesMap[i.variable.MimeType].unmarshal(i.variable.Value, a)
}

/************************************************************************/
// BATCH INPUTS
/************************************************************************/

type inputs struct {
	vars []*Variable
	err  error
}

func newInputs(err error, vars ...*Variable) Inputs {
	return &inputs{vars: vars, err: err}
}

func (v inputs) Get(name string) Bindable {
	found, ok := lo.Find(v.vars, func(variable *Variable) bool {
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
	result *StageResult
	err    error
}

func newResult(err error, r *StageResult) Bindable {
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
		return getRawFromPb(r.result.GetData())
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
