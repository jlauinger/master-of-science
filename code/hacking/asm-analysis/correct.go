package main

import (
	"reflect"
	"unsafe"
)

//go:noinline
func CorrectCast(s string) (b []byte) {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHeader.Data = strHeader.Data
	sliceHeader.Cap = strHeader.Len
	sliceHeader.Len = strHeader.Len
	return
}

func main() {
	s := "hallo"
	_ = CorrectCast(s)
}