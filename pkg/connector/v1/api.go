package connectorv1

/************************************************************************/
// CONFIGURATION
/************************************************************************/

type ConfigType string

const (
	ConfigTypeYaml ConfigType = "yaml"
	ConfigTypeJson ConfigType = "json"
)

type Configuration interface {
	Bind(target any) error
	Raw() interface{}
}

/************************************************************************/
// CONTEXT
/************************************************************************/

type (
	StartContext interface {
		Config() Configuration
	}
)

/************************************************************************/
// MODELS
/************************************************************************/

type (
	OutboundRequest interface {
		Body() Bindable
		Headers() Headers
	}

	Bindable interface {
		Raw() ([]byte, error)
		Bind(any) error
	}
)

/************************************************************************/
// API
/************************************************************************/

type Headers = map[string]any

type Connector interface {
	Start(ctx StartContext) error
	Stop() error
}

type ConnectorService interface {
	Start() error
}

type InboundConnector interface {
	Connector
	Forward(payload []byte, headers Headers)
}

type OutboundConnector interface {
	Connector
	HandleOutboundRequest(request OutboundRequest) (any, Headers, error) // see line 26, please
}

type OutboundRequest1 struct {
	Payload []byte
	Headers Headers
}

type OutboundResponse2 struct {
	Payload []byte
	Headers Headers
}
