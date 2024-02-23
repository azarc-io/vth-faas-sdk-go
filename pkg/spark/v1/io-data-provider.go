package sparkv1

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/nats-io/nats.go"
)

type ioDataProvider struct {
	ctx          context.Context
	bucket       string
	nc           *nats.Conn
	stageResults map[string]*BindableValue
}

type bindableInput struct {
	iodp      *ioDataProvider
	stageName string
	mimeType  string
	data      []byte
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

func (iodp *ioDataProvider) NewInput(stageName string, value *BindableValue) Bindable {
	return &bindableInput{
		data:      value.Value,
		iodp:      iodp,
		stageName: stageName,
		mimeType:  value.MimeType,
	}
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
