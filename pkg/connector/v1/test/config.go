package test

import (
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"strings"
)

type Config struct {
	UserConfig         []byte              `json:"user_config"`
	Ingress            []IngressConfig     `json:"ingress_config"`
	InboundDescriptors []InboundDescriptor `json:"inbound_descriptors"`
}

type IngressConfig struct {
	Name       string `json:"name"`
	BindHost   string `json:"host"`
	BindPort   int    `json:"port"`
	ExtAddress string `json:"external_address"`
}

func (i IngressConfig) ExternalAddress() string {
	return i.ExtAddress
}

func (i IngressConfig) InternalPort() int {
	return i.BindPort
}

func (i IngressConfig) InternalHost() string {
	return i.BindHost
}

type InboundDescriptor struct {
	ID           string                  `json:"id"`
	ReadableName string                  `json:"name"`
	MsgName      string                  `json:"message_name"`
	Mime         string                  `json:"mime_type"`
	Type         connectorv1.MessageType `json:"type"`
	Options      []byte                  `json:"options"`
}

func (m InboundDescriptor) Name() string {
	return m.ReadableName
}

func (m InboundDescriptor) MessageName() string {
	return m.MsgName
}

func (m InboundDescriptor) MimeType() string {
	return m.Mime
}

func (m InboundDescriptor) MessageType() connectorv1.MessageType {
	return m.Type
}

func (m InboundDescriptor) Config() connectorv1.Bindable {
	var tp string
	if m.Mime != "" {
		parts := strings.Split(m.Mime, "/")
		tp = parts[len(parts)-1]
	}
	return connectorv1.NewBindable(m.Options, connectorv1.BindableType(tp))
}
