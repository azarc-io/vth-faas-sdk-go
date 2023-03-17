package connectorv1

/************************************************************************/
// CONFIGURATION
/************************************************************************/

type Bindable interface {
	Raw() ([]byte, error)
	Bind(any) error
}

/************************************************************************/
// LOGGING
/************************************************************************/

type (
	Logger interface {
		Error(err error, format string, v ...interface{})
		Fatal(err error, format string, v ...interface{}) // this will crash the service
		Info(format string, v ...interface{})
		Warn(format string, v ...interface{})
		Debug(format string, v ...interface{})
	}
)

/************************************************************************/
// CONTEXT
/************************************************************************/

type (
	StartContext interface {
		Config() Bindable
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

type Headers = map[string]string

type MessageType string

const MessageTypeInbound MessageType = "inbound"
const MessageTypeOutbound MessageType = "outbound"

type (
	MessageDescriptor interface {
		Name() string
		MessageName() string
		MimeType() string
		MessageType() MessageType
		Config() Bindable
	}

	request interface {
		Body() Bindable
		Headers() Headers
	}

	InboundRequest interface {
		request
		MessageName() string
		MimeType() string
	}

	OutboundResponse interface {
		request
	}

	OutboundRequest interface {
		request
		MessageName() string
		MimeType() string
	}

	InboundResponse interface {
		request
	}

	InboundDescriptor interface {
		MessageDescriptor
	}

	OutboundDescriptor interface {
		MessageDescriptor
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

type ConnectorWorker interface {
	Run()
}

type InboundConnector interface {
	Connector
}

type OutboundConnector interface {
	Connector
	HandleOutboundRequest(request OutboundRequest) (any, Headers, error) // see line 26, please
}
