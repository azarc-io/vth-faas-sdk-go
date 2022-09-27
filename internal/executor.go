package internal

import "github.com/azarc-io/vth-faas-sdk-go/pkg/api"

type executor struct {
}

func NewExecutor() api.StageChain {
	return &executor{}
}

func (e executor) Stage(name string, sdf api.StageDefinitionFn) api.StageChain {
	//TODO implement me
	panic("implement me")

	// stageClient.GetStage
	// if ok and stage.status == pending
	//   stageClient.SetState("stageName", "starting")
	//   execute the StageDefinitionFn()
	//   if StageDefinitionFn() returns and error stageClient.SetState("stageName", "failed")
	//   else send StageDefinitionFn() return to verathread stageClient.SetState("stageName", "complete", stageResult)
}

func (e executor) Complete(fn api.CompletionDefinitionFn) api.CompleteChain {
	//TODO implement me
	panic("implement me")
}

func (e executor) Compensate(fn api.CompensateDefinitionFn) api.CompensateChain {
	//TODO implement me
	panic("implement me")
}

func (e executor) Canceled(fn api.CancelDefinitionFn) api.CanceledChain {
	//TODO implement me
	panic("implement me")
}

func (e executor) Run() {
	//TODO implement me
	panic("implement me")
}
