package tests

// TODO Move tests to /pkg/spark and get them working again
//import (
//	_ "embed"
//	"testing"
//
//	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
//	"github.com/stretchr/testify/assert"
//)
//
//var (
//	//go:embed testdata/err_msg_name_stage1_not_unique
//	errMsgNameStage1NotUnique string
//	//go:embed testdata/err_msg_no_stage_on_root
//	errMsgNoStagesOnRoot string
//	//go:embed testdata/err_msg_inner_nodes_same_stage_name
//	errMsgInnerNodesSameStageName string
//	//go:embed testdata/err_msg_empty_stage_name
//	errMsgEmptyStageName string
//)
//
//func Test(t *testing.T) {
//	tests := []struct {
//		name             string
//		chainFn          func() (*sdk_v1.BuilderChain, error)
//		expectedErrorMsg string
//	}{
//		{
//			name: "should return no validation errors",
//			chainFn: func() (*sdk_v1.BuilderChain, error) {
//				return sdk_v1.NewChain(
//					sdk_v1.NewNode().
//						Stage("stage1", noOpStage).
//						Stage("stage2", noOpStage).
//						Stage("stage3", noOpStage).
//						Complete("complete", noOpComplete).
//						Compensate(sdk_v1.NewNode().Stage("compensate", noOpStage).Build()).
//						Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).Build()).
//						Build()).
//					Build()
//			},
//		},
//		{
//			name: "should return validation error: same name used in two stages",
//			chainFn: func() (*sdk_v1.BuilderChain, error) {
//				return sdk_v1.NewChain(
//					sdk_v1.NewNode().
//						Stage("stage1", noOpStage).
//						Stage("stage1", noOpStage).
//						Stage("stage3", noOpStage).
//						Complete("complete", noOpComplete).
//						Compensate(sdk_v1.NewNode().Stage("compensate", noOpStage).Build()).
//						Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).Build()).
//						Build()).
//					Build()
//			},
//			expectedErrorMsg: errMsgNameStage1NotUnique,
//		},
//		{
//			name: "should return validation error: node without stages",
//			chainFn: func() (*sdk_v1.BuilderChain, error) {
//				return sdk_v1.NewChain(
//					sdk_v1.NewNode().
//						Complete("complete", noOpComplete).
//						Compensate(sdk_v1.NewNode().Stage("compensate", noOpStage).Build()).
//						Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).Build()).
//						Build()).
//					Build()
//			},
//			expectedErrorMsg: errMsgNoStagesOnRoot,
//		},
//		{
//			name: "should return validation error: multiple inner stages with the same name",
//			chainFn: func() (*sdk_v1.BuilderChain, error) {
//				return sdk_v1.NewChain(
//					sdk_v1.NewNode().Stage("stage1", noOpStage).
//						Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).
//							Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).
//								Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).
//									Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).Build()).
//									Build(),
//								).Build()).
//							Build(),
//						).Build()).
//					Build()
//			},
//			expectedErrorMsg: errMsgInnerNodesSameStageName,
//		},
//		{
//			name: "should return validation error: stage with empty name",
//			chainFn: func() (*sdk_v1.BuilderChain, error) {
//				return sdk_v1.NewChain(
//					sdk_v1.NewNode().
//						Stage("stage1", noOpStage).
//						Stage("stage2", noOpStage).
//						Stage("stage3", noOpStage).
//						Complete("", noOpComplete).
//						Compensate(sdk_v1.NewNode().Stage("compensate", noOpStage).Build()).
//						Cancelled(sdk_v1.NewNode().Stage("canceled", noOpStage).Build()).
//						Build()).
//					Build()
//			},
//			expectedErrorMsg: errMsgEmptyStageName,
//		},
//	}
//	for _, test := range tests {
//		t.Run(test.name, func(t *testing.T) {
//			_, err := test.chainFn()
//			if test.expectedErrorMsg != "" {
//				if err == nil {
//					t.Errorf("error expected: %s, got: <nil>", test.expectedErrorMsg)
//					return
//				}
//				assert.Equal(t, test.expectedErrorMsg, err.Error())
//			} else if err != nil {
//				t.Errorf("unexpected error: %s", err.Error())
//			}
//		})
//	}
//}
//
//var noOpStage = func(ctx sdk_v1.StageContext) (any, sdk_v1.StageError) { return nil, nil }
//var noOpComplete = func(ctx sdk_v1.CompleteContext) sdk_v1.StageError { return nil }
