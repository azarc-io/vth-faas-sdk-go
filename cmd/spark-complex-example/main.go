package main

import (
	"context"

	impl "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-complex-example/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

func main() {
	spark := impl.NewSpark()
	worker, err := sparkv1.NewSparkWorker(context.Background(), spark)
	if err != nil {
		panic(err)
	}
	worker.Run()
}
