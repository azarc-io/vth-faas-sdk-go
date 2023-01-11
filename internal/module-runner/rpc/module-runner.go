//go:build exclude

package module_runner

import (
	"context"
	"encoding/base64"
	"github.com/azarc-io/vth-faas-sdk-go/internal/common"
	"github.com/azarc-io/vth-faas-sdk-go/internal/healthz"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v3"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

// TODO: Taken copy for tomorrows jono and wael to remember

type Runner interface {
	Stop() error
}

type runner struct {
	sparks map[string]*sparkClient
	health *healthz.Checker
}

func (r runner) Stop() error {
	for _, s := range r.sparks {
		if !s.pluginClient.Exited() {
			s.pluginClient.Kill()
		}
	}

	return nil
}

func (r runner) initHealthz(cfg *config) {
	// TODO support TLS once support for platforms other than kubernetes are added to Verathread
	if cfg.Health != nil && cfg.Health.Enabled {
		r.health = healthz.NewChecker(&healthz.Config{
			RuntimeTTL: time.Second * 5,
		})

		go func() {
			http.Handle("/healthz", r.health.Handler())

			// nosemgrep
			if err := http.ListenAndServe(cfg.healthBindTo(), nil); err != nil { // nosemgrep
				panic(err)
			}
		}()
	}
}

type sparkClient struct {
	id           string
	name         string
	pluginClient *plugin.Client
	rpcClient    plugin.ClientProtocol
	Workflow     *JobWorkflow
}

func RunModule(ctx context.Context, cfg *config) (Runner, error) {
	runner := runner{sparks: make(map[string]*sparkClient)}

	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	for _, s := range cfg.Sparks {
		m, _ := yaml.Marshal(map[string]any{
			"config": map[string]string{
				"id":   s.Id,
				"name": s.Name,
			},
			"temporal": cfg.Temporal,
		})
		cmd := exec.Command(path.Join(cfg.BinBasePath, s.Name))
		cmd.Env = os.Environ()
		cmd.Env = append(cmd.Env, "SPARK_SECRET="+base64.StdEncoding.EncodeToString(m))

		if s.Config != "" {
			cmd.Env = append(cmd.Env, "CONFIG_SECRET="+base64.StdEncoding.EncodeToString([]byte(s.Config)))
		}

		// We're a host! Start by launching the plugin process.
		pc := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "BASIC_PLUGIN",
				MagicCookieValue: s.Id,
			},
			Plugins: pluginMap,
			Cmd:     cmd,
			Logger:  logger,
		})

		sparkId := s.Id
		go func() {
			// Connect via RPC
			rpcClient, err := pc.Client()
			if err != nil {
				log.Fatal(err)
			}

			// Request the plugin
			raw, err := rpcClient.Dispense("spark")
			if err != nil {
				log.Fatal(err)
			}

			sa := raw.(sparkv1.SparkRpcApi)
			log.Printf("pong: %s", sa.Greet(&sparkv1.IBlackboard{
				Value: "module runner",
				GetVal: func() string {
					return "its from here"
				},
			}))

			runner.sparks[sparkId] = &sparkClient{
				id:           sparkId,
				name:         s.Name,
				pluginClient: pc,
				rpcClient:    rpcClient,
				Workflow: &JobWorkflow{
					SparkId: sparkId,
					Client:  rpcClient,
					Chain:   &common.SparkChain{},
				},
			}
		}()

		//// Request the plugin
		//raw, err := rpcClient.Dispense("greeter")
		//if err != nil {
		//	log.Fatal(err)
		//}

		runner.initHealthz(cfg)
	}

	return &runner, nil
}

// pluginMap is the map of plugins we can dispense.
var pluginMap = map[string]plugin.Plugin{
	"spark": &sparkv1.SparkPlugin{},
}
