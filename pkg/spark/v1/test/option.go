package module_test_runner

import (
	"encoding/json"

	sparkv1 "github.com/azarc-io/vth-faas-sdk-go/pkg/spark/v1"
	"gopkg.in/yaml.v3"
)

/************************************************************************/
// SPARK OPTIONS
/************************************************************************/

type testOpts struct {
	configBasePath string
}

type Option = func(je *testOpts) *testOpts

func WithBasePath(configBasePath string) Option {
	return func(jw *testOpts) *testOpts {
		jw.configBasePath = configBasePath
		return jw
	}
}

func WithSparkConfigYAML(d []byte) sparkv1.Option {
	var m map[string]any
	if err := yaml.Unmarshal(d, &m); err != nil {
		panic(err)
	}

	return sparkv1.WithSparkConfig(m)
}

func WithSparkConfigJSON(d []byte) sparkv1.Option {
	var m map[string]any
	if err := json.Unmarshal(d, &m); err != nil {
		panic(err)
	}

	return sparkv1.WithSparkConfig(m)
}
