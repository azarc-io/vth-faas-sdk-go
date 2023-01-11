package module_test_runner

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
