package fixtures

import (
	"reflect"
	"runtime"
	"unsafe"
)

func UnsafeCastString(str string) []byte {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := &reflect.SliceHeader{
		Data: stringHeader.Data,
		Cap: stringHeader.Len,
		Len: stringHeader.Len,
	}
	return *(*[]byte)(unsafe.Pointer(sliceHeader))
}

func SaferCastString(str string) (b []byte) {
	stringHeader := (*reflect.StringHeader)(unsafe.Pointer(&str))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHeader.Len = stringHeader.Len
	sliceHeader.Cap = stringHeader.Len
	sliceHeader.Data = stringHeader.Data
	runtime.KeepAlive(str)
	return
}