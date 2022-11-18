package spark_v1

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type sparkContextLogger struct {
	log      zerolog.Logger
	metadata map[string]any
	sync.Mutex
}

func (s *sparkContextLogger) Info(format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Info().Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) Warn(format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Warn().Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) Debug(format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Debug().Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) Error(err error, format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Error().Err(err).Fields(s.metadata).Msgf(format, v...)
}

func (s *sparkContextLogger) AddFields(k string, v any) Logger {
	s.Lock()
	defer s.Unlock()
	s.metadata[k] = v
	return s
}

var skipFrameCount = 3

func NewLogger() Logger {
	return &sparkContextLogger{
		metadata: map[string]any{},
		log:      log.With().Str("module", "spark_worker").CallerWithSkipFrameCount(skipFrameCount).Stack().Logger(),
	}
}
