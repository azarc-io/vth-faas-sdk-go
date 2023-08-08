package connectorv1

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

const agentTokenHeader = "X-Token"

type requestDoer interface {
	Do(req *http.Request) (*http.Response, error)
}

type forwarderOption func(forwarder) forwarder

type forwarder struct {
	httpClient requestDoer
	config     *connectorConfig
	logger     Logger
}

type forwardData struct {
	Tenant        string            `json:"tenant"`
	MsgName       string            `json:"message_name"`
	ConnectorID   string            `json:"connector_id"`
	ArcID         string            `json:"arc_id"`
	EnvironmentID string            `json:"environment_id"`
	StageID       string            `json:"stage_id"`
	HeadersMap    map[string]string `json:"headers"`
	Payload       json.RawMessage   `json:"payload"`
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
	// TODO: Body must be JSON object for now but we must change to bytes after agent update
	if len(body) > 0 {
		if err := json.Unmarshal(body, &map[string]any{}); err != nil {
			f.logger.Error(err, "request body must be a valid json object")
			return nil, errors.New("request body must be a valid json object")
		}
	}
	if headers == nil {
		headers = Headers{}
	}
	req := forwardData{
		Tenant:        f.config.Tenant,
		MsgName:       name,
		ConnectorID:   f.config.Id,
		ArcID:         f.config.ArcID,
		EnvironmentID: f.config.EnvironmentID,
		StageID:       f.config.StageID,
		HeadersMap:    headers,
		Payload:       body,
	}

	data, err := json.Marshal(req)
	if err != nil {
		f.logger.Error(err, "[forwarder] failed to marshal request")
		return nil, err
	}
	request, err := http.NewRequest(http.MethodPost, f.config.Agent.forwarderURL(), bytes.NewReader(data))
	if err != nil {
		f.logger.Error(err, "[forwarder] failed to create a new request")
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set(agentTokenHeader, f.config.Agent.Token)
	f.logger.Debug("[forwarder] %s %s: %s", request.Method, request.URL.String(), string(data))

	response, err := f.httpClient.Do(request)
	if err != nil {
		f.logger.Error(err, "[forwarder] failed to do the request")
		return nil, err
	}
	body, err = io.ReadAll(response.Body)
	if err != nil {
		f.logger.Error(err, "[forwarder] failed to read response body")
		return nil, err
	}
	f.logger.Debug("[forwarder] response body: %s", string(body))
	if response.StatusCode != http.StatusOK {
		return nil, &HttpError{
			HttpCode: response.StatusCode,
			Reason:   response.Status,
			Raw:      body,
		}
	}
	var resp forwardData
	err = json.Unmarshal(body, &resp)
	if err != nil {
		f.logger.Error(err, "[forwarder] failed to decode response body. body: %s", string(body))
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
	if fwd.logger == nil {
		fwd.logger = noopLogger{}
	}
	return &fwd
}

func withRequestDoer(doer requestDoer) forwarderOption {
	return func(f forwarder) forwarder {
		f.httpClient = doer
		return f
	}
}

func withLogger(logger Logger) forwarderOption {
	return func(f forwarder) forwarder {
		f.logger = logger
		return f
	}
}
