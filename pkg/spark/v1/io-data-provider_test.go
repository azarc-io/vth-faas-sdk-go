package sparkv1

import (
	"context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIoDataProvider(t *testing.T) {
	t.Run("Stage IO", func(t *testing.T) {
		iop := &ioDataProvider{
			ctx:          context.Background(),
			stageResults: map[string]*BindableValue{},
		}
		input, _ := codec.Encode("hello world")

		t.Run("put stage result: success", func(t *testing.T) {
			sr, err := iop.PutStageResult("dummy-stage", input)
			assert.NoError(t, err)
			assert.Equal(t, "hello world", sr.String())
		})

		t.Run("get stage result: success", func(t *testing.T) {
			sr, err := iop.GetStageResult("dummy-stage")
			assert.NoError(t, err)

			var res any
			assert.NoError(t, sr.Bind(&res))
			assert.Equal(t, "hello world", res)
		})

		t.Run("get stage result: fail", func(t *testing.T) {
			t.SkipNow()
			_, err := iop.GetStageResult("dummy-error")
			assert.ErrorContains(t, err, "error getting stage result (500): foo bar error")
		})
	})

	t.Run("Spark IO", func(t *testing.T) {
		iop := &ioDataProvider{
			ctx:          context.Background(),
			stageResults: map[string]*BindableValue{},
		}

		t.Run("get input: success", func(t *testing.T) {
			sr := iop.NewInput("c789", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte(`"foo bar"`),
			})

			assert.Equal(t, "foo bar", sr.String())
		})

		t.Run("get input: fail", func(t *testing.T) {
			t.SkipNow()
			sr := iop.NewInput("c567", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("foo bar"),
			})

			var a any
			err := sr.Bind(&a)
			assert.ErrorContains(t, err, "error retrieving input data (400): correlationID (c567), reference (my-input-ref): foo bar")
			assert.Nil(t, a)
		})

		t.Run("post output: success", func(t *testing.T) {
			sr, err := iop.NewOutput("c789", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("some data"),
			})

			assert.NoError(t, err)
			assert.Empty(t, sr.String())
			assert.Equal(t, []byte("some data"), sr.(*BindableValue).Value)
		})

		t.Run("post output: fail", func(t *testing.T) {
			t.SkipNow()
			sr, err := iop.NewOutput("c567", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("some data"),
			})

			assert.ErrorContains(t, err, "error setting output (400)")
			assert.Nil(t, sr)
		})
	})
}
