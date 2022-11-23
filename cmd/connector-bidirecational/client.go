package main

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
)

type mockClient struct {
	connected bool
	address   *config
}

func (c mockClient) Connect() error {
	c.connected = true
	return nil
}

func (c mockClient) Disconnect() error {
	c.connected = false
	return nil
}

func (c mockClient) DoExternalRequest(body []byte, headers connectorv1.Headers) ([]byte, map[string]any, error) {
	return []byte(`"result": "hello!"`), map[string]any{}, nil
}
