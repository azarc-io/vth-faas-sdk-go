package connectorv1

import (
	"encoding/json"
	"github.com/rs/zerolog/log"
)

type ConnectorOpts struct {
	log            Logger
	forwarder      Forwarder
	config         []byte
	configBasePath string
}

type Option = func(je *ConnectorOpts) *ConnectorOpts

func WithConfig(cfg any) Option {
	d, err := json.Marshal(cfg)
	if err != nil {
		log.Fatal().Err(err).Msgf("unable to serialise config")
	}
	return func(opts *ConnectorOpts) *ConnectorOpts {
		opts.config = d
		return opts
	}
}

func WithConfigBasePath(path string) Option {
	return func(opts *ConnectorOpts) *ConnectorOpts {
		opts.configBasePath = path
		return opts
	}
}

func WithForwarder(fwd Forwarder) Option {
	return func(opts *ConnectorOpts) *ConnectorOpts {
		opts.forwarder = fwd
		return opts
	}
}
