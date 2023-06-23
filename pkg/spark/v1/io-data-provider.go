package sparkv1

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/codec"
	"io"
	"net/http"
	"net/url"
)

type ioDataProvider struct {
	ctx     context.Context
	baseUrl string
	apiKey  string
}

func (iodp *ioDataProvider) GetStageResult(workflowID, runID, stageName, correlationID string) (Bindable, error) {
	url, err := iodp.getServiceUrl(workflowID, runID, correlationID, stageName)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Token", iodp.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bd, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error getting stage result (%d): %s", resp.StatusCode, string(bd))
	}

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var val Value
	err = codec.DecodeAndBind(bd, codec.MimeTypeJson, &val)
	if err != nil {
		return nil, err
	}

	return NewBindable(val), nil
}

func (iodp *ioDataProvider) PutStageResult(workflowID, runID, stageName, correlationID string, stageValue []byte) (Bindable, error) {
	url, err := iodp.getServiceUrl(workflowID, runID, correlationID, stageName)
	if err != nil {
		return nil, err
	}

	d, _ := json.Marshal(Value{
		Value:    stageValue,
		MimeType: "application/json",
	})

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(d))
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Token", iodp.apiKey)
	req.Header.Set("Content-Type", string(codec.MimeTypeJson))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bd, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("error putting stage result (%d): %s", resp.StatusCode, string(bd))
	}

	val, _ := codec.Encode(iodp.getKey(workflowID, runID, stageName))
	return NewBindable(Value{
		Value:    val,
		MimeType: string(codec.MimeTypeText),
	}), nil
}

func (iodp *ioDataProvider) getKey(workflowID, runID, stageName string) string {
	return fmt.Sprintf("%s-%s-%s", workflowID, runID, stageName)
}

func (iodp *ioDataProvider) getServiceUrl(workflowID, runID, correlationID, stageName string) (*url.URL, error) {
	stageEntryID := iodp.getKey(workflowID, runID, stageName)
	return url.Parse(fmt.Sprintf("%s/%s/%s/%s", iodp.baseUrl, "stage-results", correlationID, stageEntryID))
}
