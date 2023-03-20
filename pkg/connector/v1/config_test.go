package connectorv1

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
	"time"
)

func verifyConnectorConfig(t *testing.T, conf *Config) {
	t.Run("ingress", func(t *testing.T) {
		assert.Equal(t, 1, len(conf.Ingress))
		assert.Equal(t, "http-8080", conf.Ingress[0].Name)
		assert.Equal(t, "0.0.0.0", conf.Ingress[0].Bind.Host)
		assert.Equal(t, 8080, conf.Ingress[0].Bind.Port)
	})

	t.Run("config", func(t *testing.T) {
		assert.Equal(t, "connector-simple-example_12345", conf.ConnectorConfig.Id)
		assert.Equal(t, "connector-simple-example", conf.ConnectorConfig.Name)
		assert.Equal(t, "tenant-id", conf.ConnectorConfig.Tenant)

		assert.Equal(t, "0.0.0.0", conf.ConnectorConfig.Health.Bind)
		assert.Equal(t, 8081, conf.ConnectorConfig.Health.Port)
		assert.Equal(t, time.Second*30, conf.ConnectorConfig.Health.Interval)

		assert.Equal(t, "127.0.0.1", conf.ConnectorConfig.Agent.Host)
		assert.Equal(t, 8031, conf.ConnectorConfig.Agent.Port)

		assert.Equal(t, "info", conf.ConnectorConfig.Log.Level)
	})
}

func TestLoadConnectorConfigFromEnvPathFileName(t *testing.T) {
	_ = os.Setenv("CONNECTOR_FILE_PATH", "./test/connector_config_1.yaml")
	defer func() {
		_ = os.Unsetenv("CONNECTOR_FILE_PATH")
	}()

	conf, err := loadConnectorConfig(&ConnectorOpts{})
	assert.NoError(t, err)

	verifyConnectorConfig(t, conf)
}

const connectorConfigSecret = "aW5ncmVzczoNCiAgLSBuYW1lOiAiaHR0cC04MDgwIg0KICAgIGVuYWJsZWQ6IHRydWUNCiAgICBwb3J0OiA4MDgwDQogICAgdHlwZTogaHR0cA0KICAgIGJpbmQ6DQogICAgICBob3N0OiAwLjAuMC4wDQogICAgICBwb3J0OiA4MDgwDQogICAgZW5kcG9pbnQ6DQogICAgICBob3N0OiBzb21lLWV4dGVybmFsLXVybC5jb20NCiAgICAgIHBvcnQ6IDQ0Mw0KICAgICAgcGF0aDogL3YxL2VuZHBvaW50LzEyMzU2DQogICAgICBwcm90b2NvbDogaHR0cHMNCmNvbmZpZzoNCiAgaWQ6IGNvbm5lY3Rvci1zaW1wbGUtZXhhbXBsZV8xMjM0NQ0KICBuYW1lOiBjb25uZWN0b3Itc2ltcGxlLWV4YW1wbGUNCiAgdGVuYW50OiB0ZW5hbnQtaWQNCiAgaGVhbHRoOg0KICAgIGVuYWJsZWQ6IGZhbHNlDQogICAgYmluZDogMC4wLjAuMA0KICAgIHBvcnQ6IDgwODENCiAgICBpbnRlcnZhbDogMzBzDQogIGFnZW50Og0KICAgIGhvc3Q6IDEyNy4wLjAuMQ0KICAgIHBvcnQ6IDgwMzENCiAgICBmb3J3YXJkZXI6DQogICAgICBwYXRoOiAvdjEvY29ubmVjdG9yLWZvcndhcmRlZC1tZXNzYWdlDQogIGxvZ2dpbmc6DQogICAgbGV2ZWw6ICJpbmZvIg=="

func TestLoadConnectorConfigFromEnvVarSecret(t *testing.T) {
	_ = os.Setenv("CONNECTOR_SECRET", connectorConfigSecret)
	defer func() {
		_ = os.Unsetenv("CONNECTOR_SECRET")
	}()

	conf, err := loadConnectorConfig(&ConnectorOpts{})
	assert.NoError(t, err)

	verifyConnectorConfig(t, conf)
}

func TestLoadConnectorConfigFromBaseFilePath(t *testing.T) {

	conf, err := loadConnectorConfig(&ConnectorOpts{configBasePath: "./test"})
	assert.NoError(t, err)

	verifyConnectorConfig(t, conf)
}

func verifyMessageDescriptors(t *testing.T, descriptors []messageDescriptor) {
	assert.Equal(t, 2, len(descriptors))
	mimeTypes := []string{"application/json", "application/yaml"}
	for i, desc := range descriptors {
		index := fmt.Sprint(i + 1)
		assert.Equal(t, "User friendly name "+index, desc.Name())
		assert.Equal(t, "message-name-"+index, desc.MessageName())
		assert.Equal(t, MessageTypeInbound, desc.MessageType())
		assert.Equal(t, mimeTypes[i], desc.MimeType())
		var messageConfig struct {
			Key string `json:"test-key" yaml:"test-key"`
		}
		err := desc.Config().Bind(&messageConfig)
		assert.NoError(t, err)
		assert.Equal(t, "test-value-"+index, messageConfig.Key)
	}
}

func TestLoadMessageDescriptorsConfigFromEnvFilePath(t *testing.T) {
	_ = os.Setenv("INBOUND_DESCRIPTOR_FILE_PATH", "./test/inbound_descriptors_config_1.yaml")
	defer func() {
		_ = os.Unsetenv("INBOUND_DESCRIPTOR_FILE_PATH")
	}()

	descriptors, err := loadMessageDescriptorsConfig(MessageTypeInbound)
	assert.NoError(t, err)

	verifyMessageDescriptors(t, descriptors)
}

const inboundDescriptorsSecret = "LSBpZDogbWVzc2FnZV9pZF8xDQogIG5hbWU6IFVzZXIgZnJpZW5kbHkgbmFtZSAxDQogIG1lc3NhZ2VfbmFtZTogIm1lc3NhZ2UtbmFtZS0xIg0KICBtaW1lX3R5cGU6IGFwcGxpY2F0aW9uL2pzb24NCiAgdHlwZTogImluYm91bmQiDQogIG9wdGlvbnM6ICMgdGhpcyBpcyBbXWJ5dGUgcmVwcmVzZW50YXRpb24gaW4geWFtbA0KICAgIC0gMTIzDQogICAgLSAxMA0KICAgIC0gMzINCiAgICAtIDMyDQogICAgLSAzMg0KICAgIC0gMzINCiAgICAtIDM0DQogICAgLSAxMTYNCiAgICAtIDEwMQ0KICAgIC0gMTE1DQogICAgLSAxMTYNCiAgICAtIDQ1DQogICAgLSAxMDcNCiAgICAtIDEwMQ0KICAgIC0gMTIxDQogICAgLSAzNA0KICAgIC0gNTgNCiAgICAtIDMyDQogICAgLSAzNA0KICAgIC0gMTE2DQogICAgLSAxMDENCiAgICAtIDExNQ0KICAgIC0gMTE2DQogICAgLSA0NQ0KICAgIC0gMTE4DQogICAgLSA5Nw0KICAgIC0gMTA4DQogICAgLSAxMTcNCiAgICAtIDEwMQ0KICAgIC0gNDUNCiAgICAtIDQ5DQogICAgLSAzNA0KICAgIC0gMTANCiAgICAtIDMyDQogICAgLSAzMg0KICAgIC0gMTI1DQotIGlkOiBtZXNzYWdlX2lkXzINCiAgbmFtZTogVXNlciBmcmllbmRseSBuYW1lIDINCiAgbWVzc2FnZV9uYW1lOiAibWVzc2FnZS1uYW1lLTIiDQogIG1pbWVfdHlwZTogYXBwbGljYXRpb24veWFtbA0KICB0eXBlOiAiaW5ib3VuZCINCiAgb3B0aW9uczoNCiAgICAtIDExNg0KICAgIC0gMTAxDQogICAgLSAxMTUNCiAgICAtIDExNg0KICAgIC0gNDUNCiAgICAtIDEwNw0KICAgIC0gMTAxDQogICAgLSAxMjENCiAgICAtIDU4DQogICAgLSAzMg0KICAgIC0gMTE2DQogICAgLSAxMDENCiAgICAtIDExNQ0KICAgIC0gMTE2DQogICAgLSA0NQ0KICAgIC0gMTE4DQogICAgLSA5Nw0KICAgIC0gMTA4DQogICAgLSAxMTcNCiAgICAtIDEwMQ0KICAgIC0gNDUNCiAgICAtIDUw"

func TestLoadMessageDescriptorsConfigFromEnvSecret(t *testing.T) {
	_ = os.Setenv("INBOUND_DESCRIPTOR_SECRET", inboundDescriptorsSecret)
	defer func() {
		_ = os.Unsetenv("INBOUND_DESCRIPTOR_SECRET")
	}()

	descriptors, err := loadMessageDescriptorsConfig(MessageTypeInbound)
	assert.NoError(t, err)

	verifyMessageDescriptors(t, descriptors)
}

func verifyUserConfig(t *testing.T, config Bindable) {
	var conf struct {
		Name  string `yaml:"name" json:"name"`
		Value string `yaml:"value" json:"value"`
	}
	err := config.Bind(&conf)
	assert.NoError(t, err)
	assert.Equal(t, "user config", conf.Name)
	assert.Equal(t, "user-defined config", conf.Value)
}

func TestLoadUserConfigWithEnvFilePath(t *testing.T) {

	t.Run("yaml", func(t *testing.T) {
		_ = os.Setenv("CONFIG_FILE_PATH", "./test/user_config_1.yaml")
		defer func() {
			_ = os.Unsetenv("CONFIG_FILE_PATH")
		}()

		conf, err := loadUserConfig(&ConnectorOpts{})
		assert.NoError(t, err)

		verifyUserConfig(t, conf)
	})

	t.Run("json", func(t *testing.T) {
		_ = os.Setenv("CONFIG_FILE_PATH", "./test/config.json")
		defer func() {
			_ = os.Unsetenv("CONFIG_FILE_PATH")
		}()

		conf, err := loadUserConfig(&ConnectorOpts{})
		assert.NoError(t, err)

		verifyUserConfig(t, conf)
	})
}

const userConfigSecret = "ew0KICAibmFtZSI6ICJ1c2VyIGNvbmZpZyIsDQogICJ2YWx1ZSI6ICJ1c2VyLWRlZmluZWQgY29uZmlnIg0KfQ0K"

func TestLoadUserConfigFromEnvSecret(t *testing.T) {
	_ = os.Setenv("CONFIG_SECRET", userConfigSecret)
	defer func() {
		_ = os.Unsetenv("CONFIG_SECRET")
	}()

	conf, err := loadUserConfig(&ConnectorOpts{})
	assert.NoError(t, err)

	verifyUserConfig(t, conf)
}

func TestLoadUserConfigFromOptionsConfig(t *testing.T) {
	conf, err := loadUserConfig(&ConnectorOpts{config: []byte(`{
		"name": "user config",
		"value": "user-defined config"
	}`)})
	assert.NoError(t, err)

	verifyUserConfig(t, conf)
}

func TestLoadUserConfigFromOptionsBasePath(t *testing.T) {
	conf, err := loadUserConfig(&ConnectorOpts{configBasePath: "./test"})
	assert.NoError(t, err)

	verifyUserConfig(t, conf)
}
