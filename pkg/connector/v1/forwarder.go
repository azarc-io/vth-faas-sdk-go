package connectorv1

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type requestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type forwarderOption func(forwarder) forwarder

type forwarder struct {
	httpClient requestDoer
	config     *connectorConfig
}

type forwardData struct {
	Tenant      string            `json:"tenant"`
	MsgName     string            `json:"message_name"`
	ConnectorID string            `json:"connector_id"`
	HeadersMap  map[string]string `json:"headers"`
	Payload     []byte            `json:"payload"`
}

func (f forwardData) Body() Bindable {
	return NewBindable(f.Payload, BindableTypeJson)
}

func (f forwardData) Headers() Headers {
	return f.HeadersMap
}

func (f forwardData) MessageName() string {
	return f.MsgName
}

func (f *forwarder) Forward(name string, body []byte, headers Headers) (InboundResponse, error) {
	req := forwardData{
		Tenant:      f.config.Tenant,
		MsgName:     name,
		ConnectorID: f.config.Id,
		HeadersMap:  headers,
		Payload:     body,
	}

	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, f.config.Agent.forwarderURL(), bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	response, err := f.httpClient.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		// TODO: how to handle?
	}
	var resp forwardData
	err = json.NewDecoder(response.Body).Decode(&resp)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func newForwarder(config *connectorConfig, opts ...forwarderOption) Forwarder {
	fwd := forwarder{config: config}
	for _, opt := range opts {
		fwd = opt(fwd)
	}
	if fwd.httpClient == nil {
		fwd.httpClient = http.DefaultClient
	}
	return &fwd
}

func withRequestDoer(doer requestDoer) forwarderOption {
	return func(f forwarder) forwarder {
		f.httpClient = doer
		return f
	}
}
