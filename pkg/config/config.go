// Copyright 2020-2022 Azarc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
