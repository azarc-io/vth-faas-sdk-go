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

type bindableInput struct {
	iodp          *ioDataProvider
	correlationID string
	reference     string
	mimeType      string
	data          []byte
}

func (b *bindableInput) Bind(a any) error {
	data, err := b.GetValue()
	if err != nil {
		return err
	}

	// data already cached
	return NewBindable(Value{
		Value:    data,
		MimeType: b.mimeType,
	}).Bind(a)
}

func (b *bindableInput) GetValue() ([]byte, error) {
	if len(b.data) == 0 {
		// first fetch data
		var err error
		if b.data, err = b.iodp.fetchInputData(b.correlationID, b.reference); err != nil {
			return b.data, err
		}
	}

	return b.data, nil
}

func (b *bindableInput) GetMimeType() string {
	return b.mimeType
}

func (b *bindableInput) String() string {
	data, err := b.GetValue()
	if err != nil {
		return err.Error()
	}
	return string(data)
}

func (iodp *ioDataProvider) NewInput(correlationID string, value *BindableValue) Bindable {
	return &bindableInput{
		iodp:          iodp,
		correlationID: correlationID,
		reference:     value.Reference,
		mimeType:      value.MimeType,
	}
}

func (iodp *ioDataProvider) NewOutput(correlationID string, value *BindableValue) (Bindable, error) {
	url, err := iodp.getOutputServiceUrl(correlationID)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url.String(), bytes.NewReader(value.Value))
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
		return nil, fmt.Errorf("error setting output (%d): %s", resp.StatusCode, string(bd))
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

	// override the reference and remove the data
	value.Reference = val.Reference
	value.Value = nil

	return value, nil
}

func (iodp *ioDataProvider) GetStageResult(workflowID, runID, stageName, correlationID string) (Bindable, error) {
	url, err := iodp.getStageResultServiceUrl(workflowID, runID, correlationID, stageName)
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
	url, err := iodp.getStageResultServiceUrl(workflowID, runID, correlationID, stageName)
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

func (iodp *ioDataProvider) fetchInputData(correlationID, reference string) ([]byte, error) {
	url, err := iodp.getInputServiceUrl(correlationID, reference)
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
		return nil, fmt.Errorf(
			"error retrieving input data (%d): correlationID (%s), reference (%s): %s",
			resp.StatusCode,
			correlationID,
			reference,
			string(bd),
		)
	}

	bd, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return bd, nil
}

func (iodp *ioDataProvider) getKey(workflowID, runID, stageName string) string {
	return fmt.Sprintf("%s-%s-%s", workflowID, runID, stageName)
}

func (iodp *ioDataProvider) getStageResultServiceUrl(workflowID, runID, correlationID, stageName string) (*url.URL, error) {
	stageEntryID := iodp.getKey(workflowID, runID, stageName)
	return url.Parse(fmt.Sprintf("%s/%s/%s/%s", iodp.baseUrl, "stage-results", correlationID, stageEntryID))
}

func (iodp *ioDataProvider) getInputServiceUrl(correlationID, reference string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("%s/%s/%s/%s", iodp.baseUrl, "input", correlationID, reference))
}

func (iodp *ioDataProvider) getOutputServiceUrl(correlationID string) (*url.URL, error) {
	return url.Parse(fmt.Sprintf("%s/%s/%s", iodp.baseUrl, "output", correlationID))
}
