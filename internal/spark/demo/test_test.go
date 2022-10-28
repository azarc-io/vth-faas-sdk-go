package demo

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/golang/mock/gomock"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPaymentTransaction(t *testing.T) {
	mailer, provider, service, controller := createTestMocks(t)
	defer controller.Finish()

	provider.EXPECT().CreateTransaction(gomock.Any()).Return(Transaction{Id: "uuid", Amount: 50}, nil)
	service.EXPECT().Reserve([]InventoryItem{{Id: "1", Name: "itemName"}}).Return(nil)
	provider.EXPECT().ConfirmTransaction(gomock.Any()).Return(nil)

	worker, sp, io := createTestSpark(t, mailer, provider, service)

	err := io.Output("jobKey",
		&handlers.Variable{Name: "transaction", MimeType: api.MimeTypeJson, Value: map[string]any{"id": "uuid", "amount": 50}},
		&handlers.Variable{Name: "another", MimeType: api.MimeTypeJson, Value: map[string]any{"key": "value"}},
		&handlers.Variable{Name: "items", MimeType: api.MimeTypeJson, Value: []any{map[string]any{"id": "1", "name": "itemName"}}})
	assert.Nil(t, err)

	err = worker.Execute(context.NewJobMetadata(ctx.Background(),
		"jobKey", "correlationId", "transactionId", nil))

	assert.Nil(t, err)

	raw, err := io.Input("jobKey", "newVar").Raw()
	assert.Nil(t, err)
	assert.Equal(t, string(raw), `"someValue"`)

	stage1Status, err := sp.Get("jobKey", "confirm_payment_transaction")
	assert.Nil(t, err)
	assert.Equal(t, lo.ToPtr(sdk_v1.StageStatus_StageCompleted), stage1Status)
}

func createTestMocks(t *testing.T) (*MockMailer, *MockPaymentProvider, *MockInventoryManagementService, *gomock.Controller) {
	mockCtrl := gomock.NewController(t)
	return NewMockMailer(mockCtrl), NewMockPaymentProvider(mockCtrl), NewMockInventoryManagementService(mockCtrl), mockCtrl
}

func createTestSpark(t *testing.T,
	mailer *MockMailer,
	provider *MockPaymentProvider,
	service *MockInventoryManagementService) (sdk_v1.Worker, sdk_v1.StageProgressHandler, sdk_v1.IOHandler) {

	checkoutSpark := NewCheckoutSpark(mailer, provider, service)

	spark, err := checkoutSpark.Spark()
	assert.Nil(t, err)

	stageProgressHandler := inmemory.NewStageProgressHandler(t)

	variablesHandler := inmemory.NewIOHandler(t)

	sparkWorker := v1.NewSparkTestWorker(t, spark,
		v1.WithStageProgressHandler(stageProgressHandler),
		v1.WithIOHandler(variablesHandler))

	return sparkWorker, stageProgressHandler, variablesHandler
}
