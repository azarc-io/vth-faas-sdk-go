package main

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/mock"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/test"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnector(t *testing.T) {
	ctrl := gomock.NewController(t)
	startCtx := mock.NewMockStartContext(ctrl)
	c := &connector{}

	loggerMock := mock.NewMockLogger(ctrl)
	loggerMock.EXPECT().Error(gomock.Any(), gomock.Any()).AnyTimes()
	startCtx.EXPECT().Log().Return(loggerMock)

	configMock := mock.NewMockBindable(ctrl)
	configMock.EXPECT().Bind(gomock.Any()).Do(func(val interface{}) error {
		conf, ok := val.(*config)
		assert.True(t, ok)
		conf.ClientOpenApiSpec = "client_spec"
		conf.ServerOpenApiSpec = "server_spec"
		conf.OutboundAddress = "outbound_address"
		return nil
	}).Return(nil)
	startCtx.EXPECT().Config().Return(configMock)

	ingressMock := mock.NewMockIngress(ctrl)
	ingressMock.EXPECT().InternalHost().Return("ingress.internal.host")
	ingressMock.EXPECT().InternalPort().Return(1234)
	startCtx.EXPECT().Ingress("http-8080").Return(ingressMock, nil)

	forwarderMock := mock.NewMockForwarder(ctrl)
	inboundResponseMock := mock.NewMockInboundResponse(ctrl)
	responseBodyMock := mock.NewMockBindable(ctrl)
	responseBodyMock.EXPECT().Raw().Return([]byte("response-body"), nil)
	inboundResponseMock.EXPECT().Body().Return(responseBodyMock)
	inboundResponseMock.EXPECT().Headers().Return(connectorv1.Headers{
		"response-header-key": "response-header-value",
	})
	forwarderMock.EXPECT().Forward("test-message", []byte("test-body"), connectorv1.Headers{
		"request-header-key": "request-header-value",
	}).Return(inboundResponseMock, nil)
	startCtx.EXPECT().Forwarder().Return(forwarderMock)

	worker, _ := test.NewTestConnectorWorker(t, c, test.WithStartContextMock(startCtx))
	worker.Run() // run is not blocking because we test only Start method

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
}
