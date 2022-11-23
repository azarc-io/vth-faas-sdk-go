package connectorv1

type connector struct {
}

func (c connector) Start() error {
	//TODO implement me
	panic("implement me")
}

func New(con Connector) ConnectorService {
	c := &connector{}

	return c
}
