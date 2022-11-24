package connectorv1

import "context"

type connector struct {
}

func (c connector) Start() error {
	//TODO implement me
	panic("implement me")
}

func New(ctx context.Context, con Connector) ConnectorService {
	c := &connector{}

	return c
}
