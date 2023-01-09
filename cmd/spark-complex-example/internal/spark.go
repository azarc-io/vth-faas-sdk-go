package spark

import (
	"fmt"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

type ComplexType struct {
	ByteType  []byte
	Int32Type int32
}

type Spark struct {
}

func (s Spark) Init(ctx sparkv1.InitContext) error {
	return nil
}

func (s Spark) Stop() {

}

func (s Spark) BuildChain(b sparkv1.Builder) sparkv1.Chain {
	return b.NewChain("chain-1").
		Stage("stage-1", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			var name string
			if err := ctx.Input("name_string").Bind(&name); err != nil {
				return nil, sparkv1.NewStageError(err)
			}
			return fmt.Sprintf("Hello, %s", name), nil
		}).
		Stage("stage-2", func(ctx sparkv1.StageContext) (any, sparkv1.StageError) {
			var message string
			if err := ctx.StageResult("stage-1").Bind(&message); err != nil {
				return nil, sparkv1.NewStageError(err)
			}

			return ComplexType{
				ByteType:  []byte(fmt.Sprintf("Message: %s", message)),
				Int32Type: 87,
			}, nil
		}).
		Complete(func(ctx sparkv1.CompleteContext) sparkv1.StageError {
			var out ComplexType

			// fetch the output of stage-1
			if err := ctx.StageResult("stage-2").Bind(&out); err != nil {
				return sparkv1.NewStageError(err)
			}

			// write the output of the spark
			if err := ctx.Output(
				sparkv1.NewVar("out1-struct", codec.MimeTypeJson, out),
				sparkv1.NewVar("out2-int32", codec.MimeTypeJson, int32(1)),
				sparkv1.NewVar("out3-string", codec.MimeTypeJson, "foobar"),
				sparkv1.NewVar("out4-float", codec.MimeTypeJson, 54.221),
				sparkv1.NewVar("out5-bool", codec.MimeTypeJson, true),
				sparkv1.NewVar("out6-array", codec.MimeTypeOctetStream.WithType("text"), []byte("my-byte-array")),
				sparkv1.NewVar("out7-json-bytes", codec.MimeTypeJson.WithType("text"), []byte(`{"foo": "bar"}`)),
				sparkv1.NewVar("out8-json-string", codec.MimeTypeJson.WithType("text"), `{"foo": "bar"}`),
			); err != nil {
				return sparkv1.NewStageError(err)
			}

			return nil
		})
}

// NewSpark creates a Spark
func NewSpark() sparkv1.Spark {
	return new(Spark)
}
