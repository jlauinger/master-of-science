// Type your code here, or load an example.
// Your function name should start with a capital letter.
package testconv

import (
    "unsafe"
    "reflect"
    "runtime"
)

func Str2Bytes1(s string) (b []byte) {
    sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}

func Str2Bytes2(s string) (b []byte) {
    sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = len(s)
	bh.Len = len(s)
	return
}

func Str2Bytes3(s string) (b []byte) {
    bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = uintptr((unsafe.Pointer)((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))
	bh.Cap = len(s)
	bh.Len = len(s)
	return
}

func Str2Bytes4(s string) (b []byte) {
    sh := (*reflect.StringHeader)(unsafe.Pointer(&s))
	bh := (*reflect.SliceHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
        runtime.KeepAlive(&s)
	return
}

type stringHeader struct {
    Data unsafe.Pointer
    Len uint
}

type byteHeader struct {
    Data unsafe.Pointer
    Len, Cap uint
}

func Str2Bytes5(s string) (b []byte) {
    sh := (*stringHeader)(unsafe.Pointer(&s))
	bh := (*byteHeader)(unsafe.Pointer(&b))
	bh.Data = sh.Data
	bh.Cap = sh.Len
	bh.Len = sh.Len
	return
}