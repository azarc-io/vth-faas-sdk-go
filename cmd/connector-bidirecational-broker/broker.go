package main

import connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"

type mockBroker struct {
	bindHost string
	bindPort int
	spec     string
	address  string
}

type subscription struct {
}

type publication struct {
	topic string
}

type message struct {
	body    []byte
	headers connectorv1.Headers
	path    string
}

func (s mockBroker) connect() error {
	return nil
}

func (s mockBroker) disconnect() error {
	return nil
}

func (s mockBroker) Subscribe(cfg *subscription, inboundRequest func(req *message) error) error {
	return nil
}

func (s mockBroker) publish(topic string, m *message) error {
	return nil
}
