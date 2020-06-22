package main

import (
	"reflect"
	"unsafe"
)

//go:noinline
func WithVariable(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader := &reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	return *(*[]byte)(unsafe.Pointer(sliceHeader))
}

func main() {
	s := "hallo"
	_ = WithVariable(s)
}