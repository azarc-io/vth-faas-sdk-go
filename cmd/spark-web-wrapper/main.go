package main

import (
	"context"
	"flag"
	"os"

	impl "github.com/azarc-io/vth-faas-sdk-go/cmd/spark-web-wrapper/internal"
	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// variables declaration
	var baseUrl string
	flag.StringVar(&baseUrl, "s", "", "Server Base URL: (required)")
	flag.Parse() // after declaring flags we need to call it
	if baseUrl == "" {
		flag.Usage()
		os.Exit(1)
	}

	spark := impl.NewSpark(baseUrl)
	worker, err := sparkv1.NewSparkWorker(ctx, spark)
	if err != nil {
		panic(err)
	}
	worker.Run()
}
