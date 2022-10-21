//go:generate mockgen -source=./demo.go -destination=./demo_mocks.go -package=demo github.com/azarc-io/vth-faas-sdk-go/internal/spark/demo
package demo

import (
	"errors"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
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
			Stage("create_payment_transaction", c.CreateTransaction()).
			Stage("reserve_inventory_items", c.ReserveInventoryItems()).
			Complete("confirm_payment_transaction", c.ConfirmPaymentTransaction()).
			Compensate(
				spark.NewNode().
					Stage("cancel_payment_transaction", c.CancelPaymentTransaction(),
						spark.WithStageStatus("create_payment_transaction", sdk_v1.StageStatus_StageCompleted)).
					Stage("restore_inventory_items", c.RestoreInventoryItems(),
						spark.WithStageStatus("create_payment_transaction", sdk_v1.StageStatus_StageCompleted)).
					Stage("send_apologies_email", c.SendApologiesEmail(),
						spark.WithStageStatus("create_payment_transaction", sdk_v1.StageStatus_StageCompleted)).
					Build()).
			Canceled(
				spark.NewNode().
					Stage("send_cancel_email", c.SendCancelEmail()).
					Build()).
			Build()).
		Build()
}

func (c CheckoutSpark) CreateTransaction() api.StageDefinitionFn {
	return func(ctx api.StageContext) (any, api.StageError) {
		variables, err := ctx.GetVariables("", "transaction", "another", "another", "another")
		if err != nil {
			ctx.Log().Error(err, "error getting transaction variable")
			return nil, sdk_errors.NewStageError(err)
		}
		if transactionVariable, ok := variables.Get("transaction"); ok {
			var transaction Transaction
			err = transactionVariable.Bind(&transaction)
			if err != nil {
				ctx.Log().Error(err, "error binding transaction variable")
				return nil, sdk_errors.NewStageError(err)
			}

			transactionCreated, err := c.paymentProvider.CreateTransaction(transaction)

			if err != nil {
				ctx.Log().Info("create_payment_transaction completed")
				return nil, sdk_errors.NewStageError(err, sdk_errors.WithRetry(10, 500))
			}
			return transactionCreated, nil
		}
		return nil, sdk_errors.NewStageError(errors.New("transaction variable not found"))
	}
}

func (c CheckoutSpark) ReserveInventoryItems() api.StageDefinitionFn {
	return func(ctx api.StageContext) (any, api.StageError) {

		itemsVar, err := ctx.GetVariable("", "items")
		if err != nil {
			return nil, sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_Canceled))
		}

		var inventoryItems []InventoryItem
		err = itemsVar.Bind(&inventoryItems)
		if err != nil {
			return nil, sdk_errors.NewStageError(err)
		}

		err = c.inventoryManagementService.Reserve(inventoryItems)

		if err != nil {
			return nil, sdk_errors.NewStageError(err, sdk_errors.WithErrorType(sdk_v1.ErrorType_Canceled))
		}
		return inventoryItems, nil
	}
}

func (c CheckoutSpark) ConfirmPaymentTransaction() api.CompleteDefinitionFn {
	return func(ctx api.CompleteContext) api.StageError {
		result, err := ctx.GetStageResult("create_payment_transaction")
		var transaction Transaction
		err = result.Bind(&transaction)
		if err != nil {
			ctx.Log().Error(err, "error binding transaction variable")
			return sdk_errors.NewStageError(err)
		}

		err = c.paymentProvider.ConfirmTransaction(transaction)

		if err != nil {
			return sdk_errors.NewStageError(err)
		}
		return nil
	}
}

func (c CheckoutSpark) CancelPaymentTransaction() api.StageDefinitionFn {
	return func(ctx api.StageContext) (any, api.StageError) {
		c.paymentProvider.CancelTransaction(Transaction{})
		return nil, nil
	}
}

func (c CheckoutSpark) RestoreInventoryItems() api.StageDefinitionFn {
	return func(ctx api.StageContext) (any, api.StageError) {
		c.inventoryManagementService.RestoreAvailability(nil)
		return nil, nil
	}
}

func (c CheckoutSpark) SendApologiesEmail() api.StageDefinitionFn {
	return func(ctx api.StageContext) (any, api.StageError) {
		c.mailer.SomethingBadHappened()
		return nil, nil
	}
}

func (c CheckoutSpark) SendCancelEmail() api.StageDefinitionFn {
	return func(ctx api.StageContext) (any, api.StageError) {
		c.mailer.Cancellation()
		return nil, nil
	}
}

type CheckoutService interface {
	//STAGES
	CreateTransaction() api.StageDefinitionFn
	ReserveInventoryItems() api.StageDefinitionFn
	//COMPLETE
	ConfirmPaymentTransaction() api.CompleteDefinitionFn
	//COMPENSATE
	CancelPaymentTransaction() api.StageDefinitionFn
	RestoreInventoryItems() api.StageDefinitionFn
	SendApologiesEmail() api.StageDefinitionFn
	//CANCEL
	SendCancelEmail() api.StageDefinitionFn

	Spark() (*spark.Chain, error)
}

type PaymentProvider interface {
	CreateTransaction(transaction Transaction) (Transaction, error)
	ConfirmTransaction(transaction Transaction) error
	CancelTransaction(transaction Transaction) error
}

type Transaction struct {
	Id     string
	Amount float64
}

type InventoryManagementService interface {
	Reserve(inventoryItem []InventoryItem) error
	RestoreAvailability(inventoryItem []InventoryItem) error
}

type InventoryItem struct {
	Id   string
	Name string
}

type Mailer interface {
	Confirmation()
	Cancellation()
	SomethingBadHappened()
}
