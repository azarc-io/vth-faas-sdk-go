package connectorv1

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/sethvargo/go-envconfig"
	"gopkg.in/yaml.v3"
	"os"
	"path"
	"reflect"
	"strings"
	"time"
)

/************************************************************************/
// CONNECTOR CONFIG
/************************************************************************/

type Config struct {
	Ingress         []ingressConfig `yaml:"ingress"`
	ConnectorConfig connectorConfig `yaml:"config"`
}

type connectorConfig struct {
	Id     string        `yaml:"id"`
	Name   string        `yaml:"name"`
	Tenant string        `yaml:"tenant"`
	Agent  *agent        `yaml:"agent"`
	Health *configHealth `yaml:"health"`
	Log    *configLog    `yaml:"logging"`
}

type ingressConfig struct {
	Name     string                `yaml:"name"`
	Enabled  bool                  `yaml:"enabled"`
	Type     string                `yaml:"type"`
	Bind     ingressBindConfig     `yaml:"bind"`
	Endpoint ingressEndpointConfig `yaml:"endpoint"`
}

func (i ingressConfig) ExternalAddress() string {
	return fmt.Sprintf("%s://%s:%d%s", i.Endpoint.Protocol, i.Endpoint.Host, i.Endpoint.Port, i.Endpoint.Path)
}

func (i ingressConfig) InternalPort() int {
	return i.Bind.Port
}

func (i ingressConfig) InternalHost() string {
	return i.Bind.Host
}

type ingressBindConfig struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}
type ingressEndpointConfig struct {
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	Path     string `yaml:"path"`
	Protocol string `yaml:"protocol"`
}

type configHealth struct {
	Enabled  bool          `yaml:"enabled"`
	Bind     string        `yaml:"bind"`
	Port     int           `yaml:"port"`
	Interval time.Duration `yaml:"interval"`
}

type configLog struct {
	Level string `env:"LOG_LEVEL" yaml:"level"`
}

type agent struct {
	Host      string `env:"AGENT_HOST" yaml:"host"`
	Port      int    `env:"AGENT_PORT" yaml:"port"`
	Forwarder struct {
		Path string `yaml:"path"`
	} `yaml:"forwarder"`
}

func (a agent) forwarderURL() string {
	return fmt.Sprintf("http://%s:%d%s", a.Host, a.Port, a.Forwarder.Path)
}

func loadConnectorConfig(opts *ConnectorOpts) (*Config, error) {
	config := &Config{}

	if os.Getenv("CONNECTOR_FILE_PATH") != "" {
		b, err := os.ReadFile(os.Getenv("CONNECTOR_FILE_PATH"))
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	if os.Getenv("CONNECTOR_SECRET") != "" {
		secret, err := base64.StdEncoding.DecodeString(os.Getenv("CONNECTOR_SECRET"))
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(secret, &config); err != nil {
			return nil, err
		}
		return config, nil
	}

	// check for a yaml config
	sparkPath := path.Join(opts.configBasePath, "connector.yaml")
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
// MESSAGE DESCRIPTORS CONFIG
/************************************************************************/

type messageDescriptor struct {
	ID           string      `json:"id" yaml:"id"`
	ReadableName string      `json:"name" yaml:"name"`
	MsgName      string      `json:"message_name" yaml:"message_name"`
	Mime         string      `json:"mime_type" yaml:"mime_type"`
	Type         MessageType `json:"type" yaml:"type"`
	Options      []byte      `json:"options" yaml:"options"`
}

func (m messageDescriptor) Name() string {
	return m.ReadableName
}

func (m messageDescriptor) MessageName() string {
	return m.MsgName
}

func (m messageDescriptor) MimeType() string {
	return m.Mime
}

func (m messageDescriptor) MessageType() MessageType {
	return m.Type
}

func (m messageDescriptor) Config() Bindable {
	var tp string
	if m.Mime != "" {
		parts := strings.Split(m.Mime, "/")
		tp = parts[len(parts)-1]
	}
	return NewBindable(m.Options, BindableType(tp))
}

func loadMessageDescriptorsConfig(messageType MessageType) ([]messageDescriptor, error) {
	var descriptors []messageDescriptor

	filePathEnvVar := os.Getenv(fmt.Sprintf("%s_DESCRIPTOR_FILE_PATH", strings.ToUpper(string(messageType))))
	secretEnvVar := os.Getenv(fmt.Sprintf("%s_DESCRIPTOR_SECRET", strings.ToUpper(string(messageType))))
	if filePathEnvVar != "" {
		b, err := os.ReadFile(filePathEnvVar)
		if err != nil {
			return nil, err
		}
		if err := yaml.Unmarshal(b, &descriptors); err != nil {
			return nil, err
		}
	} else if secretEnvVar != "" {
		secret, err := base64.StdEncoding.DecodeString(secretEnvVar)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(secret, &descriptors); err != nil {
			return nil, err
		}
	}

	return descriptors, nil
}

/************************************************************************/
// USER CONFIG
/************************************************************************/

func loadUserConfig(opts *ConnectorOpts) (Bindable, error) {
	configSecretEnvVar := os.Getenv("CONFIG_SECRET")
	if configSecretEnvVar != "" {
		bytes, err := base64.StdEncoding.DecodeString(configSecretEnvVar)
		if err != nil {
			return nil, err
		}
		return NewBindable(bytes, BindableTypeJson), nil
	}

	configFilePathEnvVar := os.Getenv("CONFIG_FILE_PATH")
	if configFilePathEnvVar != "" {
		bytes, err := os.ReadFile(configFilePathEnvVar)
		if err != nil {
			return nil, err
		}
		parts := strings.Split(configFilePathEnvVar, ".")
		tp := parts[len(parts)-1]
		return NewBindable(bytes, BindableType(tp)), nil
	}

	if opts.config != nil {
		return NewBindable(opts.config, BindableTypeUnknown), nil
	}

	yamlFilePath := path.Join(opts.configBasePath, "config.yaml")
	if _, err := os.Stat(yamlFilePath); err == nil {
		bytes, err := os.ReadFile(yamlFilePath)
		if err != nil {
			return nil, err
		}
		return NewBindable(bytes, BindableTypeYaml), nil
	}

	jsonFilePath := path.Join(opts.configBasePath, "config.json")
	if _, err := os.Stat(jsonFilePath); err == nil {
		bytes, err := os.ReadFile(jsonFilePath)
		if err != nil {
			return nil, err
		}
		return NewBindable(bytes, BindableTypeJson), nil
	}

	return nil, errors.New("config not provided")
}

/************************************************************************/
// Bindable
/************************************************************************/

type BindableType string

const (
	BindableTypeJson    = "json"
	BindableTypeYaml    = "yaml"
	BindableTypeUnknown = ""
)

type bindable struct {
	data     []byte
	dataType BindableType
}

func (r *bindable) Raw() ([]byte, error) {
	return r.data, nil
}

func (r *bindable) Bind(target any) error {
	if reflect.ValueOf(target).Kind() != reflect.Ptr {
		return fmt.Errorf("expected bind parameter to be ptr; is %s", reflect.ValueOf(target).Kind().String())
	}

	switch r.dataType {
	case BindableTypeYaml:
		return yaml.Unmarshal(r.data, target)
	case BindableTypeJson, BindableTypeUnknown:
		return json.Unmarshal(r.data, target)
	default:
		return fmt.Errorf("can not load config, unsupported extension: %s", r.dataType)
	}
}

func NewBindable(data []byte, tp BindableType) Bindable {
	return &bindable{data: data, dataType: tp}
}
