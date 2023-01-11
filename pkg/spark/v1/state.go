package sparkv1

type JobState struct {
	JobContext   *JobMetadata
	StageResults map[string]Bindable
}
