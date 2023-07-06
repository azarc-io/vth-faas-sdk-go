package module_runner

import (
	"context"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"time"
)

/************************************************************************/
// RUNNER CONFIG
/************************************************************************/

type config struct {
	BinBasePath string          `yaml:"bin_base_path"`
	Health      *configHealth   `yaml:"health"`
	Log         *configLog      `yaml:"logging"`
	Sparks      []*configSpark  `yaml:"sparks"`
	Temporal    *configTemporal `yaml:"temporal"`
	IOServer    *ioServer       `yaml:"io_server"`
}

func defaultConfig() *config {
	return &config{
		BinBasePath: "./",
	}
}

type configSpark struct {
	Id             string         `yaml:"id"`              // Id is unique hash to identify this combination of Name and Config
	Name           string         `yaml:"name"`            // Name of the binary to execute
	QueueGroup     string         `yaml:"queue_group"`     // QueueGroup name of execution group
	Config         string         `yaml:"config"`          // Config Deprecated: will be JSON string with config details
	ConfigServer   *configServer  `yaml:"config_server"`   // ConfigServer which is used to retrieve startup config
	StartupTimeout *time.Duration `yaml:"startup_timeout"` // StartupTimeout amount of time to wait for spark to start before error
}

type configServer struct {
	Url    string `yaml:"url" json:"url,omitempty"`
	ApiKey string `yaml:"api_key" json:"api_key,omitempty"`
}

type configHealth struct {
	Enabled bool   `env:"HEALTH_ENABLED" yaml:"enabled"`
	Bind    string `env:"SERVER_BIND" yaml:"bind"`
	Port    int    `env:"HEALTH_PORT" yaml:"port"`
}

type configLog struct {
	Level string `env:"LOG_LEVEL" yaml:"level"`
}

type configTemporal struct {
	Address   string `yaml:"address"`
	Namespace string `yaml:"namespace"`
}

type ioServer struct {
	Url    string `env:"IO_SERVER_URL" yaml:"url"`
	ApiKey string `env:"IO_SERVER_API_KEY" yaml:"api_key"`
}

func (m *config) healthBindTo() string {
	return fmt.Sprintf("%s:%d", m.Health.Bind, m.Health.Port)
}

func LoadModuleConfig(opts ...ModuleOption) (*config, error) {
	config := defaultConfig()

	if os.Getenv("MODULE_FILE_PATH") != "" {
		b, err := os.ReadFile(os.Getenv("MODULE_FILE_PATH"))
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	if os.Getenv("MODULE_SECRET") != "" {
		if err := yaml.Unmarshal([]byte(os.Getenv("MODULE_SECRET")), &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// Note: this is not run if env vars are used (above)
	mo := moduleOpts{}
	for _, opt := range opts {
		opt(&mo)
	}

	config.BinBasePath = mo.binBasePath

	// check for a yaml config
	sparkPath := path.Join(mo.configBasePath, "module.yaml")
	if _, err := os.Stat(sparkPath); err == nil {
		b, err := os.ReadFile(sparkPath)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			return nil, err
		}
	}

	if err := envconfig.Process(context.Background(), config); err != nil {
		fmt.Printf("error loading configuration: %s", err.Error())
		return nil, err
	}
	return config, nil
}
