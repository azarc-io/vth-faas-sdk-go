package test

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/mock"
	"github.com/golang/mock/gomock"
	"testing"
)

type stopContext struct {
	ctrl   *gomock.Controller
	logger connectorv1.Logger

	LoggerMock *mock.MockLogger
}

func (c *stopContext) Log() connectorv1.Logger {
	return c.logger
}

func NewStopContext(t *testing.T) *stopContext {
	ctrl := gomock.NewController(t)

	return &stopContext{
		ctrl:   ctrl,
		logger: noopLogger{},
	}
}

func (c *stopContext) WithLoggerMock() *stopContext {
	c.LoggerMock = mock.NewMockLogger(c.ctrl)
	c.logger = c.LoggerMock
	return c
}
