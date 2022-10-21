package demo

import (
	ctx "context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/worker"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/config"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDemoSparkBuilder(t *testing.T) {

	// spark dependencies initialization
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mailer := NewMockMailer(mockCtrl)
	paymentProvider := NewMockPaymentProvider(mockCtrl)
	inventoryManagementService := NewMockInventoryManagementService(mockCtrl)
	paymentProvider.EXPECT().CreateTransaction(gomock.Any()).Return(Transaction{Id: "uuid", Amount: 50}, nil)
	inventoryManagementService.EXPECT().Reserve([]InventoryItem{{Id: "1", Name: "itemName"}}).Return(nil)
	paymentProvider.EXPECT().ConfirmTransaction(gomock.Any()).Return(nil)

	checkout := NewCheckoutSpark(mailer, paymentProvider, inventoryManagementService)

	cfg, err := config.NewMock(map[string]string{"APP_ENVIRONMENT": "test", "AGENT_SERVER_PORT": "0", "MANAGER_SERVER_PORT": "0"})
	if err != nil {
		t.Error(err)
	}

	// mock handlers initialization
	stageProgressHandler := inmemory.NewStageProgressHandler(t)
	var1, _ := sdk_v1.NewVariable("transaction", "application/json", map[string]any{"id": "uuid", "amount": 50})
	var2, _ := sdk_v1.NewVariable("another", "application/json", map[string]any{"key": "value"})
	var3, _ := sdk_v1.NewVariable("items", "application/json", []any{map[string]any{"id": "1", "name": "itemName"}})
	variablesHandler := inmemory.NewVariableHandler(t,
		sdk_v1.NewSetVariablesRequest("jobKey", "", var1, var2, var3),
	)

	// get the spark chain
	spark, err := checkout.Spark()

	sparkWorker, err := worker.NewSparkWorker(cfg, spark,
		worker.WithStageProgressHandler(stageProgressHandler),
		worker.WithVariableHandler(variablesHandler))
	if err != nil {
		t.Error(err)
	}

	err = sparkWorker.Run(context.NewJobMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", nil))

	if err != nil {
		t.Fatal(err)
	}

}
