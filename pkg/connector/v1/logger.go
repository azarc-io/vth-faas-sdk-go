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
