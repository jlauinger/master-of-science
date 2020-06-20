package main

import (
	"reflect"
	"runtime"
	"time"
	"unsafe"
)

func main() {
	s := "hello"
	_ = CorrectCast(s)
	_ = CorrectCastWithKeepAlive(s)
	_ = OneStatement(s)
	_ = WithVariable(s)
	_ = WithVariableAndDelay(s)
}

//go:noinline
func CorrectCast(s string) (b []byte) {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHeader.Data = strHeader.Data
	sliceHeader.Cap = strHeader.Len
	sliceHeader.Len = strHeader.Len
	return
}

//go:noinline
func CorrectCastWithKeepAlive(s string) (b []byte) {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	sliceHeader.Data = strHeader.Data
	sliceHeader.Cap = strHeader.Len
	sliceHeader.Len = strHeader.Len
	runtime.KeepAlive(s)
	return
}

//go:noinline
func OneStatement(s string) []byte {
	strHeader := (*reflect.StringHeader)(unsafe.Pointer(&s))
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: strHeader.Data,
		Cap:  strHeader.Len,
		Len:  strHeader.Len,
	}))
}

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

//go:noinline
func WithVariableAndDelay(s string) []byte {
	sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	sliceHeader := &reflect.SliceHeader{
		Data: sh.Data,
		Len:  sh.Len,
		Cap:  sh.Len,
	}
	time.Sleep(1 * time.Nanosecond)
	return *(*[]byte)(unsafe.Pointer(sliceHeader))
}