package mock

import _ "github.com/golang/mock/mockgen/model"

//go:generate mockgen -destination=./mock_forwarder.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Forwarder
//go:generate mockgen -destination=./mock_connector.go -package mock github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1 Connector
