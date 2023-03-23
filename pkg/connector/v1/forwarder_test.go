package connectorv1

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

type mockHttpDoer struct {
	DoFunc func(req *http.Request) (*http.Response, error)
}

func (m mockHttpDoer) Do(req *http.Request) (*http.Response, error) {
	return m.DoFunc(req)
}

func TestForward(t *testing.T) {
	agentConfig := &agent{
		Host: "test.agent",
		Port: 8080,
		Forwarder: struct {
			Path string `yaml:"path"`
		}{
			Path: "/forward",
		},
	}
	connectorConfig := &connectorConfig{
		Id:     "connector-id",
		Name:   "connector-name",
		Tenant: "connector-tenant",
		Agent:  agentConfig,
	}
	mockClient := mockHttpDoer{DoFunc: func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://test.agent:8080/forward", req.URL.String())
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		var data forwardData
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
		assert.Equal(t, connectorConfig.Id, data.ConnectorID)
		assert.Equal(t, connectorConfig.Tenant, data.Tenant)
		assert.Equal(t, "test-name", data.MsgName)
		assert.Equal(t, 1, len(data.HeadersMap))
		assert.Equal(t, "test-value", data.HeadersMap["test-header"])
		assert.Equal(t, []byte("test-body"), data.Payload)

		resp := forwardData{
			Payload: []byte("response-data"),
			HeadersMap: map[string]string{
				"response-header": "response-value",
			},
		}
		respBytes, _ := json.Marshal(resp)
		return &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respBytes)),
		}, nil
	}}
	fwd := newForwarder(connectorConfig, withRequestDoer(mockClient))

	resp, err := fwd.Forward("test-name", []byte("test-body"), map[string]string{
		"test-header": "test-value",
	})
	assert.NoError(t, err)

	respData, err := resp.Body().Raw()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Headers()))
	assert.Equal(t, "response-value", resp.Headers()["response-header"])
	assert.Equal(t, []byte("response-data"), respData)
}
