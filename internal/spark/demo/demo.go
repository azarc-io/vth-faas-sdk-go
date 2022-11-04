//go:generate mockgen -source=./demo.go -destination=./demo_mocks.go -package=demo github.com/azarc-io/vth-faas-sdk-go/internal/spark/demo
package demo

import (
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1/models"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
)

type CheckoutSpark struct {
	mailer                     Mailer
	paymentProvider            PaymentProvider
	inventoryManagementService InventoryManagementService
}

func NewCheckoutSpark(mailer Mailer, paymentProvider PaymentProvider, inventoryManagementService InventoryManagementService) *CheckoutSpark {
	return &CheckoutSpark{mailer, paymentProvider, inventoryManagementService}
}

func (c CheckoutSpark) Spark() (*spark.Chain, error) {
	return spark.NewChain(
		spark.NewNode().
			Stage("create_payment_transaction", c.CreateTransaction).
			Stage("reserve_inventory_items", c.ReserveInventoryItems).
			Complete("confirm_payment_transaction", c.ConfirmPaymentTransaction).
			Compensate(
				spark.NewNode().
					Stage("cancel_payment_transaction", c.CancelPaymentTransaction,
						spark.WithStageStatus("create_payment_transaction", v1.StageStatus_STAGE_STATUS_COMPLETED)).
					Stage("restore_inventory_items", c.RestoreInventoryItems,
						spark.WithStageStatus("create_payment_transaction", v1.StageStatus_STAGE_STATUS_COMPLETED)).
					Stage("send_apologies_email", c.SendApologiesEmail,
						spark.WithStageStatus("create_payment_transaction", v1.StageStatus_STAGE_STATUS_COMPLETED)).
					Build()).
			Cancelled(
				spark.NewNode().
					Stage("send_cancel_email", c.SendCancelEmail).
					Build()).
			Build()).
		Build()
}

func (c CheckoutSpark) CreateTransaction(ctx v1.StageContext) (any, v1.StageError) {
	inputs := ctx.Inputs("transaction", "another", "another", "another")

	var transaction Transaction
	err := inputs.Get("transaction").Bind(&transaction)
	if err != nil {
		ctx.Log().Error(err, "error binding transaction variable")
		return nil, sdk_errors.NewStageError(err)
	}

	transactionCreated, err := c.paymentProvider.CreateTransaction(transaction)

	if err != nil {
		ctx.Log().Info("create_payment_transaction completed")
		return nil, sdk_errors.NewStageError(err, sdk_errors.WithRetry(10, 500)) //nolint:gomnd
	}
	return transactionCreated, nil
}

func (c CheckoutSpark) ReserveInventoryItems(ctx v1.StageContext) (any, v1.StageError) {
	var inventoryItems []InventoryItem
	err := ctx.Input("items").Bind(&inventoryItems)
	if err != nil {
		return nil, sdk_errors.NewStageError(err)
	}

	err = c.inventoryManagementService.Reserve(inventoryItems)

	if err != nil {
		return nil, sdk_errors.NewStageError(err, sdk_errors.WithErrorType(v1.ErrorType_ERROR_TYPE_CANCELLED))
	}
	return inventoryItems, nil
}

func (c CheckoutSpark) ConfirmPaymentTransaction(ctx v1.CompleteContext) v1.StageError {
	var transaction Transaction
	err := ctx.StageResult("create_payment_transaction").Bind(&transaction)
	if err != nil {
		ctx.Log().Error(err, "error binding transaction variable")
		return sdk_errors.NewStageError(err)
	}
	err = c.paymentProvider.ConfirmTransaction(transaction)
	if err != nil {
		return sdk_errors.NewStageError(err)
	}

	err = ctx.Output(&models.Variable{Name: "newVar", MimeType: api.MimeTypeJSON, Value: "someValue"})
	if err != nil {
		return sdk_errors.NewStageError(err)
	}

	return nil
}

func (c CheckoutSpark) CancelPaymentTransaction(ctx v1.StageContext) (any, v1.StageError) {
	c.paymentProvider.CancelTransaction(Transaction{}) //nolint:errcheck
	return nil, nil
}

func (c CheckoutSpark) RestoreInventoryItems(ctx v1.StageContext) (any, v1.StageError) {
	c.inventoryManagementService.RestoreAvailability(nil) //nolint:errcheck
	return nil, nil
}

func (c CheckoutSpark) SendApologiesEmail(ctx v1.StageContext) (any, v1.StageError) {
	c.mailer.SomethingBadHappened()
	return nil, nil
}

func (c CheckoutSpark) SendCancelEmail(ctx v1.StageContext) (any, v1.StageError) {
	c.mailer.Cancellation()
	return nil, nil
}

type CheckoutService interface {
	CreateTransaction() v1.StageDefinitionFn
	ReserveInventoryItems() v1.StageDefinitionFn
	ConfirmPaymentTransaction() v1.CompleteDefinitionFn
	CancelPaymentTransaction() v1.StageDefinitionFn
	RestoreInventoryItems() v1.StageDefinitionFn
	SendApologiesEmail() v1.StageDefinitionFn
	SendCancelEmail() v1.StageDefinitionFn
	Spark() (*spark.Chain, error)
}

type PaymentProvider interface {
	CreateTransaction(transaction Transaction) (Transaction, error)
	ConfirmTransaction(transaction Transaction) error
	CancelTransaction(transaction Transaction) error
}

type Transaction struct {
	ID     string
	Amount float64
}

type InventoryManagementService interface {
	Reserve(inventoryItem []InventoryItem) error
	RestoreAvailability(inventoryItem []InventoryItem) error
}

type InventoryItem struct {
	ID   string
	Name string
}

type Mailer interface {
	Confirmation()
	Cancellation()
	SomethingBadHappened()
}
