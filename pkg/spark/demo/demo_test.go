package demo

//import (
//	ctx "context"
//	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/context"
//	"testing"
//
//	v12 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
//
//	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
//	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/worker/v1"
//	"github.com/golang/mock/gomock"
//	"github.com/samber/lo"
//	"github.com/stretchr/testify/assert"
//)
//
//func TestPaymentTransaction(t *testing.T) {
//	mailer, provider, service, controller := createTestMocks(t)
//	defer controller.Finish()
//
//	provider.EXPECT().CreateTransaction(gomock.Any()).Return(Transaction{ID: "uuid", Amount: 50}, nil)
//	service.EXPECT().Reserve([]InventoryItem{{ID: "1", Name: "itemName"}}).Return(nil)
//	provider.EXPECT().ConfirmTransaction(gomock.Any()).Return(nil)
//
//	worker, sp, io := createTestSpark(t, mailer, provider, service)
//
//	err := io.Output("jobKey",
//		&v12.Variable{Name: "transaction", MimeType: api.MimeTypeJSON, Value: map[string]any{"id": "uuid", "amount": 50}},
//		&v12.Variable{Name: "another", MimeType: api.MimeTypeJSON, Value: map[string]any{"key": "value"}},
//		&v12.Variable{Name: "items", MimeType: api.MimeTypeJSON, Value: []any{map[string]any{"id": "1", "name": "itemName"}}})
//	assert.Nil(t, err)
//
//	err = worker.Execute(context.NewSparkMetadata(ctx.Background(),
//		"jobKey", "correlationId", "transactionId", nil))
//
//	assert.Nil(t, err)
//
//	raw, err := io.Input("jobKey", "newVar").Raw()
//	assert.Nil(t, err)
//	assert.Equal(t, string(raw), `"someValue"`)
//
//	stage1Status, err := sp.Get("jobKey", "confirm_payment_transaction")
//	assert.Nil(t, err)
//	assert.Equal(t, lo.ToPtr(v12.StageStatus_STAGE_STATUS_COMPLETED), stage1Status)
//}
//
//func createTestMocks(t *testing.T) (*MockMailer, *MockPaymentProvider, *MockInventoryManagementService, *gomock.Controller) {
//	mockCtrl := gomock.NewController(t)
//	return NewMockMailer(mockCtrl), NewMockPaymentProvider(mockCtrl), NewMockInventoryManagementService(mockCtrl), mockCtrl
//}
//
//func createTestSpark(t *testing.T,
//	mailer *MockMailer,
//	provider *MockPaymentProvider,
//	service *MockInventoryManagementService) (v12.Worker, v12.StageProgressHandler, v12.IOHandler) {
//
//	checkoutSpark := NewCheckoutSpark(mailer, provider, service)
//
//	spark, err := checkoutSpark.Spark()
//	assert.Nil(t, err)
//
//	stageProgressHandler := v12.NewStageProgressHandler(t)
//
//	variablesHandler := v12.NewIOHandler(t)
//
//	sparkWorker := v1.NewSparkTestWorker(t, spark,
//		v12.WithStageProgressHandler(stageProgressHandler),
//		v12.WithIOHandler(variablesHandler))
//
//	return sparkWorker, stageProgressHandler, variablesHandler
//}
