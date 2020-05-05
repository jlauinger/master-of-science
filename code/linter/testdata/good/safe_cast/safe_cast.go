package good

import (
	"reflect"
	"runtime"
	"unsafe"
)

func SaferCastString(str string) (b []byte) {
	strH := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sH := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sH.Len = strH.Len
	sH.Cap = strH.Len
	sH.Data = strH.Data
	runtime.KeepAlive(str)
	return
}
