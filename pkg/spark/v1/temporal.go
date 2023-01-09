package sparkv1

import (
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

var (
	DefaultRetryPolicy = &RetryPolicy{
		InitialInterval:        time.Second * 5,
		BackoffCoefficient:     2,
		MaximumInterval:        time.Minute * 5,
		MaximumAttempts:        1,
		NonRetryableErrorTypes: nil,
	}

	DefaultActivityOptions = &ActivityOptions{
		workflow.ActivityOptions{
			StartToCloseTimeout: time.Minute * 10,
		},
	}
)

type RetryPolicy struct {
	// Backoff interval for the first retry. If BackoffCoefficient is 1.0 then it is used for all retries.
	// If not set or set to 0, a default interval of 1s will be used.
	InitialInterval time.Duration `yaml:"initial_interval"`

	// Coefficient used to calculate the next retry backoff interval.
	// The next retry interval is previous interval multiplied by this coefficient.
	// Must be 1 or larger. Default is 2.0.
	BackoffCoefficient float64 `yaml:"backoff_coefficient"`

	// Maximum backoff interval between retries. Exponential backoff leads to interval increase.
	// This value is the cap of the interval. Default is 100x of initial interval.
	MaximumInterval time.Duration `yaml:"maximum_interval"`

	// Maximum number of attempts. When exceeded the retries stop even if not expired yet.
	// If not set or set to 0, it means unlimited, and rely on activity ScheduleToCloseTimeout to stop.
	MaximumAttempts int32 `yaml:"maximum_attempts"`

	// Non-Retriable errors. This is optional. Temporal server will stop retry if error type matches this list.
	// Note:
	//  - cancellation is not a failure, so it won't be retried,
	//  - only StartToClose or Heartbeat timeouts are retryable.
	NonRetryableErrorTypes []string `yaml:"non_retryable_error_types"`
}

func (rp RetryPolicy) GetTemporalPolicy() *temporal.RetryPolicy {
	return &temporal.RetryPolicy{
		InitialInterval:        rp.InitialInterval,
		BackoffCoefficient:     rp.BackoffCoefficient,
		MaximumInterval:        rp.MaximumInterval,
		MaximumAttempts:        rp.MaximumAttempts,
		NonRetryableErrorTypes: rp.NonRetryableErrorTypes,
	}
}

type ActivityOptions struct {
	workflow.ActivityOptions
}

func (ao ActivityOptions) GetTemporalActivityOptions() workflow.ActivityOptions {
	return ao.ActivityOptions
}

// Logger
type TemporalLogger struct {
}

func (f *TemporalLogger) Debug(msg string, keyvals ...interface{}) {
	log.Debug().Msgf(msg, keyvals...)
}

func (f *TemporalLogger) Info(msg string, keyvals ...interface{}) {
	log.Info().Msgf(msg, keyvals...)
}

func (f *TemporalLogger) Warn(msg string, keyvals ...interface{}) {
	log.Warn().Msgf(msg, keyvals...)
}

func (f *TemporalLogger) Error(msg string, keyvals ...interface{}) {
	log.Error().Msgf(msg, keyvals...)
}
