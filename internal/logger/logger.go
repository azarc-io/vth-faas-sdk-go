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

package logger

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type sparkContextLogger struct {
	log      zerolog.Logger
	metadata map[string]any
}

func (s *sparkContextLogger) Info(format string, v ...any) {
	log.Info().Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) Warn(format string, v ...any) {
	log.Warn().Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) Debug(format string, v ...any) {
	log.Debug().Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) Error(err error, format string, v ...any) {
	log.Error().Err(err).Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) AddFields(k string, v any) sdk_v1.Logger {
	s.metadata[k] = v
	return s
}

func NewLogger() sdk_v1.Logger {
	return &sparkContextLogger{
		metadata: map[string]any{},
		log:      log.With().Str("module", "spark_worker").CallerWithSkipFrameCount(3).Stack().Logger(),
	}
}
