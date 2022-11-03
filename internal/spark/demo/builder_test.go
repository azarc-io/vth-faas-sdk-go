package demo

import (
	ctx "context"
	"testing"

	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
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
	stageProgressHandler := inmemory.NewStageProgressHandler(t)

	variablesHandler := inmemory.NewIOHandler(t)
	err := variablesHandler.Output("jobKey",
		&handlers.Variable{Name: "transaction", MimeType: api.MimeTypeJSON, Value: map[string]any{"id": "uuid", "amount": 50}},
		&handlers.Variable{Name: "another", MimeType: api.MimeTypeJSON, Value: map[string]any{"key": "value"}},
		&handlers.Variable{Name: "items", MimeType: api.MimeTypeJSON, Value: []any{map[string]any{"id": "1", "name": "itemName"}}})
	assert.Nil(t, err)

	// get the spark chain
	spark, err := checkout.Spark()
	if err != nil {
		t.Error(err)
		return
	}

	sparkWorker := v1.NewSparkTestWorker(t, spark,
		v1.WithStageProgressHandler(stageProgressHandler),
		v1.WithIOHandler(variablesHandler))

	err = sparkWorker.Execute(context.NewSparkMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))

	if err != nil {
		t.Fatal(err)
	}

}
