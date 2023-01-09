package common

//********************************************************************************************
// MODULE API
//********************************************************************************************

type (
	GetStageResultRequest struct {
		WorkflowId string
		RunId      string
		StageName  string
	}

	GetStageResultResponse struct {
		Data []byte
	}
)

type ModuleApi interface {
	GetStageResult(req *GetStageResultRequest) (*GetStageResultResponse, error)
}
