package e2e_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestWorkerE2E(t *testing.T) {
	dummyRequestBody := []byte(`{"key":"request"}`)

	_ = os.Setenv("INBOUND_DESCRIPTOR_FILE_PATH", "../fixtures/inbound_descriptors_config_1.yaml")
	_ = os.Setenv("CONFIG_FILE_PATH", "../fixtures/user_config_1.yaml")

	waitChan := make(chan struct{}, 1)
	healthCheckChan := make(chan struct{}, 1)
	// setup agent server
	server := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		assert.NoError(t, err)

		var forwardData struct {
			Tenant      string            `json:"tenant"`
			MsgName     string            `json:"message_name"`
			ConnectorID string            `json:"connector_id"`
			HeadersMap  map[string]string `json:"headers"`
			Payload     json.RawMessage   `json:"payload"`
		}
		err = json.Unmarshal(body, &forwardData)
		assert.NoError(t, err)
		assert.Equal(t, "tenant-id", forwardData.Tenant)
		assert.Equal(t, "message-name-1", forwardData.MsgName)
		assert.Equal(t, "connector-simple-example_12345", forwardData.ConnectorID)
		assert.Equal(t, map[string]string{"key": "value"}, forwardData.HeadersMap)
		assert.Equal(t, dummyRequestBody, []byte(forwardData.Payload))

		resp := map[string]any{
			"headers": map[string]string{
				"key": "header-value",
			},
			"payload": map[string]string{
				"response": "response-body",
			},
		}
		err = json.NewEncoder(writer).Encode(resp)
		assert.NoError(t, err)
	}))
	defer server.Close()
	serverUrl, err := url.Parse(server.URL)
	assert.NoError(t, err)
	port, err := strconv.Atoi(serverUrl.Port())
	assert.NoError(t, err)

	// modify connector configuration to use correct ports
	configBytes, err := os.ReadFile("../fixtures/connector_config_1.yaml")
	assert.NoError(t, err)
	var config connectorv1.Config
	err = yaml.Unmarshal(configBytes, &config)
	assert.NoError(t, err)
	config.ConnectorConfig.Agent.Host = strings.Split(serverUrl.Host, ":")[0]
	config.ConnectorConfig.Agent.Port = port

	listener, err := net.Listen("tcp", ":0")
	assert.NoError(t, err)

	config.ConnectorConfig.Health.Enabled = true
	config.ConnectorConfig.Health.Bind = "127.0.0.1"
	config.ConnectorConfig.Health.Port = listener.Addr().(*net.TCPAddr).Port
	err = listener.Close()
	assert.NoError(t, err)

	configBytes, err = yaml.Marshal(config)
	assert.NoError(t, err)
	configSecret := base64.StdEncoding.EncodeToString(configBytes)
	_ = os.Setenv("CONNECTOR_SECRET", configSecret)

	// setup connector mock
	ctrl := gomock.NewController(t)
	connector := mock.NewMockConnector(ctrl)
	connector.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx connectorv1.StartContext) error {
		defer func() {
			waitChan <- struct{}{}
		}()

		ctx.RegisterPeriodicHealthCheck("health-check", func() error {
			healthCheckChan <- struct{}{}
			return nil
		})
		resp, err := ctx.Forwarder().Forward("message-name-1", dummyRequestBody, map[string]string{
			"key": "value",
		})
		assert.NoError(t, err)
		var respBody struct {
			Response string `json:"response"`
		}
		err = resp.Body().Bind(&respBody)
		assert.NoError(t, err)
		assert.Equal(t, "response-body", respBody.Response)
		assert.Equal(t, "header-value", resp.Headers()["key"])
		return nil
	})
	connector.EXPECT().Stop(gomock.Any()).Return(nil)

	worker, err := connectorv1.NewConnectorWorker(connector)
	assert.NoError(t, err)

	go func() {
		<-waitChan
		// verify health check
		resp, err := http.Get(fmt.Sprintf(`http://localhost:%d/healthz`, config.ConnectorConfig.Health.Port))
		assert.NoError(t, err)
		body, err := io.ReadAll(resp.Body)
		assert.NoError(t, err)
		var healthResponse struct {
			Status   string         `json:"status"`
			Failures map[string]any `json:"failures"`
		}
		err = json.Unmarshal(body, &healthResponse)
		assert.NoError(t, err)
		assert.Equal(t, "OK", healthResponse.Status)
		assert.Equal(t, 0, len(healthResponse.Failures))

		select {
		case <-healthCheckChan: // health check func has run
		case <-time.After(10 * time.Second): // health check func hasn't run
			t.Fail()
		}
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	worker.Run()
}
