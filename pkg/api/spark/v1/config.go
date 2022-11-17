package spark_v1

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"time"

	"github.com/sethvargo/go-envconfig"
)

type Config struct {
	Config struct {
		Health *configHealth `yaml:"health"`
		Agent  *configAgent  `yaml:"agent"`
		Server *configServer `yaml:"server"`
		Log    *configLog    `yaml:"logging"`
		App    *configApp    `yaml:"app"`
	}
}

type configHealth struct {
	Enabled string `env:"HEALTH_ENABLED" yaml:"enabled"`
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

func (m Config) AgentAddress() string {
	return fmt.Sprintf("%s:%d", m.Config.Agent.Host, m.Config.Agent.Port)
}

func (m Config) ServerAddress() string {
	return fmt.Sprintf("%s:%d", m.Config.Server.Bind, m.Config.Server.Port)
}

func loadConfig() (*Config, error) {
	config := defaultConfig()

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

func defaultConfig() *Config {
	c := &Config{}

	return c
}
