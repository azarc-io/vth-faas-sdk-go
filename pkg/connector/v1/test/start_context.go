package test

import (
	"errors"
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

type startContext struct {
	ctrl               *gomock.Controller
	userConfig         []byte
	inboundDescriptors []InboundDescriptor
	logger             connectorv1.Logger
	forwarder          connectorv1.Forwarder
	ingress            []IngressConfig
	healthCheckers     map[string]connectorv1.HealthCheckFunc

	LoggerMock    *mock.MockLogger
	ForwarderMock *mock.MockForwarder
}

func (c *startContext) Ingress(name string) (connectorv1.Ingress, error) {
	for _, ing := range c.ingress {
		if ing.Name == name {
			return &ing, nil
		}
	}
	return nil, errors.New("ingress not found")
}

func (c *startContext) InboundDescriptors() []connectorv1.InboundDescriptor {
	descriptors := make([]connectorv1.InboundDescriptor, len(c.inboundDescriptors))
	for i := range c.inboundDescriptors {
		descriptors[i] = c.inboundDescriptors[i]
	}
	return descriptors
}

func (c *startContext) OutboundDescriptors() []connectorv1.OutboundDescriptor {
	return nil
}

func (c *startContext) Forwarder() connectorv1.Forwarder {
	return c.forwarder
}

func (c *startContext) Log() connectorv1.Logger {
	return c.logger
}

func (c *startContext) RegisterPeriodicHealthCheck(name string, fn connectorv1.HealthCheckFunc) {
	c.healthCheckers[name] = fn
}

func (c *startContext) Config() connectorv1.Bindable {
	return connectorv1.NewBindable(c.userConfig, connectorv1.BindableTypeJson)
}

func NewStartContext(t *testing.T, config *Config) *startContext {
	ctrl := gomock.NewController(t)
	forwarderMock := mock.NewMockForwarder(ctrl)

	return &startContext{
		ctrl:               ctrl,
		userConfig:         config.UserConfig,
		inboundDescriptors: config.InboundDescriptors,
		logger:             noopLogger{},
		forwarder:          forwarderMock,
		ingress:            config.Ingress,
		healthCheckers:     make(map[string]connectorv1.HealthCheckFunc),

		ForwarderMock: forwarderMock,
	}
}

func (c *startContext) WithLoggerMock() *startContext {
	c.LoggerMock = mock.NewMockLogger(c.ctrl)
	c.logger = c.LoggerMock
	return c
}

type InboundResponse struct {
	HeadersMap connectorv1.Headers
	Payload    []byte
}

func (f *InboundResponse) Body() connectorv1.Bindable {
	return connectorv1.NewBindable(f.Payload, connectorv1.BindableTypeJson)
}

func (f *InboundResponse) Headers() connectorv1.Headers {
	return f.HeadersMap
}

func (c *startContext) MockForward(messageName string, body any, headers any, response *InboundResponse, responseErr error) {
	c.ForwarderMock.EXPECT().Forward(messageName, body, headers).Return(response, responseErr)
}
