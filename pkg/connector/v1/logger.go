package connectorv1

import "github.com/azarc-io/vth-faas-sdk-go/internal/common"

type connectorContextLogger struct {
	*common.Logger
}

func newLogger(cfg *configLog) (Logger, error) {
	logger, err := common.NewLogger("connector_context", cfg.Level)
	if err != nil {
		return nil, err
	}
	return &connectorContextLogger{
		logger,
	}, nil
}

type noopLogger struct {
}

func (n noopLogger) Error(_ error, _ string, _ ...interface{}) {
}

func (n noopLogger) Fatal(_ error, _ string, _ ...interface{}) {
}

func (n noopLogger) Info(_ string, _ ...interface{}) {
}

func (n noopLogger) Warn(_ string, _ ...interface{}) {
}

func (n noopLogger) Debug(_ string, _ ...interface{}) {
}
