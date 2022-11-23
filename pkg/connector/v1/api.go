package connectorv1

/************************************************************************/
// CONFIGURATION
/************************************************************************/

type ConfigType string
type ConfigSchemaType int

const (
	ConfigTypeYaml ConfigType = "yaml"
	ConfigTypeJson ConfigType = "json"
)

const (
	InboundSchema ConfigSchemaType = iota
	OutboundSchema
)

type Configuration interface {
	Bind(target any) error
	Raw() interface{}
}

/************************************************************************/
// LOGGING
/************************************************************************/

type (
	Logger interface {
		LogError(err error, format string, v ...interface{})
		LogFatal(err error, format string, v ...interface{}) // this will crash the service
		LogInfo(format string, v ...interface{})
		LogWarn(format string, v ...interface{})
		LogDebug(format string, v ...interface{})
	}
)

/************************************************************************/
// CONTEXT
/************************************************************************/

type (
	StartContext interface {
		Config() Configuration
		Ingress() (Ingress, error)
		ForwardingContext() ForwardingContext
		InboundDescriptors() []InboundDescriptor
		OutboundDescriptors() []OutboundDescriptor
	}

	StopContext interface {
		Logger
	}

	ForwardingContext interface {
		Logger
		Forward(name string, body []byte, headers Headers) (InboundResponse, error)
	}
)

/************************************************************************/
// INGRESS
/************************************************************************/

type Ingress interface {
	IngressHost() string
	IngressPort() int
}

/************************************************************************/
// MODELS
/************************************************************************/

type (
	OutboundRequest interface {
		Body() Bindable
		Headers() Headers
		MessageName() string
		MimeType() string
	}

	Bindable interface {
		Raw() ([]byte, error)
		Bind(any) error
	}

	InboundResponse interface {
		MessageName() string
		Body() Bindable
		Headers() Headers
	}

	InboundDescriptor interface {
		OutboundRequest
		Config() Bindable
	}

	OutboundDescriptor interface {
		OutboundRequest
		Config() Bindable
	}
)

/************************************************************************/
// API
/************************************************************************/

type Headers = map[string]any

type Connector interface {
	Start(ctx StartContext) error
	Stop(ctx StopContext) error
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
