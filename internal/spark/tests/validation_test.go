package tests

import (
	_ "embed"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	//go:embed testdata/err_msg_name_stage1_not_unique
	errMsgNameStage1NotUnique string
	//go:embed testdata/err_msg_no_stage_on_root
	errMsgNoStagesOnRoot string
	//go:embed testdata/err_msg_inner_nodes_same_stage_name
	errMsgInnerNodesSameStageName string
	//go:embed testdata/err_msg_empty_stage_name
	errMsgEmptyStageName string
)

func Test(t *testing.T) {
	tests := []struct {
		name             string
		chainFn          func() (*spark.Chain, error)
		expectedErrorMsg string
	}{
		{
			name: "should return no validation errors",
			chainFn: func() (*spark.Chain, error) {
				return spark.NewChain(
					spark.NewNode().
						Stage("stage1", noOpStage).
						Stage("stage2", noOpStage).
						Stage("stage3", noOpStage).
						Complete("complete", noOpComplete).
						Compensate(spark.NewNode().Stage("compensate", noOpStage).Build()).
						Cancelled(spark.NewNode().Stage("canceled", noOpStage).Build()).
						Build()).
					Build()
			},
		},
		{
			name: "should return validation error: same name used in two stages",
			chainFn: func() (*spark.Chain, error) {
				return spark.NewChain(
					spark.NewNode().
						Stage("stage1", noOpStage).
						Stage("stage1", noOpStage).
						Stage("stage3", noOpStage).
						Complete("complete", noOpComplete).
						Compensate(spark.NewNode().Stage("compensate", noOpStage).Build()).
						Cancelled(spark.NewNode().Stage("canceled", noOpStage).Build()).
						Build()).
					Build()
			},
			expectedErrorMsg: errMsgNameStage1NotUnique,
		},
		{
			name: "should return validation error: node without stages",
			chainFn: func() (*spark.Chain, error) {
				return spark.NewChain(
					spark.NewNode().
						Complete("complete", noOpComplete).
						Compensate(spark.NewNode().Stage("compensate", noOpStage).Build()).
						Cancelled(spark.NewNode().Stage("canceled", noOpStage).Build()).
						Build()).
					Build()
			},
			expectedErrorMsg: errMsgNoStagesOnRoot,
		},
		{
			name: "should return validation error: multiple inner stages with the same name",
			chainFn: func() (*spark.Chain, error) {
				return spark.NewChain(
					spark.NewNode().Stage("stage1", noOpStage).
						Cancelled(spark.NewNode().Stage("canceled", noOpStage).
							Cancelled(spark.NewNode().Stage("canceled", noOpStage).
								Cancelled(spark.NewNode().Stage("canceled", noOpStage).
									Cancelled(spark.NewNode().Stage("canceled", noOpStage).Build()).
									Build(),
								).Build()).
							Build(),
						).Build()).
					Build()
			},
			expectedErrorMsg: errMsgInnerNodesSameStageName,
		},
		{
			name: "should return validation error: stage with empty name",
			chainFn: func() (*spark.Chain, error) {
				return spark.NewChain(
					spark.NewNode().
						Stage("stage1", noOpStage).
						Stage("stage2", noOpStage).
						Stage("stage3", noOpStage).
						Complete("", noOpComplete).
						Compensate(spark.NewNode().Stage("compensate", noOpStage).Build()).
						Cancelled(spark.NewNode().Stage("canceled", noOpStage).Build()).
						Build()).
					Build()
			},
			expectedErrorMsg: errMsgEmptyStageName,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := test.chainFn()
			if test.expectedErrorMsg != "" {
				if err == nil {
					t.Errorf("error expected: %s, got: <nil>", test.expectedErrorMsg)
					return
				}
				assert.Equal(t, test.expectedErrorMsg, err.Error())
			} else if err != nil {
				t.Errorf("unexpected error: %s", err.Error())
			}
		})
	}
}

var noOpStage = func(ctx sdk_v1.StageContext) (any, sdk_v1.StageError) { return nil, nil }
var noOpComplete = func(ctx sdk_v1.CompleteContext) sdk_v1.StageError { return nil }
