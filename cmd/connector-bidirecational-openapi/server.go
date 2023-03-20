package main

import connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"

type mockServer struct {
	bindHost  string
	bindPort  int
	spec      string
	onRequest func(path string, body []byte, headers connectorv1.Headers) (*response, error)
}

type response struct {
	body    []byte
	headers connectorv1.Headers
	path    string
}

func (s mockServer) start() error {
	return nil
}

func (s mockServer) stop() error {
	return nil
}
