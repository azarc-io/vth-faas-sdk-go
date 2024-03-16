package sparkv1

import (
	"context"
	"encoding/json"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1/util"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIoDataProvider(t *testing.T) {
	port, err := util.GetFreeTCPPort()
	if err != nil {
		t.Fatal(err)
	}

	s, err := util.RunServerOnPort(port, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	defer s.Shutdown()
	s.Start()

	nc, js := util.GetNatsClient(port)
	defer nc.Close()

	store, err := js.CreateObjectStore(context.Background(), jetstream.ObjectStoreConfig{
		Bucket: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Stage IO", func(t *testing.T) {
		iop := &ioDataProvider{
			ctx:          context.Background(),
			stageResults: map[string]*BindableValue{},
			store:        store,
		}
		input, _ := codec.Encode("hello world")

		t.Run("put stage result: success", func(t *testing.T) {
			sr, err := iop.PutStageResult("dummy-stage", input)
			assert.NoError(t, err)
			assert.Equal(t, "hello world", sr.String())
		})

		t.Run("get stage result: success", func(t *testing.T) {
			sr, err := iop.GetStageResult("dummy-stage")
			assert.NoError(t, err)

			var res any
			assert.NoError(t, sr.Bind(&res))
			assert.Equal(t, "hello world", res)
		})

		t.Run("get stage result: fail", func(t *testing.T) {
			_, err := iop.GetStageResult("dummy-error")
			assert.ErrorContains(t, err, "stage result not found")
		})
	})

	t.Run("Spark IO", func(t *testing.T) {
		iop := &ioDataProvider{
			ctx:          context.Background(),
			stageResults: map[string]*BindableValue{},
			store:        store,
		}

		t.Run("get input: success", func(t *testing.T) {
			b, err := json.Marshal(map[string]any{
				"test": &BindableValue{
					MimeType: string(codec.MimeTypeJson),
					Value:    []byte(`"foo bar"`),
				},
			})
			if err != nil {
				t.Fatal(err)
			}

			if _, err := store.PutBytes(context.Background(), "test", b); err != nil {
				t.Fatal(err)
			}

			if err := iop.LoadVariables("test"); err != nil {
				t.Fatal(err)
			}

			sr := iop.NewInput("test", "c789", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
			})

			assert.Equal(t, "foo bar", sr.String())
		})

		t.Run("get input: fail", func(t *testing.T) {
			sr := iop.NewInput("missing", "c567", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("foo bar"),
			})

			var a any
			err := sr.Bind(&a)
			assert.ErrorContains(t, err, "variable not found: error retrieving input data: stage (c567), name (missing)")
			assert.Nil(t, a)
		})

		t.Run("post output: success", func(t *testing.T) {
			sr, err := iop.NewOutput("c789", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("some data"),
			})

			assert.NoError(t, err)
			assert.Empty(t, sr.String())
			assert.Equal(t, []byte("some data"), sr.(*BindableValue).Value)
		})

		t.Run("post output: fail", func(t *testing.T) {
			t.Skipf("no longer fails because outputs are stored in memory until the spark completed")
			sr, err := iop.NewOutput("c567", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("some data"),
			})

			assert.ErrorContains(t, err, "error setting output (400)")
			assert.Nil(t, sr)
		})
	})
}
