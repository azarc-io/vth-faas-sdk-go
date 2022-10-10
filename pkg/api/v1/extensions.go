package sdk_v1

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/structpb"
)

// TODO yaml, xml, json, toml, csv <- just the ones that golang support
// we need to use the correct encoder based on the mime_type field of the message

func (x *Variable) Raw() ([]byte, error) {
	return x.Value.MarshalJSON()
}

func (x *Variable) Bind(a any) error {
	return serdesMap[x.MimeType].unmarshal(x.Value, a)
}

func NewSetStageResultReq(jobKey, name string, data any) (*SetStageResultRequest, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	pbValue, err := structpb.NewValue(b)
	if err != nil {
		return nil, err
	}
	return &SetStageResultRequest{
		JobKey: jobKey,
		Name:   name,
		Result: &StageResult{Data: pbValue},
	}, nil
}

func NewVariable(name, mimeType string, value any) (*Variable, error) {
	pbValue, err := serdesMap[mimeType].marshal(value)
	if err != nil {
		return nil, err
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

func NewGetVariablesRequest(jobKey, stage string, names ...string) *GetVariablesRequest {
	vr := &GetVariablesRequest{
		Stage:  stage,
		JobKey: jobKey,
	}
	for _, name := range names {
		vr.Name = append(vr.Name, name)
	}
	return vr
}

func NewSetVariablesRequest(jobKey, stage string, variables ...*Variable) *SetVariablesRequest {
	m := map[string]*Variable{}
	for _, v := range variables {
		m[v.Name] = v
	}
	return &SetVariablesRequest{Stage: stage, JobKey: jobKey, Variables: m}
}

func NewGetStageStatusReq(jobKey, stageName string) *GetStageStatusRequest {
	return &GetStageStatusRequest{JobKey: jobKey, Name: stageName}
}

func Ptr[T any](t T) *T {
	return &t
}

type serdes struct {
	unmarshal func(value *structpb.Value, a any) error
	marshal   func(a any) (*structpb.Value, error)
}

var serdesMap = map[string]serdes{
	"application/json": {
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
				v := map[string]interface{}{}
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
