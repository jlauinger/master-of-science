package bad

import (
	"reflect"
	"unsafe"
)

type Protocol struct {
	Foo int
	Sh *reflect.SliceHeader
}

func UnsafeStringIntoProtocol(str string) []byte {
	strH := (*reflect.StringHeader)(unsafe.Pointer(&str))
	protocol := Protocol{}
	protocol.Sh.Len = strH.Len
	protocol.Sh.Cap = strH.Len
	protocol.Sh.Data = strH.Data
	return *(*[]byte)(unsafe.Pointer(protocol.Sh))
}
