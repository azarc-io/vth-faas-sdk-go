package main

import (
	module_runner "github.com/azarc-io/vth-faas-sdk-go/internal/module-runner"
	"github.com/azarc-io/vth-faas-sdk-go/internal/signals"
)

func main() {
	cfg, err := module_runner.LoadModuleConfig(
		module_runner.WithBinBasePath("bin"),
		module_runner.WithBasePath("cmd/module-runner"))

	if err != nil {
		panic(err)
	}

	runner, err := module_runner.RunModule(cfg)
	if err != nil {
		panic(err)
	}

	<-signals.SetupSignalHandler()

	if err := runner.Stop(); err != nil {
		panic(err)
	}
}
