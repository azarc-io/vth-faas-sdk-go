package main

import (
	_ "embed"
	"encoding/json"
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

//go:embed connector_test_config.json
var configBytes []byte

func TestConnector(t *testing.T) {
	// initialize connector
	c := &connector{}

	// load config
	var config test.Config
	err := json.Unmarshal(configBytes, &config)
	assert.NoError(t, err)

	// create a start context
	ctx := test.NewStartContext(t, &config)

	// setup expectations
	ctx.MockForward("test-message", []byte("test-body"), connectorv1.Headers{
		"request-header-key": "request-header-value",
	}, &test.InboundResponse{
		HeadersMap: connectorv1.Headers{
			"response-header-key": "response-header-value",
		},
		Payload: []byte("response-body"),
	}, nil)

	// call Start method
	err = c.Start(ctx)

	// verify results
	assert.NoError(t, err)

	assert.Equal(t, "client_spec", c.config.ClientOpenApiSpec)
	assert.Equal(t, "server_spec", c.config.ServerOpenApiSpec)
	assert.Equal(t, "outbound_address", c.config.OutboundAddress)

	assert.NotNil(t, c.client)
	assert.True(t, c.client.connected)

	assert.NotNil(t, c.server)
	assert.Equal(t, "ingress.internal.host", c.server.bindHost)
	assert.Equal(t, 1234, c.server.bindPort)
	assert.Equal(t, "server_spec", c.server.spec)

	assert.NotNil(t, c.server.onRequest)
	resp, err := c.server.onRequest("test-message", []byte("test-body"), connectorv1.Headers{
		"request-header-key": "request-header-value",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "test-message", resp.path)
	assert.Equal(t, []byte("response-body"), resp.body)
	assert.Equal(t, connectorv1.Headers{
		"response-header-key": "response-header-value",
	}, resp.headers)

	// create a stop context
	stopCtx := test.NewStopContext(t).WithLoggerMock()

	// setup expectations
	stopCtx.LoggerMock.EXPECT().Info("stopping")

	// call Stop method
	err = c.Stop(stopCtx)

	// verify results
	assert.NoError(t, err)
	assert.False(t, c.client.connected)
}
