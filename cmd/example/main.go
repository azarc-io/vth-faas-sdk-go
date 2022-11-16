package example

import (
	impl "github.com/azarc-io/vth-faas-sdk-go/cmd/example/internal"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/api/spark/v1"
)

func main() {
	spark := impl.NewSpark()
	worker, err := sdk_v1.NewSparkWorker(nil, spark)
	if err != nil {
		panic(err)
	}
	worker.Run()
}
