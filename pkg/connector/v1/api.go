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
		Ingress(name string) (Ingress, error)
		InboundDescriptors() []InboundDescriptor
		OutboundDescriptors() []OutboundDescriptor
		Forwarder() Forwarder
		Log() Logger
		RegisterPeriodicHealthCheck(name string, fn HealthCheckFunc)
	}

	StopContext interface {
		Log() Logger
	}
)

type HealthCheckFunc func() error

/************************************************************************/
// Forwarding
/************************************************************************/

type (
	Forwarder interface {
		Forward(name string, body []byte, headers Headers) (InboundResponse, error)
	}
)

/************************************************************************/
// INGRESS
/************************************************************************/

type Ingress interface {
	ExternalAddress() string
	InternalPort() int
	InternalHost() string
}

/************************************************************************/
// MODELS
/************************************************************************/

type Headers = map[string]any

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

type Connector interface {
	// Start can be used to initialise your server or client
	// Note: This function should return and not block
	Start(ctx StartContext) error
	Stop(ctx StopContext) error
}

type ConnectorService interface {
	Start() error
}

type InboundConnector interface {
	Connector
}

type OutboundConnector interface {
	Connector
	HandleOutboundRequest(request OutboundRequest) (any, Headers, error) // see line 26, please
}
