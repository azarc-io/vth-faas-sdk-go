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
	dummyRequestBody := []byte(`{"key":"request"}`)
	dummyResponseBody := []byte(`{"key":"response"}`)
	agentConfig := &agent{
		Host:  "test.agent",
		Port:  8080,
		Token: "test-token",
		Forwarder: struct {
			Path string `yaml:"path"`
		}{
			Path: "/forward",
		},
	}
	connectorConfig := &connectorConfig{
		Id:            "connector-id",
		Name:          "connector-name",
		Tenant:        "connector-tenant",
		ArcID:         "arc-id",
		EnvironmentID: "env-id",
		StageID:       "stage-id",
		Agent:         agentConfig,
	}
	mockClient := mockHttpDoer{DoFunc: func(req *http.Request) (*http.Response, error) {
		assert.Equal(t, "http://test.agent:8080/forward", req.URL.String())
		assert.Equal(t, "application/json", req.Header.Get("Content-Type"))
		assert.Equal(t, agentConfig.Token, req.Header.Get("X-Token"))

		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		var data forwardData
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
		assert.Equal(t, connectorConfig.Id, data.ConnectorID)
		assert.Equal(t, connectorConfig.ArcID, data.ArcID)
		assert.Equal(t, connectorConfig.EnvironmentID, data.EnvironmentID)
		assert.Equal(t, connectorConfig.StageID, data.StageID)
		assert.Equal(t, connectorConfig.Tenant, data.Tenant)
		assert.Equal(t, "test-name", data.MsgName)
		assert.Equal(t, 1, len(data.HeadersMap))
		assert.Equal(t, "test-value", data.HeadersMap["test-header"])
		assert.Equal(t, dummyRequestBody, []byte(data.Payload))

		resp := forwardData{
			Payload: dummyResponseBody,
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

	resp, err := fwd.Forward("test-name", dummyRequestBody, map[string]string{
		"test-header": "test-value",
	})
	assert.NoError(t, err)

	respData, err := resp.Body().Raw()
	assert.NoError(t, err)

	assert.Equal(t, 1, len(resp.Headers()))
	assert.Equal(t, "response-value", resp.Headers()["response-header"])
	assert.Equal(t, dummyResponseBody, respData)
}

func TestForwardWithNilHeader(t *testing.T) {
	dummyEmptyBody := []byte(`{}`)
	agentConfig := &agent{
		Host:  "test.agent",
		Port:  8080,
		Token: "test-token",
		Forwarder: struct {
			Path string `yaml:"path"`
		}{
			Path: "/forward",
		},
	}
	connectorConfig := &connectorConfig{
		Id:            "connector-id",
		Name:          "connector-name",
		Tenant:        "connector-tenant",
		ArcID:         "arc-id",
		EnvironmentID: "env-id",
		StageID:       "stage-id",
		Agent:         agentConfig,
	}
	mockClient := mockHttpDoer{DoFunc: func(req *http.Request) (*http.Response, error) {
		body, err := io.ReadAll(req.Body)
		assert.NoError(t, err)
		var data forwardData
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)
		assert.Equal(t, Headers{}, data.HeadersMap)

		resp := forwardData{
			Payload: dummyEmptyBody,
		}
		respBytes, _ := json.Marshal(resp)
		return &http.Response{
			Status:     http.StatusText(http.StatusOK),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader(respBytes)),
		}, nil
	}}
	fwd := newForwarder(connectorConfig, withRequestDoer(mockClient))

	resp, err := fwd.Forward("test-name", dummyEmptyBody, nil)
	assert.NoError(t, err)

	respData, err := resp.Body().Raw()
	assert.NoError(t, err)

	assert.Equal(t, dummyEmptyBody, respData)
}
