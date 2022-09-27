package errors

import "errors"

var (
	StageDoesNotExist = errors.New("stage does not exists")
	BindValueFailed   = errors.New("bind value failed")
)
