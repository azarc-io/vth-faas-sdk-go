package sparkv1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"reflect"
	"strings"
)

/************************************************************************/
// SPARK CONFIG
/************************************************************************/

type Config struct {
	Id         string          `yaml:"id"`
	Name       string          `yaml:"Name"`
	QueueGroup string          `yaml:"queue_group"`
	Health     *configHealth   `yaml:"health"`
	Server     *configServer   `yaml:"plugin"`
	Log        *configLog      `yaml:"logging"`
	App        *configApp      `yaml:"app"`
	Temporal   *configTemporal `yaml:"temporal"`
}

type configHealth struct {
	Enabled bool   `env:"HEALTH_ENABLED" yaml:"enabled"`
	Bind    string `env:"SERVER_BIND" yaml:"bind"`
	Port    int    `env:"HEALTH_PORT" yaml:"port"`
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

type configTemporal struct {
	Address   string `yaml:"address"`
	Namespace string `yaml:"namespace"`
}

func (m *Config) serverAddress() string {
	return fmt.Sprintf("%s:%d", m.Server.Bind, m.Server.Port)
}

func (m *Config) healthBindTo() string {
	return fmt.Sprintf("%s:%d", m.Health.Bind, m.Health.Port)
}

func loadSparkConfig(opts *SparkOpts) (*Config, error) {
	config := &Config{}

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
		secret, err := base64.StdEncoding.DecodeString(os.Getenv("SPARK_SECRET"))
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(secret, &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// check for a yaml config
	sparkPath := path.Join(opts.configBasePath, "spark.yaml")
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

/************************************************************************/
// USER CONFIG
/************************************************************************/

type bindableConfig struct {
	b        []byte
	filePath string
	opts     *SparkOpts
}

func newBindableConfig(opts *SparkOpts) BindableConfig {
	c := &bindableConfig{opts: opts}

	var err error

	// optional config path environment
	yamlFilePath := path.Join(opts.configBasePath, "config.yaml")
	jsonFilePath := path.Join(opts.configBasePath, "config.json")

	if os.Getenv("CONFIG_SECRET") != "" {
		c.opts.configType = ConfigTypeJson
		c.b, err = base64.StdEncoding.DecodeString(os.Getenv("CONFIG_SECRET"))
		if err != nil {
			panic(err)
		}
	} else if os.Getenv("CONFIG_FILE_PATH") != "" {
		c.filePath = os.Getenv("CONFIG_FILE_PATH")
		if c.b, err = os.ReadFile(c.filePath); err != nil {
			panic(err)
		}
	} else if opts.config != nil {
		return c
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
	if r.opts.config != nil {
		return r.opts.config, nil
	}
	return r.b, nil
}

func (r *bindableConfig) Bind(target any) error {
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("expected bind parameter to be ptr; is %s", reflect.ValueOf(target).Kind().String())
	}

	if r.opts.config != nil {
		return json.Unmarshal(r.opts.config, target)
	}

	if len(r.b) > 1 {
		return json.Unmarshal(r.b, target)
	}

	if strings.HasSuffix(r.filePath, ".yaml") {
		return yaml.Unmarshal(r.b, target)
	} else if strings.HasSuffix(r.filePath, ".json") {
		return json.Unmarshal(r.b, target)
	}

	return fmt.Errorf("can not load config, unsupported extension: %s", r.filePath)
}
