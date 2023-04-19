package module_runner

import (
	"encoding/base64"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/internal/healthz"
	hclog "github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"gopkg.in/yaml.v3"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"time"
)

const (
	requestTokenHeader = "X-Token"
)

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
}

func RunModule(cfg *config) (Runner, error) {
	r := runner{sparks: make(map[string]*sparkClient)}

	// Create an hclog.Logger
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Debug,
	})

	for _, s := range cfg.Sparks {
		cmd := exec.Command(path.Join(cfg.BinBasePath, s.Name))
		cmd.Env = os.Environ()

		cfgPath := path.Join(os.TempDir(), fmt.Sprintf("%s-%s-config.json", s.Name, s.Id))

		var cfgData []byte
		// Check if config server is used
		if s.ConfigServer != nil {
			var err error
			cfgData, err = getConfigFromServer(s)
			if err != nil {
				return nil, err
			}
		} else {
			// Deprecated: Move to using config server
			cfgData = []byte(s.Config)
		}

		if err := os.WriteFile(cfgPath, cfgData, fs.ModePerm); err != nil {
			return nil, err
		}

		cmd.Env = append(cmd.Env, fmt.Sprintf("CONFIG_FILE_PATH=%s", cfgPath))

		// Create the config options for the spark runner
		m, _ := yaml.Marshal(map[string]any{
			"id":          s.Id,
			"name":        s.Name,
			"queue_group": s.QueueGroup,
			"temporal":    cfg.Temporal,
		})

		cmd.Env = append(cmd.Env, "SPARK_SECRET="+base64.StdEncoding.EncodeToString(m))

		startupTimeout := time.Second * 20
		if s.StartupTimeout != nil {
			startupTimeout = *s.StartupTimeout
		}

		// We're a host! Start by launching the plugin process.
		pc := plugin.NewClient(&plugin.ClientConfig{
			HandshakeConfig: plugin.HandshakeConfig{
				ProtocolVersion:  1,
				MagicCookieKey:   "BASIC_PLUGIN",
				MagicCookieValue: s.Id,
			},
			Plugins:      map[string]plugin.Plugin{},
			Cmd:          cmd,
			Logger:       logger,
			StartTimeout: startupTimeout,
			//TODO: Investigate graceful shutdown time, currently defaults to 2s:
			//  https://github.com/hashicorp/go-plugin/pull/222/files
		})

		sparkId := s.Id
		go func() {
			// Connect via RPC
			rpcClient, err := pc.Client()
			if err != nil {
				log.Fatal(err)
			}

			r.sparks[sparkId] = &sparkClient{
				id:           sparkId,
				name:         s.Name,
				pluginClient: pc,
				rpcClient:    rpcClient,
			}
		}()
	}

	r.initHealthz(cfg)

	return &r, nil
}

func getConfigFromServer(s *configSpark) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, s.ConfigServer.Url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set(requestTokenHeader, s.ConfigServer.ApiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("unable to request spark config: (%s): %w", s.Id, err)
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unable to fetch spark config: (%s): %w", s.Id, err)
	}

	cfgData, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read spark config: (%s): %w", s.Id, err)
	}

	return cfgData, nil
}
