package v1

type Headers = map[string]any

type Connector interface {
	Start() error
	Stop() error
}

type InboundConnector interface {
	Connector
	Forward(payload []byte, headers Headers)
}

type OutboundConnector interface {
	Connector
	HandleOutboundRequest(request OutboundRequest) (OutboundResponse, error) // see line 26, please
}

type OutboundRequest struct {
	Payload []byte
	Headers Headers
}

type OutboundResponse struct {
	// what fields do we put into this struct?
	// I think we should return a sdk.Variable{name, mime-type, value} from OutboundConnector.HandleOutboundRequest
	// in the end, a spark will consume that data, I guess. So, it seems reasonable and simpler to return a variable here
}

// if the only reason of this interface is to provide configuration to the connectors
// maybe we should use environment variables
type ConnectorContext interface {
}

//// pseudocode is here just for reference TODO remove before send it to main
//
// package hello_world
//
//import "context"
//
//type connector struct {
//	// Requires developer to implement Start / Stop functions
//	api.InboundConnector
//
//	// Causes the SDK to enable inbound server Agent -> Connector
//	api.OutboundConnector
//
//	// Services defined by the developer
//	mq       interface{}
//	server   interface{}
//	apiCtx   api.ConnectorContext
//}
//
//// Defined by the developer to handle their inbound requests
//func (c connector) handle(ctx echo.Context) {
//	// This is part of the SDK and is how the developer forwards a request to the Agent
//	rsp, err := c.Forward(ctx.Body(), ctx.Headers())
//	// Developer returns responses or errors back to requester
//	return rsp, err
//}
//
//// api.Connector, api.OutboundConnector
//func (c connector) Start() error {
//	// connect to ibmmq
//	// loop overinbound  message types
//	for in := c.apiCtx.GetInoundDefinitions() {
//		// subscribe to topic on IBMMQ
//		c.mq.Subscribe(in.String("topic"), c.handle)
//	}
//	// loop over outbound message types
//	for in := c.apiCtx.GetOutboundDefinitions() {
//		// create publisher for ibmmq
//	}
//}
//
//// api.Connector, api.OutboundConnector
//func (c connector) Stop() error {
//
//}
//
//// api.OutboundConnector
//func (c connector) HandleOutboundRequest(req *api_v1.OutboundRequest) (*api_v1.OutboundResponse, error) {
//	// Defined by developer
//	rsp, err := c.server.Handle(req.Payload, req.Headers)
//	// Return result to Agent
//	return rsp, err
//}
//
//func main() {
//	apiCtx := api.NewContext()
//	c := api.NewConnector(&connector{apiCtx: apiCtx}, apiCtx)
//	c.Run()
//}
