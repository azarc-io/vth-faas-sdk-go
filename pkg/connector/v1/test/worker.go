package test

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/mock"
	"testing"
)

type worker struct {
	t         *testing.T
	connector connectorv1.Connector
	startCtx  connectorv1.StartContext
	stopCtx   connectorv1.StopContext
	signaler  chan struct{}
}

func (w *worker) Run() {
	if w.startCtx != nil {
		err := w.connector.Start(w.startCtx)
		if err != nil {
			w.t.Error(err)
		}
	}
	if w.startCtx != nil && w.stopCtx != nil {
		<-w.signaler
	}
	if w.stopCtx != nil {
		err := w.connector.Stop(w.stopCtx)
		if err != nil {
			w.t.Error(err)
		}
	}
}

func NewTestConnectorWorker(t *testing.T, connector connectorv1.Connector, opts ...Option) (connectorv1.ConnectorWorker, chan struct{}) {
	signaler := make(chan struct{}, 1)
	w := &worker{t: t, connector: connector, signaler: signaler}
	for _, option := range opts {
		w = option(w)
	}
	return w, signaler
}

type Option func(*worker) *worker

func WithStartContextMock(ctx *mock.MockStartContext) Option {
	return func(w *worker) *worker {
		w.startCtx = ctx
		return w
	}
}
func WithStopContextMock(ctx *mock.MockStopContext) Option {
	return func(w *worker) *worker {
		w.stopCtx = ctx
		return w
	}
}
