package internal

import "github.com/azarc-io/vth-faas-sdk-go/pkg/api"

type Job struct {
	// options Options
}

func (j Job) WithService(name string, fn api.StageDefinitionFn) {
	// job.options =
}
