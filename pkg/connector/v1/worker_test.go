package connectorv1_test

import (
	"fmt"
	connectorv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/connector/v1/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"os"
	"syscall"
	"testing"
	"time"
)

type forwardResponse struct {
	HeadersMap connectorv1.Headers `json:"headers"`
	Payload    []byte              `json:"payload"`
}

func (f forwardResponse) Body() connectorv1.Bindable {
	return connectorv1.NewBindable(f.Payload, connectorv1.BindableTypeJson)
}

func (f forwardResponse) Headers() connectorv1.Headers {
	return f.HeadersMap
}

func TestWorker(t *testing.T) {
	ctrl := gomock.NewController(t)
	connector := mock.NewMockConnector(ctrl)
	forwarder := mock.NewMockForwarder(ctrl)

	_ = os.Setenv("CONNECTOR_FILE_PATH", "./fixtures/connector_config_1.yaml")
	_ = os.Setenv("INBOUND_DESCRIPTOR_FILE_PATH", "./fixtures/inbound_descriptors_config_1.yaml")
	_ = os.Setenv("CONFIG_FILE_PATH", "./fixtures/user_config_1.yaml")

	waitChan := make(chan struct{}, 1)
	connector.EXPECT().Start(gomock.Any()).DoAndReturn(func(ctx connectorv1.StartContext) error {
		defer func() {
			waitChan <- struct{}{}
		}()
		t.Run("connector/config", func(t *testing.T) {
			var conf struct {
				Name  string `yaml:"name"`
				Value string `yaml:"value"`
			}
			err := ctx.Config().Bind(&conf)
			assert.NoError(t, err)
			assert.Equal(t, "user config", conf.Name)
			assert.Equal(t, "user-defined config", conf.Value)
		})

		t.Run("connector/ingress", func(t *testing.T) {
			ingr, err := ctx.Ingress("http-8080")
			assert.NoError(t, err)
			assert.NotNil(t, ingr)
			assert.Equal(t, 8080, ingr.InternalPort())
			assert.Equal(t, "0.0.0.0", ingr.InternalHost())
			assert.Equal(t, "https://some-external-url.com:443/v1/endpoint/12356", ingr.ExternalAddress())

		})
		t.Run("connector/inbound_descriptors", func(t *testing.T) {
			descriptors := ctx.InboundDescriptors()
			assert.Equal(t, 2, len(descriptors))
			mimeTypes := []string{"application/json", "application/yaml"}
			for i, desc := range descriptors {
				index := fmt.Sprint(i + 1)
				assert.Equal(t, "User friendly name "+index, desc.Name())
				assert.Equal(t, "message-name-"+index, desc.MessageName())
				assert.Equal(t, connectorv1.MessageTypeInbound, desc.MessageType())
				assert.Equal(t, mimeTypes[i], desc.MimeType())
				var messageConfig struct {
					Key string `json:"test-key" yaml:"test-key"`
				}
				err := desc.Config().Bind(&messageConfig)
				assert.NoError(t, err)
				assert.Equal(t, "test-value-"+index, messageConfig.Key)
			}
		})

		t.Run("connector/forwarder", func(t *testing.T) {
			resp, err := ctx.Forwarder().Forward("message-name-1", []byte("some-data"), map[string]string{
				"key": "value",
			})
			assert.NoError(t, err)
			var respBody struct {
				Response string `json:"response"`
			}
			err = resp.Body().Bind(&respBody)
			assert.NoError(t, err)
			assert.Equal(t, "response-body", respBody.Response)
			assert.Equal(t, "header-value", resp.Headers()["key"])
		})
		return nil
	})

	forwarder.EXPECT().Forward("message-name-1", []byte("some-data"), map[string]string{
		"key": "value",
	}).Return(forwardResponse{
		Payload: []byte(`{"response": "response-body"}`),
		HeadersMap: connectorv1.Headers{
			"key": "header-value",
		},
	}, nil)

	connector.EXPECT().Stop(gomock.Any()).Return(nil)

	worker, err := connectorv1.NewConnectorWorker(connector, connectorv1.WithForwarder(forwarder))
	assert.NoError(t, err)

	go func() {
		select {
		case <-waitChan:
		case <-time.After(10 * time.Second):
			t.Fail()
		}
		_ = syscall.Kill(syscall.Getpid(), syscall.SIGINT)
	}()
	worker.Run()
}
