package sparkv1

import (
	"context"
	"time"
)

//********************************************************************************************
// INBOUND JOB METADATA
//********************************************************************************************

// JobMetadata the context for the spark we want to execute on a module
// TODO this type should come from the Module Library
type JobMetadata struct {
	SparkId                string             `json:"spark_id"` // id of the spark to execute
	JobKeyValue            string             `json:"job_key"`
	CorrelationIdValue     string             `json:"correlation_id"`
	TransactionIdValue     string             `json:"transaction_id"`
	RetryCount             uint               `json:"retry_count"`
	RetryBackoff           time.Duration      `json:"retry_backoff"`
	RetryBackoffMultiplier uint               `json:"retry_backoff_multiplier"`
	JobPid                 *JobPid            `json:"job_pid,omitempty"`
	VariablesBucket        string             `json:"variables_bucket"`
	VariablesKey           string             `json:"variables_key"`
	Model                  string             `json:"model,omitempty"`
	Inputs                 ExecuteSparkInputs `json:"-"`
}

type JobPid struct {
	Address   string `json:"Address"`
	Id        string `json:"Id"`
	RequestId uint32 `json:"request_id"`
}

type JobContext struct {
	context.Context
	Metadata *JobMetadata
}

func (jc *JobContext) JobKey() string {
	return jc.Metadata.JobKeyValue
}

func (jc *JobContext) CorrelationID() string {
	return jc.Metadata.CorrelationIdValue
}

func (jc *JobContext) TransactionID() string {
	return jc.Metadata.TransactionIdValue
}
