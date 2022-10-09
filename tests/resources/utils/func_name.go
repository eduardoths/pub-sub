package utils

import (
	"reflect"
	"runtime"
)

func FuncName(fn any) string {
	return runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
}
