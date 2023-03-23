package main

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
)

type mockClient struct {
	connected bool
	address   string
	spec      string
}

func (c *mockClient) connect() error {
	c.connected = true
	return nil
}

func (c *mockClient) disconnect() error {
	c.connected = false
	return nil
}

func (c *mockClient) DoExternalRequest(
	endpoint string,
	mimeType string,
	body []byte,
	headers connectorv1.Headers,
) ([]byte, connectorv1.Headers, error) {
	return []byte(`"result": "hello!"`), map[string]string{}, nil
}
