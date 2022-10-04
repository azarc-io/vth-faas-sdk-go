package context

import (
	"github.com/azarc-io/vth-faas-sdk-go/pkg/api"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/v1"
)

type Compensation struct {
	jobContext   Job
	stageContext Stage
}

func NewCompensationContext(jobCtx Job) Compensation {
	return Compensation{jobContext: jobCtx}
}

// TODO stage name inside compensation we must concatened a prefix like 'compensate' to the name of the stage
func (c Compensation) Stage(name string, sdf api.StageDefinitionFn) api.StageChain {
	//TODO implement me
	panic("implement me")
}

func (c Compensation) WithStageStatus(names []string, value any) bool {
	//TODO implement me
	panic("implement me")
}

func (c Compensation) GetVariable(s string) sdk_v1.Variable {
	//TODO implement me
	panic("implement me")
}

func (c Compensation) SetVariable(name string, value any, mimeType string) error {
	//TODO implement me
	panic("implement me")
}
