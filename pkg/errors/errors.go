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

package errors

import (
	"encoding/json"
	"errors"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"google.golang.org/protobuf/types/known/structpb"
)

var (
	StageDoesNotExist = errors.New("stage does not exists")
	BindValueFailed   = errors.New("bind value failed")
	VariableNotFound  = errors.New("variable not found")

	errorTypeToStageStatusMapper = map[sdk_v1.ErrorType]sdk_v1.StageStatus{
		sdk_v1.ErrorType_Retry:    sdk_v1.StageStatus_StageFailed,
		sdk_v1.ErrorType_Skip:     sdk_v1.StageStatus_StageSkipped,
		sdk_v1.ErrorType_Canceled: sdk_v1.StageStatus_StageCanceled,
		sdk_v1.ErrorType_Failed:   sdk_v1.StageStatus_StageFailed,
	}
)

func ErrorTypeToStageStatusMapper(errType sdk_v1.ErrorType) sdk_v1.StageStatus {
	if err, ok := errorTypeToStageStatusMapper[errType]; ok {
		return err
	}
	return sdk_v1.StageStatus_StageFailed
}

type Option = func(err *Stage) *Stage

type Stage struct {
	err       error
	errorType sdk_v1.ErrorType
	errorCode uint32
	metadata  map[string]any
	retry     *RetryConfig
}

type RetryConfig struct {
	times         uint
	backoffMillis uint
}

func NewStageError(err error, opts ...Option) *Stage {
	stg := &Stage{err: err}
	for _, opt := range opts {
		stg = opt(stg)
	}
	return stg
}

func (s *Stage) ErrorType() sdk_v1.ErrorType {
	return s.errorType
}

func (s *Stage) Code() uint32 {
	return s.errorCode
}

func (s *Stage) Error() string {
	return s.err.Error()
}

func (s *Stage) Metadata() map[string]any {
	return s.metadata
}

func (s *Stage) ToErrorMessage() *sdk_v1.Error {
	err := &sdk_v1.Error{
		Error:     s.err.Error(),
		ErrorCode: s.errorCode,
		ErrorType: s.errorType,
	}
	if s.metadata != nil {
		err.Metadata, _ = structpb.NewValue(s.metadata)
	}
	if s.retry != nil {
		err.Retry = &sdk_v1.RetryStrategy{Backoff: uint32(s.retry.backoffMillis), Count: uint32(s.retry.times)}
	}
	return err
}

func WithErrorType(errorType sdk_v1.ErrorType) Option {
	return func(err *Stage) *Stage {
		err.errorType = errorType
		return err
	}
}

func WithErrorCode(code uint32) Option {
	return func(err *Stage) *Stage {
		err.errorCode = code
		return err
	}
}

func WithMetadata(metadata any) Option {
	return func(err *Stage) *Stage {
		err.parseMetadata(metadata)
		return err
	}
}

func WithRetry(times, backoffMillis uint) Option {
	return func(err *Stage) *Stage {
		err.retry = &RetryConfig{times, backoffMillis}
		err.errorType = sdk_v1.ErrorType_Retry
		return err
	}
}

func (s *Stage) parseMetadata(metadata any) {
	m := map[string]any{}
	if metadata != nil {
		mdBytes, _ := json.Marshal(metadata)
		_ = json.Unmarshal(mdBytes, &m)
	}
	s.metadata = m
}
