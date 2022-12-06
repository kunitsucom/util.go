package runtimez

import (
	"path"
	"runtime"
)

func FuncName(skip int) (funcName string) {
	pc, _, _, _ := runtime.Caller(skip + 1) //nolint:dogsled
	return path.Base(runtime.FuncForPC(pc).Name())
}

func FullFuncName(skip int) (funcName string) {
	pc, _, _, _ := runtime.Caller(skip + 1) //nolint:dogsled
	return runtime.FuncForPC(pc).Name()
}
