package connectorv1

import (
	"context"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/internal/healthz"
	"github.com/azarc-io/vth-faas-sdk-go/internal/signals"
	"net/http"
	"time"
)

const runtimeTTL = time.Minute

type worker struct {
	connector          Connector
	opts               *ConnectorOpts
	config             *connectorConfig
	ingress            []ingressConfig
	userConfig         Bindable
	inboundDescriptors []messageDescriptor

	health       healthChecker
	healthServer *http.Server
}

type healthChecker interface {
	Register(name string, period time.Duration, fn healthz.CheckFunc)
	Handler() http.Handler
}

func (w *worker) Run() {

	if w.healthServer != nil {
		go func() {
			if err := w.healthServer.ListenAndServe(); err != http.ErrServerClosed {
				panic(err)
			}
		}()
		defer func() {
			if err := w.healthServer.Shutdown(context.Background()); err != nil {
				w.opts.log.Error(err, "failed to shutdown health server")
			}
		}()
	}

	startCtx := startContext{
		userConfig:         w.userConfig,
		inboundDescriptors: w.inboundDescriptors,
		logger:             w.opts.log,
		forwarder:          w.opts.forwarder,
		health:             w.health,
		healthConfig:       w.config.Health,
		ingress:            w.ingress,
	}

	err := w.connector.Start(&startCtx)
	if err != nil {
		panic(err)
	}

	// wait for signal to shut down
	<-signals.SetupSignalHandler()

	stopCtx := stopContext{logger: w.opts.log}
	err = w.connector.Stop(&stopCtx)
	if err != nil {
		panic(err)
	}

}

func (w *worker) loadConfiguration() error {
	c, err := loadConnectorConfig(w.opts)
	if err != nil {
		return err
	}
	w.config = &c.ConnectorConfig
	w.ingress = c.Ingress

	inboundDescriptors, err := loadMessageDescriptorsConfig(MessageTypeInbound)
	if err != nil {
		return err
	}
	w.inboundDescriptors = inboundDescriptors

	userConfig, err := loadUserConfig(w.opts)
	if err != nil {
		return err
	}
	w.userConfig = userConfig

	return nil
}

func (w *worker) initHealthz() {
	if w.config.Health != nil && w.config.Health.Enabled {
		w.health = healthz.NewChecker(&healthz.Config{
			RuntimeTTL: runtimeTTL,
		})

		http.Handle("/healthz", w.health.Handler())
		w.healthServer = &http.Server{Addr: fmt.Sprintf("%s:%d", w.config.Health.Bind, w.config.Health.Port)}
	}
}

func NewConnectorWorker(connector Connector, options ...Option) (ConnectorWorker, error) {
	w := worker{
		connector: connector,
		opts:      &ConnectorOpts{},
	}
	for _, opt := range options {
		w.opts = opt(w.opts)
	}

	err := w.loadConfiguration()
	if err != nil {
		return nil, err
	}

	if w.opts.log == nil {
		w.opts.log, err = newLogger(w.config.Log)
		if err != nil {
			return nil, err
		}
	}

	if w.opts.forwarder == nil {
		w.opts.forwarder = newForwarder(w.config)
	}

	w.initHealthz()

	return &w, nil
}
