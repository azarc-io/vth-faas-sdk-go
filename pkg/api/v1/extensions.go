package sdk_v1

import (
	"encoding/json"
	"google.golang.org/protobuf/types/known/structpb"
)

// TODO yaml, xml, json, toml, csv <- just the ones that golang support
// we need to use the correct encoder based on the mime_type field of the message

func (x *Variable) Raw() []byte {
	return nil
}

func (x *Variable) Bind(a any) error {
	return nil
}

func NewSetStageResultReq(jobKey, name string, data any) *SetStageResultRequest {
	b, err := json.Marshal(data)
	if err != nil {
		panic("this should panic?") // TODO proper error handling
	}
	pbValue, err := structpb.NewValue(b)
	if err != nil {
		println("ERROR CREATING NEW VARIABLE::>> ", err.Error()) // TODO fix me
	}
	return &SetStageResultRequest{
		JobKey: jobKey,
		Name:   name,
		Result: &StageResult{Data: pbValue},
	}
}

func NewVariable(name, mimeType string, value any) *Variable {
	pbValue, err := structpb.NewValue(value)
	if err != nil {
		println("ERROR CREATING NEW VARIABLE::>> ", err.Error()) // TODO fix me
	}
	return &Variable{
		Name:     name,
		Value:    pbValue,
		MimeType: mimeType,
	}
}

func NewSetJobStatusReq(key string, status JobStatus, err ...Error) *SetJobStatusRequest {
	req := &SetJobStatusRequest{Key: key, Status: status}
	if len(err) > 0 {
		req.Err = &err[0]
	}
	return req
}

func NewStageResultReq(jobKey, stageName string) *GetStageResultRequest {
	return &GetStageResultRequest{
		Name:   stageName,
		JobKey: jobKey,
	}
}

func NewSetStageStatusReq(jobKey, stageName string, status StageStatus, err ...Error) *SetStageStatusRequest {
	sssr := &SetStageStatusRequest{
		Name:   stageName,
		JobKey: jobKey,
		Status: status,
	}
	if len(err) > 0 {
		sssr.Err = &err[0]
	}
	return sssr
}

func Ptr[T any](t T) *T {
	return &t
}
