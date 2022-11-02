// Copyright 2020-2022 Azarc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package context

import (
	"context"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type JobMetadata struct {
	ctx             context.Context
	jobKey          string
	correlationId   string
	transactionId   string
	lastActiveStage *sdk_v1.LastActiveStage
}

func NewJobMetadata(ctx context.Context, jobKey, correlationId, transactionId string, lastActiveStage *sdk_v1.LastActiveStage) JobMetadata {
	return JobMetadata{ctx: ctx, jobKey: jobKey, correlationId: correlationId, transactionId: transactionId, lastActiveStage: lastActiveStage}
}

func NewJobMetadataFromGrpcRequest(ctx context.Context, req *sdk_v1.ExecuteJobRequest) JobMetadata {
	return JobMetadata{
		ctx:             ctx,
		jobKey:          req.Key,
		correlationId:   req.CorrelationId,
		transactionId:   req.TransactionId,
		lastActiveStage: req.LastActiveStage,
	}
}

func (j JobMetadata) JobKey() string {
	return j.jobKey
}

func (j JobMetadata) CorrelationID() string {
	return j.correlationId
}

func (j JobMetadata) TransactionID() string {
	return j.transactionId
}

func (j JobMetadata) Ctx() context.Context {
	return j.ctx
}

func (j JobMetadata) LastActiveStage() *sdk_v1.LastActiveStage {
	return j.lastActiveStage
}
