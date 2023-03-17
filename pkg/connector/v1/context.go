package connectorv1

import (
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/internal/healthz"
)

/************************************************************************/
// START CONTEXT
/************************************************************************/

type startContext struct {
	userConfig         Bindable
	inboundDescriptors []messageDescriptor
	logger             Logger
	forwarder          Forwarder
	health             healthChecker
	healthConfig       *configHealth
	ingress            []ingressConfig
}

func (c *startContext) Ingress(name string) (Ingress, error) {
	for _, ing := range c.ingress {
		if ing.Name == name {
			return &ing, nil
		}
	}
	return nil, errors.New("ingress not found")
}

func (c *startContext) InboundDescriptors() []InboundDescriptor {
	descriptors := make([]InboundDescriptor, len(c.inboundDescriptors))
	for i := range c.inboundDescriptors {
		descriptors[i] = c.inboundDescriptors[i]
	}
	return descriptors
}

func (c *startContext) OutboundDescriptors() []OutboundDescriptor {
	//TODO implement me
	return nil
}

func (c *startContext) Forwarder() Forwarder {
	return c.forwarder
}

func (c *startContext) Log() Logger {
	return c.logger
}

func (c *startContext) RegisterPeriodicHealthCheck(name string, fn HealthCheckFunc) {
	c.health.Register(name, c.healthConfig.Interval, healthz.CheckFunc(fn))
}

func (c *startContext) Config() Bindable {
	return c.userConfig
}

/************************************************************************/
// STOP CONTEXT
/************************************************************************/

type stopContext struct {
	logger Logger
}

func (c *stopContext) Log() Logger {
	return c.logger
}
