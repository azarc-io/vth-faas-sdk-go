package main

import connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"

/************************************************************************/
// TYPES
/************************************************************************/

type connector struct {
	client mockClient
	config *config
}

type config struct {
}

/************************************************************************/
// connectorv1.OutboundConnector IMPLEMENTATION
/************************************************************************/

func (c connector) HandleOutboundRequest(ctx connectorv1.OutboundRequest) (any, connectorv1.Headers, error) {
	var (
		requestBytes, responseBytes []byte
		headers                     connectorv1.Headers
		err                         error
	)

	if requestBytes, err = ctx.Body().Raw(); err != nil {
		return nil, nil, err
	}

	responseBytes, headers, err = c.client.DoExternalRequest(requestBytes, ctx.Headers())

	return responseBytes, headers, err
}

/************************************************************************/
// connectorv1.Connector IMPLEMENTATION
/************************************************************************/

func (c connector) Start(ctx connectorv1.StartContext) error {
	if err := ctx.Config().Bind(&c.config); err != nil {
		return err
	}

	// create a client and connect to some external service
	c.client = mockClient{address: c.config}

	return c.client.Connect()
}

func (c connector) Stop() error {
	return c.client.Disconnect()
}

func newConnector() connectorv1.OutboundConnector {
	return &connector{}
}

func main() {
	service := connectorv1.New(newConnector())
	if err := service.Start(); err != nil {
		panic(err)
	}
}
