package main

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
)

/************************************************************************/
// TYPES
/************************************************************************/

type connector struct {
	config     *config
	broker     *mockBroker
	publishers map[string]*publication
}

type config struct {
	BrokerAddress string `json:"broker_address"`
}

type request struct {
	ctx     connectorv1.ForwardingContext
	path    string
	body    []byte
	headers map[string]interface{}
}

/************************************************************************/
// connectorv1.OutboundConnector IMPLEMENTATION
/************************************************************************/

func (c connector) HandleOutboundRequest(ctx connectorv1.OutboundRequest) (any, connectorv1.Headers, error) {
	var (
		requestBytes []byte
		err          error
		cfg          = c.publishers[ctx.MessageName()]
	)

	if requestBytes, err = ctx.Body().Raw(); err != nil {
		return nil, nil, err
	}

	err = c.broker.publish(cfg.topic, &message{
		body: requestBytes,
		// any headers to include, the agent will provide the following headers out of the box
		// - X-Request-Id
		// - X-Transaction-Id
		// - X-Correlation-Id
		// - X-Tenant-Id
		// - X-Workflow-Version
		// - X-Process-Id
		// - X-Timestamp
		// - Content-Type
		headers: ctx.Headers(),
	})

	return nil, nil, err
}

/************************************************************************/
// connectorv1.Connector IMPLEMENTATION
/************************************************************************/

// Start called by the sdk when the service has started successfully
// you can access custom configuration at this point and set up your clients/servers
// you can also read the message descriptors for inbound and outbound message types
// from the context
func (c connector) Start(ctx connectorv1.StartContext) error {
	// fetch user configured parameters for your connector, this is a json
	// payload that matches your configuration schema as set in the connector.yaml file
	if err := ctx.Config().Bind(&c.config); err != nil {
		return err
	}

	// create a client and connect to some external service
	// external address is provided through configuration
	c.broker = &mockBroker{address: c.config.BrokerAddress}
	// establish a connection, lets just pretend this is a http2 client
	if err := c.broker.connect(); err != nil {
		return err
	}

	// iterate over the inbound message descriptors and setup each one
	for _, descriptor := range ctx.InboundDescriptors() {
		var subCfg *subscription
		if err := descriptor.Config().Bind(&subCfg); err != nil {
			return err
		}

		if err := c.broker.Subscribe(subCfg, c.handleInboundRequest); err != nil {
			return err
		}
	}

	// iterate over the outbound message descriptors and setup each one
	// lets us cache any configuration so we don't have to bind it on every request
	for _, descriptor := range ctx.OutboundDescriptors() {
		var pubCfg *publication
		if err := descriptor.Config().Bind(&pubCfg); err != nil {
			return err
		}
		c.publishers[descriptor.MessageName()] = pubCfg
	}

	return nil
}

// Stop called by the sdk when the service is asked to shut down
// you can gracefully terminate any clients/servers at this point
func (c connector) Stop(_ connectorv1.StopContext) error {
	return c.broker.disconnect()
}

/************************************************************************/
// INBOUND HANDLING
/************************************************************************/

// handleInboundRequest handles inbound requests from the server e.g. open api server
func (c connector) handleInboundRequest(req *request) (string, []byte, connectorv1.Headers, error) {
	response, err := req.ctx.Forward(req.path, req.body, req.headers)
	if err != nil {
		req.ctx.LogError(err, "could not handle inbound request")
		return "", nil, nil, err
	}

	rawBody, err := response.Body().Raw()
	if err != nil {
		req.ctx.LogError(err, "could not fetch response from agent")
		return "", nil, nil, err
	}

	return response.MessageName(), rawBody, response.Headers(), nil
}

/************************************************************************/
// ENTRY POINT
/************************************************************************/

func main() {
	service := connectorv1.New(&connector{publishers: map[string]*publication{}})
	if err := service.Start(); err != nil {
		panic(err)
	}
}
