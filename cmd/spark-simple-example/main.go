package main

import (
	"context"
	impl "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-simple-example/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	spark := impl.NewSpark()
	worker, err := sparkv1.NewSparkWorker(ctx, spark)
	if err != nil {
		panic(err)
	}
	worker.Run()
}
