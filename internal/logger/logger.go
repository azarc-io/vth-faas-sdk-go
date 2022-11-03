package logger

import (
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
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

var skipFrameCount = 3

func NewLogger() sdk_v1.Logger {
	return &sparkContextLogger{
		metadata: map[string]any{},
		log:      log.With().Str("module", "spark_worker").CallerWithSkipFrameCount(skipFrameCount).Stack().Logger(),
	}
}
