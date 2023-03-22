package mock

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=./mock_forwarder.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Forwarder
//go:generate mockgen -destination=./mock_connector.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Connector
//go:generate mockgen -destination=./mock_start_context.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 StartContext
//go:generate mockgen -destination=./mock_stop_context.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 StopContext
//go:generate mockgen -destination=./mock_ingress.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Ingress
//go:generate mockgen -destination=./mock_inbound_descriptor.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 InboundDescriptor
//go:generate mockgen -destination=./mock_outbound_descriptor.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 OutboundDescriptor
//go:generate mockgen -destination=./mock_logger.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Logger
//go:generate mockgen -destination=./mock_bindable.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Bindable
//go:generate mockgen -destination=./mock_inbound_response.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 InboundResponse
