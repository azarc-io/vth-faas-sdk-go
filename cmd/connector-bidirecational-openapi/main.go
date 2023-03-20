package main

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
)

/************************************************************************/
// TYPES
/************************************************************************/

type connector struct {
	config *config
	client *mockClient
	server *mockServer
}

type config struct {
	ClientOpenApiSpec string `json:"client_open_api_spec"`
	ServerOpenApiSpec string `json:"server_open_api_spec"`
	OutboundAddress   string `json:"outbound_address"`
}

type request struct {
	forwarder connectorv1.Forwarder
	path      string
	body      []byte
	headers   connectorv1.Headers
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

	responseBytes, headers, err = c.client.DoExternalRequest(
		// this is the message name, in this demo the message name is not editable by the user
		// because it is parsed from the open api spec, thus the message name = open api endpoint already
		ctx.MessageName(),
		// this is the mimetype as defined in the message spec, in this example we are using open api
		// so the in the ui the messages have already been broken up by mimetype
		ctx.MimeType(),
		// the body of the request
		requestBytes,
		// any headers to include, the agent will provide the following headers out of the box
		// - X-Request-Id
		// - X-Transaction-Id
		// - X-Correlation-Id
		// - X-Tenant-Id
		// - X-Workflow-Version
		// - X-Process-Id
		// - X-Timestamp
		// - Content-Type
		ctx.Headers(),
	)

	return responseBytes, headers, err
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
	c.client = &mockClient{address: c.config.OutboundAddress, spec: c.config.ClientOpenApiSpec}
	// establish a connection, lets just pretend this is a http2 client
	if err := c.client.connect(); err != nil {
		return err
	}

	// request the ingress configuration from the context, if ingress is
	// not enabled in the connector.yaml then this will error
	ingress, err := ctx.Ingress("http-8080")
	if err != nil {
		return err
	}

	// create a server and expose it on the given ip and port
	// ingress host will either be the ip address of the service or 0.0.0.0
	// ingress port will use the configured port from the connector.yaml when unit testing locally
	// ingress port will be provided by the agent when deployed through verathread
	c.server = &mockServer{bindHost: ingress.InternalHost(), bindPort: ingress.InternalPort(), spec: c.config.ServerOpenApiSpec}
	// register a handler with our mock server, you have to wrap the handler so that you can
	// pass a forwarding context to your actual handler, that will give you access to everything
	// you need to handle an inbound request
	c.server.onRequest = func(path string, body []byte, headers connectorv1.Headers) (*response, error) {
		rPath, rBody, rHeaders, err := c.handleInboundRequest(&request{
			forwarder: ctx.Forwarder(),
			path:      path,
			body:      body,
			headers:   headers,
		}, ctx.Log())

		if err != nil {
			return nil, err
		}

		return &response{body: rBody, headers: rHeaders, path: rPath}, nil
	}
	// start the server
	go func() {
		if err := c.server.start(); err != nil {
			panic(err)
		}
	}()

	return nil
}

// Stop called by the sdk when the service is asked to shut down
// you can gracefully terminate any clients/servers at this point
func (c connector) Stop(ctx connectorv1.StopContext) error {
	if err := c.server.stop(); err != nil {
		ctx.Log().Error(err, "failed to gracefully stop the server")
	}
	return c.client.disconnect()
}

/************************************************************************/
// INBOUND HANDLING
/************************************************************************/

// handleInboundRequest handles inbound requests from the server e.g. open api server
func (c connector) handleInboundRequest(req *request, logger connectorv1.Logger) (string, []byte, connectorv1.Headers, error) {
	response, err := req.forwarder.Forward(req.path, req.body, req.headers)
	if err != nil {
		logger.Error(err, "could not handle inbound request")
		return "", nil, nil, err
	}

	rawBody, err := response.Body().Raw()
	if err != nil {
		logger.Error(err, "could not fetch response from agent")
		return "", nil, nil, err
	}

	return req.path, rawBody, response.Headers(), nil
}

/************************************************************************/
// ENTRY POINT
/************************************************************************/

func main() {
	service, err := connectorv1.NewConnectorWorker(&connector{})
	if err != nil {
		panic(err)
	}
	service.Run()
}
