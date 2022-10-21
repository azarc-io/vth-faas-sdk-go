package logger

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
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

func (s *sparkContextLogger) AddFields(k string, v any) api.Logger {
	s.metadata[k] = v
	return s
}

func NewLogger() api.Logger {
	return &sparkContextLogger{
		metadata: map[string]any{},
		log:      log.With().Str("module", "spark_worker").CallerWithSkipFrameCount(3).Stack().Logger(),
	}
}
