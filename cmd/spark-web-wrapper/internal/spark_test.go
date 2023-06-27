package spark_test

import (
	"context"
	_ "embed"
	spark "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-web-wrapper/internal"
	helpers "github.com/azarc-io/vth-faas-sdk-go/internal/test_helpers"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1/test"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

//go:embed fixtures/example1-s1-result.json
var example1_s1_result []byte

//go:embed fixtures/example1-s1-expected-input.json
var example1_s1_expected_input []byte

//go:embed fixtures/example1-s2-result.json
var example1_s2_result []byte

//go:embed fixtures/example1-s2-expected-input.json
var example1_s2_expected_input []byte

//go:embed fixtures/example1-c1-result.json
var example1_c1_result []byte

//go:embed fixtures/example1-c1-expected-input.json
var example1_c1_expected_input []byte

//go:embed fixtures/example1-spec.json
var example1_Spec []byte

//go:embed fixtures/example1-stage-error-with-retry.json
var example1_stage_error_with_retry []byte

func TestShouldInitialiseSparkAndRunMultiStages(t *testing.T) {
	svr := helpers.GetTestHttpServerWithRequests(t, []helpers.Request{
		{http.MethodGet, 200, "/basepath/spec", example1_Spec, nil, nil},
		{http.MethodPost, 200, "/basepath/init", nil, []byte(`{"Foo":"my-bar-from-config"}`), nil},
		{http.MethodPost, 200, "/basepath/stages/My-Stage-1", example1_s1_result, example1_s1_expected_input, nil},
		{http.MethodPost, 200, "/basepath/stages/My-Stage-2", example1_s2_result, example1_s2_expected_input, nil},
		{http.MethodPost, 200, "/basepath/complete/My-Complete", example1_c1_result, example1_c1_expected_input, nil},
	})

	ctx := module_test_runner.NewTestJobContext(context.Background(), "test", "cid", "tid", module_test_runner.Inputs{
		"myKey": {
			Value:    "anything",
			MimeType: "",
		},
		"foo": {
			Value:    12345,
			MimeType: codec.MimeTypeJson,
		},
	})

	worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark(svr.URL+"/basepath"))
	assert.Nil(t, err)

	result, err := worker.Execute(ctx, sparkv1.WithSparkConfig(map[string]any{"Foo": "my-bar-from-config"}))
	if !assert.Nil(t, err) {
		return
	}

	t.Run("supports result: string", func(t *testing.T) {
		var v string
		assert.NoError(t, result.Bind("v-string", &v))
		assert.Equal(t, "hello back", v)
	})
	t.Run("supports result: int", func(t *testing.T) {
		var v int
		assert.NoError(t, result.Bind("v-int", &v))
		assert.Equal(t, 1234, v)
	})
	t.Run("supports result: float", func(t *testing.T) {
		var v float64
		assert.NoError(t, result.Bind("v-float", &v))
		assert.Equal(t, 789.123, v)
	})
	t.Run("supports result: bool", func(t *testing.T) {
		var v bool
		assert.NoError(t, result.Bind("v-bool", &v))
		assert.Equal(t, true, v)
	})
	t.Run("supports result: object", func(t *testing.T) {
		var v map[string]any
		assert.NoError(t, result.Bind("v-object", &v))
		assert.Equal(t, map[string]any{"foo": "bar"}, v)
	})

	worker.AssertStageCompleted("My-Stage-1")
	worker.AssertStageCompleted("My-Stage-2")
	worker.AssertStageCompleted("main_complete")
	worker.AssertStageOrder("My-Stage-1", "My-Stage-2")
}

func TestShouldErrorOnStage1(t *testing.T) {
	t.Run("User Error: 500", func(t *testing.T) {
		svr := helpers.GetTestHttpServerWithRequests(t, []helpers.Request{
			{http.MethodGet, 200, "/basepath/spec", example1_Spec, nil, nil},
			{http.MethodPost, 200, "/basepath/init", nil, nil, nil},
			{http.MethodPost, 500, "/basepath/stages/My-Stage-1", example1_stage_error_with_retry, nil, nil},
		})

		ctx := module_test_runner.NewTestJobContext(context.Background(), "test", "cid", "tid", module_test_runner.Inputs{})
		worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark(svr.URL+"/basepath"))
		assert.Nil(t, err)

		_, err = worker.Execute(ctx)
		if !assert.Error(t, err) {
			return
		}
		assert.IsType(t, &sparkv1.ExecuteSparkError{}, err)

		e := *err.(*sparkv1.ExecuteSparkError)
		assert.Equal(t, "My-Stage-1", e.StageName)
		assert.Equal(t, sparkv1.ErrorCode("ERR-1234"), e.ErrorCode)
		assert.Equal(t, "EEEK I broke", e.ErrorMessage)
		assert.Equal(t, map[string]any{"another": float64(1234), "foo": "bar"}, e.Metadata)
	})

	t.Run("Server Error: 502", func(t *testing.T) {
		svr := helpers.GetTestHttpServerWithRequests(t, []helpers.Request{
			{http.MethodGet, 200, "/basepath/spec", example1_Spec, nil, nil},
			{http.MethodPost, 200, "/basepath/init", nil, nil, nil},
			{http.MethodPost, 502, "/basepath/stages/My-Stage-1", []byte("dummy timeout issue"), nil, nil},
		})

		ctx := module_test_runner.NewTestJobContext(context.Background(), "test", "cid", "tid", module_test_runner.Inputs{})
		worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark(svr.URL+"/basepath"))
		assert.Nil(t, err)

		_, err = worker.Execute(ctx)
		if !assert.Error(t, err) {
			return
		}
		assert.IsType(t, &sparkv1.ExecuteSparkError{}, err)

		e := *err.(*sparkv1.ExecuteSparkError)
		assert.Equal(t, "My-Stage-1", e.StageName)
		assert.Equal(t, sparkv1.ErrorCodeGeneric, e.ErrorCode)
		assert.Equal(t, "unable to process complete stage response: http code 502: dummy timeout issue", e.ErrorMessage)
	})
}

func TestShouldErrorOnStageComplete(t *testing.T) {
	svr := helpers.GetTestHttpServerWithRequests(t, []helpers.Request{
		{http.MethodGet, 200, "/basepath/spec", example1_Spec, nil, nil},
		{http.MethodPost, 200, "/basepath/init", nil, nil, nil},
		{http.MethodPost, 200, "/basepath/stages/My-Stage-1", example1_s1_result, nil, nil},
		{http.MethodPost, 200, "/basepath/stages/My-Stage-2", example1_s2_result, nil, nil},
		{http.MethodPost, 500, "/basepath/complete/My-Complete", example1_stage_error_with_retry, nil, nil},
	})

	ctx := module_test_runner.NewTestJobContext(context.Background(), "test", "cid", "tid", module_test_runner.Inputs{})
	worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark(svr.URL+"/basepath"))
	assert.Nil(t, err)

	_, err = worker.Execute(ctx)
	if !assert.Error(t, err) {
		return
	}
	assert.IsType(t, &sparkv1.ExecuteSparkError{}, err)

	e := *err.(*sparkv1.ExecuteSparkError)
	assert.Equal(t, "main_complete", e.StageName)
	assert.Equal(t, sparkv1.ErrorCode("ERR-1234"), e.ErrorCode)
	assert.Equal(t, "EEEK I broke", e.ErrorMessage)
	assert.Equal(t, map[string]any{"another": float64(1234), "foo": "bar"}, e.Metadata)
}
