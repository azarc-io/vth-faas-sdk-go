package common

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type Logger struct {
	log      zerolog.Logger
	metadata map[string]any
	sync.Mutex
}

func (s *Logger) Info(format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Info().Fields(s.metadata).Msgf(format, v...)
}

func (s *Logger) Warn(format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Warn().Fields(s.metadata).Msgf(format, v...)
}

func (s *Logger) Debug(format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Debug().Fields(s.metadata).Msgf(format, v...)
}

func (s *Logger) Error(err error, format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Error().Err(err).Fields(s.metadata).Msgf(format, v...)
}

func (s *Logger) Fatal(err error, format string, v ...any) {
	s.Lock()
	defer s.Unlock()
	log.Fatal().Err(err).Fields(s.metadata).Msgf(format, v...)
}

func (s *Logger) AddFields(k string, v any) {
	s.Lock()
	defer s.Unlock()
	s.metadata[k] = v
}

var skipFrameCount = 3

func NewLogger(module string, level string) (*Logger, error) {
	lvl, err := zerolog.ParseLevel(level)
	if err != nil {
		return nil, err
	}
	return &Logger{
		metadata: map[string]any{},
		log:      log.With().Str("module", module).CallerWithSkipFrameCount(skipFrameCount).Stack().Logger().Level(lvl),
	}, nil
}
