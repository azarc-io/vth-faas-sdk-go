package spark_test

import (
	"context"
	"encoding/json"
	spark "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-complex-example/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Should_Chain_Multiple_Inputs_And_Outputs(t *testing.T) {
	ctx := module_test_runner.NewTestJobContext(context.Background(), "test", "cid", "tid", sparkv1.ExecuteSparkInputs{
		"name_string": {
			Value:    []byte(`"Bob"`),
			MimeType: "application/text",
		},
	})

	worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark())
	assert.Nil(t, err)

	result, err := worker.Execute(ctx)
	if !assert.NoError(t, err) {
		return
	}

	t.Run("out1-struct", func(t *testing.T) {
		res := new(spark.ComplexType)
		assert.NoError(t, result.Bind("out1-struct", res))
		assert.Equal(t, "Message: Hello, Bob", string(res.ByteType))
		assert.Equal(t, int32(87), res.Int32Type)
	})

	t.Run("out2-int32", func(t *testing.T) {
		res := new(int32)
		assert.NoError(t, result.Bind("out2-int32", res))
		assert.Equal(t, int32(1), *res)
	})

	t.Run("out3-string", func(t *testing.T) {
		res := new(string)
		assert.NoError(t, result.Bind("out3-string", res))
		assert.Equal(t, "foobar", *res)
	})

	t.Run("out4-float", func(t *testing.T) {
		res := new(float32)
		assert.NoError(t, result.Bind("out4-float", res))
		assert.Equal(t, float32(54.221), *res)
	})

	t.Run("out5-bool", func(t *testing.T) {
		res := new(bool)
		assert.NoError(t, result.Bind("out5-bool", res))
		assert.Equal(t, true, *res)
	})

	t.Run("out6-array", func(t *testing.T) {
		res := new([]byte)
		assert.NoError(t, result.Bind("out6-array", res))
		assert.Equal(t, "my-byte-array", string(*res))
	})

	t.Run("out7-json-bytes", func(t *testing.T) {
		var res []byte
		assert.NoError(t, result.Bind("out7-json-bytes", &res))
		out := make(map[string]string)
		assert.NoError(t, json.Unmarshal(res, &out))
		assert.Equal(t, "bar", out["foo"])
	})

	t.Run("out8-json-string", func(t *testing.T) {
		var res string
		assert.NoError(t, result.Bind("out8-json-string", &res))
		out := make(map[string]string)
		assert.NoError(t, json.Unmarshal([]byte(res), &out))
		assert.Equal(t, "bar", out["foo"])
	})
}
