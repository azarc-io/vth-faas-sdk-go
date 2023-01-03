package spark_v1

import (
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/internal/common"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/internal/gen/azarc/sdk/spark/v1"
	"os"
)

/************************************************************************/
// FACTORIES
/************************************************************************/

func newSetStageResultReq(ctx SparkContext, name string, data interface{}) (*sparkv1.SetStageResultRequest, error) {
	b, err := sparkv1.MarshalBinary(data)

	req := &sparkv1.SetStageResultRequest{
		Key:      ctx.JobKey(),
		Name:     name,
		Data:     b,
		Metadata: &sparkv1.RequestMetadata{Metadata: ctx.RequestMetadata()},
	}

	// include known env vars
	updateMetadata(req.Metadata)

	return req, err
}

func newVariable(name, mimeType string, value interface{}) (*sparkv1.Variable, error) {
	var pbValue []byte
	var err error

	switch v := value.(type) {
	case string:
		pbValue = []byte(v)
	default:
		pbValue, err = sparkv1.MarshalBinary(v)
	}
	if err != nil {
		return nil, fmt.Errorf("error creating variable named '%s': %w", name, err)
	}
	return &sparkv1.Variable{
		Data:     pbValue,
		MimeType: mimeType,
	}, nil
}

func newStageResultReq(ctx SparkContext, stageName string) *sparkv1.GetStageResultRequest {
	req := &sparkv1.GetStageResultRequest{
		Name:     stageName,
		Key:      ctx.JobKey(),
		Metadata: &sparkv1.RequestMetadata{Metadata: ctx.RequestMetadata()},
	}

	// include known env vars
	updateMetadata(req.Metadata)

	return req
}

func newSetStageStatusReq(ctx SparkContext, stageName string, status sparkv1.StageStatus, err ...*sparkv1.Error) *sparkv1.SetStageStatusRequest {
	sssr := &sparkv1.SetStageStatusRequest{
		Name:     stageName,
		Key:      ctx.JobKey(),
		Status:   status,
		Metadata: &sparkv1.RequestMetadata{Metadata: ctx.RequestMetadata()},
	}

	// include known env vars
	updateMetadata(sssr.Metadata)

	if len(err) > 0 {
		sssr.Err = err[0]
	}

	return sssr
}

func newGetVariablesRequest(ctx SparkContext, names ...string) *sparkv1.GetInputsRequest {
	vr := &sparkv1.GetInputsRequest{
		Key:      ctx.JobKey(),
		Metadata: &sparkv1.RequestMetadata{Metadata: ctx.RequestMetadata()},
	}
	// include known env vars
	updateMetadata(vr.Metadata)
	vr.Names = append(vr.Names, names...)
	return vr
}

// updateMetadata populates request metadata with values from env vars if set
func updateMetadata(metadata *sparkv1.RequestMetadata) {
	if metadata.Metadata == nil {
		metadata.Metadata = map[string]string{}
	}

	if val, ok := os.LookupEnv("TASK_ID"); ok {
		metadata.Metadata["taskId"] = val
	}
}

func newSetVariablesRequest(ctx SparkContext, variables ...*Var) (*sparkv1.SetOutputsRequest, error) {
	m := map[string]*sparkv1.Variable{}
	for _, v := range variables {
		variable, err := newVariable(v.Name, v.MimeType, v.Value)
		if err != nil {
			return nil, err
		}
		m[v.Name] = variable
	}
	req := &sparkv1.SetOutputsRequest{
		Key:       ctx.JobKey(),
		Variables: m,
		Metadata:  &sparkv1.RequestMetadata{Metadata: ctx.RequestMetadata()},
	}

	// include known env vars
	updateMetadata(req.Metadata)

	return req, nil
}

func newGetStageStatusReq(ctx SparkContext, stageName string) *sparkv1.GetStageStatusRequest {
	req := &sparkv1.GetStageStatusRequest{
		Key:      ctx.JobKey(),
		Name:     stageName,
		Metadata: &sparkv1.RequestMetadata{Metadata: ctx.RequestMetadata()},
	}

	// include known env vars
	updateMetadata(req.Metadata)

	return req
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

	if err := sparkv1.UnmarshalBinaryTo(i.variable.Data, a, i.variable.MimeType); err != nil {
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

	return sparkv1.ConvertBytes(r.result.Data, common.MimeTypeJSON)
}

func (r *result) Bind(a interface{}) error {
	if r.err != nil {
		return r.err
	}

	return sparkv1.UnmarshalBinaryTo(r.result.Data, a, common.MimeTypeJSON)
}

func (r *result) String() string {
	b, _ := r.Raw()
	return string(b)
}
