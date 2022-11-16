package demo

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/context"
	"testing"

	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/worker/v1"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestDemoSparkBuilder(t *testing.T) {

	// spark dependencies initialization
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mailer := NewMockMailer(mockCtrl)
	paymentProvider := NewMockPaymentProvider(mockCtrl)
	inventoryManagementService := NewMockInventoryManagementService(mockCtrl)
	paymentProvider.EXPECT().CreateTransaction(gomock.Any()).Return(Transaction{ID: "uuid", Amount: 50}, nil)
	inventoryManagementService.EXPECT().Reserve([]InventoryItem{{ID: "1", Name: "itemName"}}).Return(nil)
	paymentProvider.EXPECT().ConfirmTransaction(gomock.Any()).Return(nil)

	checkout := NewCheckoutSpark(mailer, paymentProvider, inventoryManagementService)

	// mock handlers initialization
	stageProgressHandler := sdk_v1.NewStageProgressHandler(t)

	variablesHandler := sdk_v1.NewIOHandler(t)
	err := variablesHandler.Output("jobKey",
		&sdk_v1.Variable{Name: "transaction", MimeType: api.MimeTypeJSON, Value: map[string]any{"id": "uuid", "amount": 50}},
		&sdk_v1.Variable{Name: "another", MimeType: api.MimeTypeJSON, Value: map[string]any{"key": "value"}},
		&sdk_v1.Variable{Name: "items", MimeType: api.MimeTypeJSON, Value: []any{map[string]any{"id": "1", "name": "itemName"}}})
	assert.Nil(t, err)

	// get the spark chain
	spark, err := checkout.Spark()
	if err != nil {
		t.Error(err)
		return
	}

	sparkWorker := v1.NewSparkTestWorker(t, spark,
		sdk_v1.WithStageProgressHandler(stageProgressHandler),
		sdk_v1.WithIOHandler(variablesHandler))

	err = sparkWorker.Execute(context.NewSparkMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))

	if err != nil {
		t.Fatal(err)
	}

}
