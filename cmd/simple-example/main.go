package main

import (
	"context"
	impl "github.com/azarc-io/vth-faas-sdk-go/cmd/simple-example/internal"
	sdk_v1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	spark := impl.NewSpark(cancel)
	worker, err := sdk_v1.NewSparkWorker(ctx, spark)
	if err != nil {
		panic(err)
	}
	worker.Run()
}
