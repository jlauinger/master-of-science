package bad

import (
	"reflect"
	"unsafe"
)

func AlsoUnsafeCastString(str string) []byte {
	strH := (*reflect.StringHeader)(unsafe.Pointer(&str))
	var sH *reflect.SliceHeader
	sH.Len = strH.Len
	sH.Cap = strH.Len
	sH.Data = strH.Data
	return *(*[]byte)(unsafe.Pointer(sH))
}

