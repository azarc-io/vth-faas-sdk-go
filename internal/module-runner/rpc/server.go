//go:build exclude

package spark_v1

import (
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"os"
)

/************************************************************************/
// TYPES
/************************************************************************/

type sparkPlugin struct {
	config *config
	worker Worker
}

/************************************************************************/
// SERVER
/************************************************************************/

func newSparkPlugin(cfg *config, worker Worker) *sparkPlugin {
	return &sparkPlugin{config: cfg, worker: worker}
}

func (s *sparkPlugin) start() error {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"spark": &SparkPlugin{Impl: s},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   "BASIC_PLUGIN",
			MagicCookieValue: s.config.Config.Id,
		},
		Plugins: pluginMap,
	})

	return nil
}

func (s *sparkPlugin) Greet(bb Blackboard) string {
	bb.SetValue("howdy back: " + s.config.Config.Id)
	return "Hello: " + s.config.Config.Id
}

func (s *sparkPlugin) stop() {
}
