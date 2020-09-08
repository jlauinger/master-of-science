package main

import (
	"reflect"
	"unsafe"
)

//go:noinline
func WithVariable1Stmt(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}))
}

func main() {
	s := "hallo"
	_ = WithVariable1Stmt(s)
}