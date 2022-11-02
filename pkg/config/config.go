package config

import (
	"context"
	"fmt"
	"github.com/lithammer/shortuuid/v4"
	"github.com/sethvargo/go-envconfig"
	"time"
)

type Config struct {
	App struct {
		Environment string `env:"APP_ENVIRONMENT,required"`
		Component   string `env:"APP_COMPONENT,default=job_worker"`
		InstanceId  string
	}
	Log struct {
		Level string `env:"LOG_LEVEL,default=info"`
	}
	AgentService struct {
		Port string `env:"AGENT_SERVER_PORT,required"`
	}
	ManagerService ManagerService
}

type ManagerService struct {
	HeartBeatInterval time.Duration `env:"MANAGER_HEARTBEAT_INTERVAL,default=5s"`
	Host              string        `env:"MANAGER_SERVER_PORT,required"`
	Port              int           `env:"MANAGER_SERVER_PORT,required"`
	RetryBackoff      time.Duration `env:"MANAGER_RETRY_BACKOFF_DURATION,default=1s"`
	MaxRetries        int           `env:"MANAGER_RETRY_ATTEMPTS,default=20"`
}

func (m ManagerService) HostPort() string {
	return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func New() (*Config, error) {
	config := &Config{}
	if err := envconfig.Process(context.Background(), config); err != nil {
		fmt.Printf("error loading configuration: %s", err.Error())
		return nil, err
	}
	return config, nil
}

func NewMock(mapper map[string]string) (*Config, error) {
	config := &Config{}
	err := envconfig.ProcessWith(context.Background(), config, envconfig.MapLookuper(mapper))
	if err != nil {
		fmt.Printf("error loading configuration: %s", err.Error())
		return nil, err
	}
	config.App.InstanceId = shortuuid.New()
	return config, nil
}
