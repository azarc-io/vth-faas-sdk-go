package test

type noopLogger struct {
}

func (n noopLogger) Error(_ error, _ string, _ ...interface{}) {
}

func (n noopLogger) Fatal(_ error, _ string, _ ...interface{}) {
}

func (n noopLogger) Info(_ string, _ ...interface{}) {
}

func (n noopLogger) Warn(_ string, _ ...interface{}) {
}

func (n noopLogger) Debug(_ string, _ ...interface{}) {
}
