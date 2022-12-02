package spark_v1

import (
	"context"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"

	"github.com/sethvargo/go-envconfig"
)

/************************************************************************/
// SPARK CONFIG
/************************************************************************/

type config struct {
	Config struct {
		Health *configHealth `yaml:"health"`
		Agent  *configAgent  `yaml:"agent"`
		Server *configServer `yaml:"server"`
		Log    *configLog    `yaml:"logging"`
		App    *configApp    `yaml:"app"`
	}
}

type configHealth struct {
	Enabled bool   `env:"HEALTH_ENABLED" yaml:"enabled"`
	Bind    string `env:"SERVER_BIND" yaml:"bind"`
	Port    int    `env:"HEALTH_PORT" yaml:"port"`
}

type configAgent struct {
	Host         string        `env:"AGENT_HOST" yaml:"host"`
	Port         int           `env:"AGENT_PORT" yaml:"port"`
	RetryBackoff time.Duration `env:"AGENT_RETRY_BACKOFF_DURATION" yaml:"retryBackoff"`
	MaxRetries   int           `env:"AGENT_RETRY_ATTEMPTS" yaml:"maxRetries"`
}

type configServer struct {
	Bind    string `env:"SERVER_BIND" yaml:"bind"`
	Port    int    `env:"SERVER_PORT" yaml:"port"`
	Enabled bool   `env:"SERVER_ENABLED" yaml:"enabled"`
}

type configLog struct {
	Level string `env:"LOG_LEVEL" yaml:"level"`
}

type configApp struct {
	Environment string `env:"APP_ENVIRONMENT" yaml:"environment"`
	Component   string `env:"APP_COMPONENT" yaml:"component"`
	InstanceID  string `env:"APP_INSTANCE_ID" yaml:"instanceId"`
}

func (m *config) agentAddress() string {
	return fmt.Sprintf("%s:%d", m.Config.Agent.Host, m.Config.Agent.Port)
}

func (m *config) serverAddress() string {
	return fmt.Sprintf("%s:%d", m.Config.Server.Bind, m.Config.Server.Port)
}

func (m *config) healthBindTo() string {
	return fmt.Sprintf("%s:%d", m.Config.Health.Bind, m.Config.Health.Port)
}

func loadSparkConfig() (*config, error) {
	config := &config{}

	if os.Getenv("SPARK_FILE_PATH") != "" {
		b, err := os.ReadFile(os.Getenv("SPARK_FILE_PATH"))
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	if os.Getenv("SPARK_SECRET") != "" {
		if err := yaml.Unmarshal([]byte(os.Getenv("SPARK_SECRET")), &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// check for a yaml config
	if _, err := os.Stat("spark.yaml"); err == nil {
		b, err := os.ReadFile("spark.yaml")
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

/************************************************************************/
// USER CONFIG
/************************************************************************/

type bindableConfig struct {
	b        []byte
	filePath string
	opts     *sparkOpts
}

func newBindableConfig(opts *sparkOpts) BindableConfig {
	c := &bindableConfig{opts: opts}

	if opts.config != nil {
		return c
	}

	var err error

	// optional config path environment
	yamlFilePath := "config.yaml"
	jsonFilePath := "config.json"

	if os.Getenv("CONFIG_SECRET") != "" {
		c.opts.configType = ConfigTypeJson
		c.b = []byte(os.Getenv("CONFIG_SECRET"))
	} else if os.Getenv("CONFIG_FILE_PATH") != "" {
		c.filePath = os.Getenv("CONFIG_FILE_PATH")
		if c.b, err = os.ReadFile(c.filePath); err != nil {
			panic(err)
		}
	} else if _, err = os.Stat(yamlFilePath); err == nil {
		c.filePath = yamlFilePath
		if c.b, err = os.ReadFile(yamlFilePath); err != nil {
			panic(err)
		}
	} else if _, err = os.Stat(jsonFilePath); err == nil {
		c.filePath = jsonFilePath
		if c.b, err = os.ReadFile(jsonFilePath); err != nil {
			panic(err)
		}
	}

	return c
}

func (r *bindableConfig) Raw() ([]byte, error) {
	return r.b, nil
}

func (r *bindableConfig) Bind(target any) error {
	if r.opts.config != nil {
		switch r.opts.configType {
		case ConfigTypeJson:
			return json.Unmarshal(r.opts.config, &target)
		case ConfigTypeYaml:
			return yaml.Unmarshal(r.opts.config, &target)
		}
	}

	if strings.HasSuffix(r.filePath, ".yaml") {
		return yaml.Unmarshal(r.b, &target)
	} else if strings.HasSuffix(r.filePath, ".json") {
		return json.Unmarshal(r.b, &target)
	}

	return fmt.Errorf("can not load config, unsupported extension: %s", r.filePath)
}
