package module_runner

/************************************************************************/
// MODULE OPTIONS
/************************************************************************/

type moduleOpts struct {
	configBasePath string
	binBasePath    string
}

type ModuleOption = func(je *moduleOpts) *moduleOpts

func WithBasePath(configBasePath string) ModuleOption {
	return func(jw *moduleOpts) *moduleOpts {
		jw.configBasePath = configBasePath
		return jw
	}
}

func WithBinBasePath(binBasePath string) ModuleOption {
	return func(je *moduleOpts) *moduleOpts {
		je.binBasePath = binBasePath
		return je
	}
}
