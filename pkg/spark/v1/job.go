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
	Inputs                 ExecuteSparkInputs `json:"inputs"`   // all inputs for the given spark id
	JobKeyValue            string             `json:"job_key"`
	CorrelationIdValue     string             `json:"correlation_id"`
	TransactionIdValue     string             `json:"transaction_id"`
	RetryCount             uint               `yaml:"retry_count"`
	RetryBackoff           time.Duration      `yaml:"retry_backoff"`
	RetryBackoffMultiplier uint               `yaml:"retry_backoff_multiplier"`
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
