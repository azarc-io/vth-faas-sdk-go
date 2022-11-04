package tests

import (
	ctx "context"
	"errors"
	"testing"

	"github.com/azarc-io/vth-faas-sdk-go/internal/context"
	"github.com/azarc-io/vth-faas-sdk-go/internal/handlers/test/inmemory"
	"github.com/azarc-io/vth-faas-sdk-go/internal/spark"
	v1 "github.com/azarc-io/vth-faas-sdk-go/internal/worker/v1"
	v12 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
	sdk_errors "github.com/azarc-io/vth-faas-sdk-go/pkg/errors"
	"github.com/samber/lo"
)

func TestResumeOnRetry(t *testing.T) {
	newSB := func() *stageBehaviour {
		return NewStageBehaviour(t, "stage1", "stage2", "stage3", "complete", "compensate", "canceled")
	}

	tests := []struct {
		name            string
		stageBehaviour  *stageBehaviour
		lastActiveStage *v12.LastActiveStage
		assertFn        func(t *testing.T, sb *stageBehaviour)
		errorType       *v12.ErrorType
	}{
		{
			name:            "should execute all stages and complete",
			stageBehaviour:  newSB(),
			lastActiveStage: nil,
			assertFn: func(t *testing.T, sb *stageBehaviour) {
				for _, stage := range []string{"stage1", "stage2", "stage3", "complete"} {
					if !sb.Executed(stage) {
						t.Errorf("stage '%s' expected to be executed", stage)
					}
				}
			},
			errorType: nil,
		},
		{
			name:           "should execute only complete stage",
			stageBehaviour: newSB(),
			lastActiveStage: &v12.LastActiveStage{
				Name: "complete",
			},
			assertFn: func(t *testing.T, sb *stageBehaviour) {
				for _, stage := range []string{"stage1", "stage2", "stage3", "canceled", "compensate"} {
					if sb.Executed(stage) {
						t.Errorf("stage '%s' expected to not be executed", stage)
					}
				}
				if !sb.Executed("complete") {
					t.Error("stage 'complete' expected to be executed")
				}
			},
			errorType: nil,
		},
		{
			name:           "should execute stage3 and complete stages",
			stageBehaviour: newSB(),
			lastActiveStage: &v12.LastActiveStage{
				Name: "stage3",
			},
			assertFn: func(t *testing.T, sb *stageBehaviour) {
				for _, stage := range []string{"stage1", "stage2", "canceled", "compensate"} {
					if sb.Executed(stage) {
						t.Errorf("stage '%s' expected to not be executed", stage)
					}
				}
				for _, stage := range []string{"stage3", "complete"} {
					if !sb.Executed(stage) {
						t.Errorf("stage '%s' expected to be executed", stage)
					}
				}
			},
			errorType: nil,
		},
		{
			name:           "should execute only stage2 and compensate",
			stageBehaviour: newSB().Change("stage2", nil, sdk_errors.NewStageError(errors.New("stage2 error"))),
			lastActiveStage: &v12.LastActiveStage{
				Name: "stage2",
			},
			assertFn: func(t *testing.T, sb *stageBehaviour) {
				for _, stage := range []string{"stage1", "stage3", "complete", "canceled"} {
					if sb.Executed(stage) {
						t.Errorf("stage '%s' expected to not be executed", stage)
					}
				}
				for _, stage := range []string{"stage2", "compensate"} {
					if !sb.Executed(stage) {
						t.Errorf("stage '%s' expected to not be executed", stage)
					}
				}
			},
			errorType: lo.ToPtr(v12.ErrorType_ERROR_TYPE_FAILED_UNSPECIFIED),
		},
		{
			name: "should execute only stage2 and cancel",
			stageBehaviour: newSB().
				Change("stage2", nil, sdk_errors.NewStageError(errors.New("stage2 cancel"), sdk_errors.WithErrorType(v12.ErrorType_ERROR_TYPE_CANCELLED))),
			lastActiveStage: &v12.LastActiveStage{
				Name: "stage2",
			},
			assertFn: func(t *testing.T, sb *stageBehaviour) {
				for _, stage := range []string{"stage1", "stage3", "complete", "compensate"} {
					if sb.Executed(stage) {
						t.Errorf("stage '%s' expected to not be executed", stage)
					}
				}
				for _, stage := range []string{"stage2", "canceled"} {
					if !sb.Executed(stage) {
						t.Errorf("stage '%s' expected to be executed", stage)
					}
				}
			},
			errorType: lo.ToPtr(v12.ErrorType_ERROR_TYPE_CANCELLED),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.stageBehaviour.ResetExecutions()
			chain := createChainForResumeOnRetryTests(t, test.stageBehaviour)
			worker := v1.NewSparkTestWorker(t, chain, v1.WithIOHandler(inmemory.NewIOHandler(t)), v1.WithStageProgressHandler(inmemory.NewStageProgressHandler(t)))
			err := worker.Execute(context.NewSparkMetadata(ctx.Background(), "jobKey", "correlationId", "transactionId", test.lastActiveStage))
			if err != nil && test.errorType == nil {
				t.Errorf("a unexpected error occured: %v", err)
			}
			if test.errorType != nil {
				if err == nil {
					t.Errorf("error '%s' is expected from chain execution, got none", test.errorType)
				} else if *test.errorType != err.ErrorType() {
					t.Errorf("error expected: %v; got: %v;", test.errorType, err.ErrorType())
				}
			}
			test.assertFn(t, test.stageBehaviour)
		})
	}
}

func createChainForResumeOnRetryTests(t *testing.T, sb *stageBehaviour) *spark.Chain {
	chain, err := spark.NewChain(
		spark.NewNode().
			Stage("stage1", stageFn("stage1", sb)).
			Stage("stage2", stageFn("stage2", sb)).
			Stage("stage3", stageFn("stage3", sb)).
			Compensate(
				spark.NewNode().Stage("compensate", stageFn("compensate", sb)).Build(),
			).
			Cancelled(
				spark.NewNode().Stage("canceled", stageFn("canceled", sb)).Build(),
			).
			Complete("complete", func(context v12.CompleteContext) v12.StageError {
				sb.exec("complete")
				return sb.shouldErr("complete")
			}).Build(),
	).Build()
	if err != nil {
		t.Fatal(err)
		return nil
	}
	return chain
}

func stageFn(name string, sm *stageBehaviour) v12.StageDefinitionFn {
	return func(context v12.StageContext) (any, v12.StageError) {
		sm.exec(name)
		return sm.shouldReturn(name), sm.shouldErr(name)
	}
}

type behaviour struct {
	executed bool
	result   any
	err      v12.StageError
}

type stageBehaviour struct {
	t *testing.T
	m map[string]*behaviour
}

func NewStageBehaviour(t *testing.T, stages ...string) *stageBehaviour {
	m := map[string]*behaviour{}
	for _, stage := range stages {
		m[stage] = &behaviour{}
	}
	return &stageBehaviour{t, m}
}

func (s *stageBehaviour) Executed(stage string) bool {
	if b, ok := s.m[stage]; ok {
		return b.executed
	}
	s.t.Fatalf("error shouldErr stage: %s not mapped", stage)
	return false

}

func (s *stageBehaviour) Change(stageName string, result any, shouldError v12.StageError) *stageBehaviour {
	s.m[stageName] = &behaviour{executed: false, err: shouldError, result: result}
	return s
}

func (s *stageBehaviour) ResetExecutions() {
	for _, v := range s.m {
		v.executed = false
	}
}

func (s *stageBehaviour) exec(stage string) {
	if b, ok := s.m[stage]; ok {
		b.executed = true
		return
	}
	s.t.Fatalf("error exec stage: %s not mapped", stage)
}

func (s *stageBehaviour) shouldErr(stage string) v12.StageError {
	if b, ok := s.m[stage]; ok {
		return b.err
	}
	s.t.Fatalf("error shouldErr stage: %s not mapped", stage)
	return nil
}

func (s *stageBehaviour) shouldReturn(stage string) any {
	if b, ok := s.m[stage]; ok {
		return b.result
	}
	s.t.Fatalf("error shouldErr stage: %s not mapped", stage)
	return nil
}
