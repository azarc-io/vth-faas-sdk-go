package sparkv1

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/nats-io/nats.go/jetstream"
)

type ioDataProvider struct {
	ctx          context.Context
	stageResults map[string]*BindableValue
	inputs       map[string]*BindableValue
	store        jetstream.ObjectStore
}

func (iodp *ioDataProvider) SetInitialInputs(inputs ExecuteSparkInputs) {
	if inputs != nil {
		iodp.inputs = inputs
	}
}

func (iodp *ioDataProvider) GetInputValue(name string) (*BindableValue, bool) {
	v, ok := iodp.inputs[name]
	return v, ok
}

func (iodp *ioDataProvider) LoadVariables(key string) error {
	b, err := iodp.store.GetBytes(iodp.ctx, key)
	if err != nil {
		if errors.Is(err, jetstream.ErrObjectNotFound) {
			if iodp.inputs == nil {
				iodp.inputs = make(map[string]*BindableValue)
			}
			return nil
		}
		return err
	}
	return json.Unmarshal(b, &iodp.inputs)
}

type bindableInput struct {
	iodp      *ioDataProvider
	stageName string
	mimeType  string
	data      []byte
	name      string
}

func NewIoDataProvider(ctx context.Context, store jetstream.ObjectStore) SparkDataIO {
	return &ioDataProvider{
		ctx:          ctx,
		store:        store,
		stageResults: make(map[string]*BindableValue),
		inputs:       make(map[string]*BindableValue),
	}
}

func (b *bindableInput) Bind(a any) error {
	data, err := b.GetValue()
	if err != nil {
		return err
	}

	// data already cached
	return NewBindable(Value{
		Value:    data,
		MimeType: b.mimeType,
	}).Bind(a)
}

func (b *bindableInput) GetValue() ([]byte, error) {
	if len(b.data) == 0 {
		// first fetch data
		if bv, ok := b.iodp.GetInputValue(b.name); ok {
			b.data = bv.Value
			return b.data, nil
		} else {
			return nil, fmt.Errorf("%w: error retrieving input data: stage (%s), name (%s)",
				ErrVariableNotFound,
				b.stageName,
				b.name,
			)
		}
	}

	return b.data, nil
}

func (b *bindableInput) GetMimeType() string {
	return b.mimeType
}

// deprecated
// todo: Must deprecate this as it can lead to unexpected issues
func (b *bindableInput) String() string {
	d, _ := b.GetValue()
	var s string
	err := json.Unmarshal(d, &s)
	if err != nil {
		return string(d)
	}
	return s
}

func (iodp *ioDataProvider) NewInput(name, stageName string, value *BindableValue) Bindable {
	bi := &bindableInput{
		iodp:      iodp,
		stageName: stageName,
		name:      name,
		mimeType:  value.MimeType,
	}

	return bi
}

func (iodp *ioDataProvider) NewOutput(stageName string, value *BindableValue) (Bindable, error) {
	iodp.stageResults[stageName] = value
	return value, nil
}

func (iodp *ioDataProvider) GetStageResult(stageName string) (Bindable, error) {
	if v, ok := iodp.stageResults[stageName]; !ok {
		return nil, errors.New("stage result not found")
	} else {
		return NewBindable(Value{
			Value:    v.Value,
			MimeType: v.MimeType,
		}), nil
	}
}

func (iodp *ioDataProvider) PutStageResult(stageName string, stageValue []byte) (Bindable, error) {
	iodp.stageResults[stageName] = NewBindable(Value{
		Value:    stageValue,
		MimeType: string(codec.MimeTypeText),
	})

	return iodp.stageResults[stageName], nil
}
