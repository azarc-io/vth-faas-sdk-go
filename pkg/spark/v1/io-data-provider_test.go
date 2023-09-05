package sparkv1

import (
	"context"
	helpers "github.com/azarc-io/vth-faas-sdk-go/internal/test_helpers"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestIoDataProvider(t *testing.T) {
	t.Run("Stage IO", func(t *testing.T) {
		entryData := []byte(`{"value":"ImhlbGxvIHdvcmxkIg==","mime_type":"application/json"}`)
		svr := helpers.GetTestHttpServerWithRequests(t, []helpers.Request{
			{http.MethodPost, 200, "/stage-results/c789/w123-r456-dummy-stage", []byte("foo bar"), entryData, func(t *testing.T, req *http.Request) {
				assert.Equal(t, codec.MimeTypeJson, codec.MimeType(req.Header.Get("Content-Type")))
				assert.Equal(t, "dummy-token", req.Header.Get("X-Token"))
			}},
			{http.MethodPost, 500, "/stage-results/c789/w123-r456-dummy-error", []byte("foo bar error"), nil, nil},
			{http.MethodGet, 200, "/stage-results/c789/w123-r456-dummy-stage", entryData, nil, func(t *testing.T, req *http.Request) {
				assert.Equal(t, "dummy-token", req.Header.Get("X-Token"))
			}},
			{http.MethodGet, 500, "/stage-results/c789/w123-r456-dummy-error", []byte("foo bar error"), nil, nil},
		})

		iop := &ioDataProvider{
			ctx:     context.Background(),
			baseUrl: svr.URL,
			apiKey:  "dummy-token",
		}
		input, _ := codec.Encode("hello world")

		t.Run("put stage result: success", func(t *testing.T) {
			sr, err := iop.PutStageResult("w123", "r456", "dummy-stage", "c789", input)
			assert.NoError(t, err)
			assert.Equal(t, "w123-r456-dummy-stage", sr.String())
		})
		t.Run("put stage result: fail", func(t *testing.T) {
			_, err := iop.PutStageResult("w123", "r456", "dummy-error", "c789", input)
			assert.ErrorContains(t, err, "error putting stage result (500): foo bar error")
		})

		t.Run("get stage result: success", func(t *testing.T) {
			sr, err := iop.GetStageResult("w123", "r456", "dummy-stage", "c789")
			assert.NoError(t, err)

			var res any
			assert.NoError(t, sr.Bind(&res))
			assert.Equal(t, "hello world", res)
		})
		t.Run("get stage result: fail", func(t *testing.T) {
			_, err := iop.GetStageResult("w123", "r456", "dummy-error", "c789")
			assert.ErrorContains(t, err, "error getting stage result (500): foo bar error")
		})
	})

	t.Run("Spark IO", func(t *testing.T) {
		entryData := []byte(`{"reference":"my output reference"}`)
		svr := helpers.GetTestHttpServerWithRequests(t, []helpers.Request{
			{http.MethodGet, 200, "/input/c789/my-input-ref", []byte(`"foo bar"`), nil, func(t *testing.T, req *http.Request) {
				assert.Equal(t, "dummy-token", req.Header.Get("X-Token"))
			}},
			{http.MethodPost, 200, "/output/c789", entryData, []byte("some data"), func(t *testing.T, req *http.Request) {
				assert.Equal(t, "dummy-token", req.Header.Get("X-Token"))
			}},
			{http.MethodGet, 400, "/input/c567/my-input-ref", []byte("foo bar"), nil, func(t *testing.T, req *http.Request) {
				assert.Equal(t, "dummy-token", req.Header.Get("X-Token"))
			}},
			{http.MethodPost, 400, "/output/c567", entryData, []byte("some data"), func(t *testing.T, req *http.Request) {
				assert.Equal(t, "dummy-token", req.Header.Get("X-Token"))
			}},
		})

		iop := &ioDataProvider{
			ctx:     context.Background(),
			baseUrl: svr.URL,
			apiKey:  "dummy-token",
		}

		t.Run("get input: success", func(t *testing.T) {
			sr := iop.NewInput("c789", &BindableValue{
				MimeType:  string(codec.MimeTypeJson),
				Reference: "my-input-ref",
			})

			assert.Equal(t, "foo bar", sr.String())
		})
		t.Run("get input: fail", func(t *testing.T) {
			sr := iop.NewInput("c567", &BindableValue{
				MimeType:  string(codec.MimeTypeJson),
				Reference: "my-input-ref",
			})

			var a any
			err := sr.Bind(&a)
			assert.ErrorContains(t, err, "error retrieving input data (400): correlationID (c567), reference (my-input-ref): foo bar")
			assert.Nil(t, a)
		})

		t.Run("post output: success", func(t *testing.T) {
			sr, err := iop.NewOutput("c789", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("some data"),
			})

			assert.NoError(t, err)
			assert.Empty(t, sr.String())
			assert.Equal(t, "my output reference", sr.(*BindableValue).Reference)
		})

		t.Run("post output: fail", func(t *testing.T) {
			sr, err := iop.NewOutput("c567", &BindableValue{
				MimeType: string(codec.MimeTypeJson),
				Value:    []byte("some data"),
			})

			assert.ErrorContains(t, err, "error setting output (400)")
			assert.Nil(t, sr)
		})
	})
}
