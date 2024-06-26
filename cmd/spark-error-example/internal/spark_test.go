package spark_test

import (
	"context"
	"testing"
	"time"

	spark "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-error-example/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1/test"
	"github.com/stretchr/testify/assert"
)

const totalTimesToErr = 10

func TestSparkErrorWithRetries(t *testing.T) {
	tests := []struct {
		name              string
		expectedOut       string
		expectedErr       string
		Times             uint
		FirstBackoffWait  time.Duration
		BackoffMultiplier uint
		Panic             bool
	}{
		{"recover after x retries", "finally I can pass after 10 failures", "", 15, 10 * time.Millisecond, 1, false},
		{"fail after x retries", "", "failures 5 of 10", 5, 15 * time.Millisecond, 2, false},
		{"fail gracefully after panic", "", "i was forced to panic :)", 10, 15 * time.Millisecond, 2, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx := module_test_runner.NewTestJobContext(context.Background(), "test", "cid", "tid", module_test_runner.Inputs{})
			worker, err := module_test_runner.NewTestRunner(t, spark.NewSpark(totalTimesToErr))
			assert.Nil(t, err)

			result, err := worker.ExecuteWithoutStageRetryOverride(ctx, sparkv1.WithSparkConfig(spark.Config{
				Panic: test.Panic,
				Retry: &sparkv1.RetryConfig{
					Times:             test.Times,
					FirstBackoffWait:  test.FirstBackoffWait,
					BackoffMultiplier: test.BackoffMultiplier,
				},
			}))

			if test.expectedErr == "" {
				// No error expected
				if !assert.Nil(t, err) {
					return
				}

				var message string
				assert.NoError(t, result.Bind("message", &message))
				assert.Equal(t, test.expectedOut, message)
				worker.AssertStageCompleted("chain-1_complete")
				return
			}

			//error expected
			assert.Contains(t, err.Error(), test.expectedErr)
		})
	}
}
