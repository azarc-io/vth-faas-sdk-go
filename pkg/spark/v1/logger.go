package sparkv1

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/common"
)

type sparkContextLogger struct {
	*common.Logger
}

func (s *sparkContextLogger) AddFields(k string, v any) Logger {
	s.Logger.AddFields(k, v)
	return s
}

func NewLogger() Logger {
	logger, _ := common.NewLogger("spark_context", "")
	return &sparkContextLogger{
		logger,
	}
}
