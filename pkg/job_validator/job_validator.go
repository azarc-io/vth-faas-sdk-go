package job_validator

import (
	"context"
	"fmt"
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
	"sort"
	"strings"
)

type validatorContext struct {
	stageNameCounter map[string]uint
}

func Check(job api.Job) error {
	ctx := validatorContext{map[string]uint{}}
	job.Execute(ctx)
	return ctx.validate()
}

func (v validatorContext) validate() error {
	if duplicatedNames := v.duplicateStageNames(); len(duplicatedNames) > 0 {
		return fmt.Errorf("invalid stage chain. stage names must be unique. the following stage names have more than one occurrence: %s", strings.Join(duplicatedNames, ", "))
	}
	return nil
}

func (v validatorContext) duplicateStageNames() []string {
	var duplicatedNames []string
	for stageName, times := range v.stageNameCounter {
		if times > 1 {
			duplicatedNames = append(duplicatedNames, stageName)
		}
	}
	sort.Strings(duplicatedNames)
	return duplicatedNames
}

func (v validatorContext) Stage(name string, sdf api.StageDefinitionFn, options ...api.StageOption) api.StageChain {
	v.stageNameCounter[name] += 1
	return v
}

func (v validatorContext) Canceled(fn api.CancelDefinitionFn) api.CanceledChain {
	fn(v)
	return v
}

func (v validatorContext) Compensate(fn api.CompensateDefinitionFn) api.CompensateChain {
	fn(v)
	return v
}

func (v validatorContext) Complete(fn api.CompletionDefinitionFn) api.CompleteChain {
	fn(v)
	return v
}

func (v validatorContext) WithStageStatus(names []string, value any) bool {
	return true
}

func (v validatorContext) Ctx() context.Context {
	return nil
}

func (v validatorContext) JobKey() string {
	return ""
}

func (v validatorContext) CorrelationID() string {
	return ""
}

func (v validatorContext) TransactionID() string {
	return ""
}

func (v validatorContext) Payload() any {
	return nil
}

func (v validatorContext) GetStage(jobKey, name string) (*sdk_v1.StageStatus, error) {
	return nil, nil
}

func (v validatorContext) GetStageResult(jobKey, stageName string) (*sdk_v1.StageResult, error) {
	return nil, nil
}

func (v validatorContext) Err() api.StageError {
	return nil
}

func (v validatorContext) GetVariable(name, stage string) (*sdk_v1.Variable, error) {
	return nil, nil
}

func (v validatorContext) SetVariable(variable *sdk_v1.SetVariableRequest) error {
	return nil
}
