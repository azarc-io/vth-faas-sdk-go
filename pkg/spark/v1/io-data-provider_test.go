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
}
