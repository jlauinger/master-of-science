# Label comparison

Match on file name by concatenating my file_name column with the package import path excluding module path, then for
the matches only take those where the module path contains the Costa project identifier. See the notebook code for
details.

Only 20 labeled instances in Costa's and my data set match. Differences:

 - Costa only looks at project code, not dependencies
 - Costa only labels files, not lines in files
 - Costa did not even label all files, there are lots of missing labels in their data set
 - Costa only did one label dimension, for some instances it matches my purpose label and for some other it matches my
   what label
   
Stats:
 - 15 instances where Costa's label is missing
 - 3 matches
 - 2 non-matches
 
In the non-matches, my labeling is better / more detailed.


## Detailed matches and mismatches

```
/mask.go
github.com/gorilla/websocket
43.0
262    gorilla/websocket/mask.go
262    gorilla/websocket

pointer-arithmetic / efficiency
262    Performance Optimization
262    Program
262    NaN

*(*uintptr)(unsafe.Pointer(uintptr(unsafe.Pointer(&b[0])) + uintptr(i))) ^= kw

MATCH: matches

--------

proto/pointer_unsafe.go
github.com/golang/protobuf
83.0
231    golang/protobuf/proto/pointer_unsafe.go
231    golang/protobuf

cast-pointer / efficiency
231    NaN
231    NaN
231    NaN

return pointer{p: (*[2]unsafe.Pointer)(unsafe.Pointer(i))[1]}

MATCH: not labeled by Costa

--------

codec/helper_unsafe.go
github.com/ugorji/go
199.0
549    ugorji/go/codec/helper_unsafe.go
549    ugorji/go

cast-struct / serialization
549    NaN
549    NaN
549    NaN

urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

MATCH: not labeled by Costa

--------

/alloc.go
github.com/yuin/gopher-lua
62.0
585    yuin/gopher-lua/alloc.go
585    yuin/gopher-lua

cast-header / efficiency
585    NaN
585    NaN
585    NaN

al.fheader = (*reflect.SliceHeader)(unsafe.Pointer(&al.fptrs))

MATCH: not labeled by Costa

--------

/murmur32.go
github.com/spaolacci/murmur3
116.0
508    spaolacci/murmur3/murmur32.go
508    spaolacci/murmur3

cast-basic / efficiency
508    NaN
508    NaN
508    NaN

k1 := *(*uint32)(unsafe.Pointer(p))

MATCH: not labeled by Costa

--------

cmp/internal/value/pointer_unsafe.go
github.com/google/go-cmp
25.0
245    google/go-cmp/cmp/internal/value/pointer_unsaf...
245    google/go-cmp

cast-pointer / reflect
245    Reflect
245    Program
245    NaN

return Pointer{unsafe.Pointer(v.Pointer()), v.Type()}

MATCH: matches

--------

proto/pointer_unsafe.go
github.com/golang/protobuf
294.0
231    golang/protobuf/proto/pointer_unsafe.go
231    golang/protobuf

delegate / atomic
231    NaN
231    NaN
231    NaN

atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(p)), unsafe.Pointer(v))

MATCH: not labeled by Costa

--------

codec/helper_unsafe.go
github.com/ugorji/go
194.0
549    ugorji/go/codec/helper_unsafe.go
549    ugorji/go

cast-struct / serialization
549    NaN
549    NaN
549    NaN

urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

MATCH: not labeled by Costa

--------

codec/helper_unsafe.go
github.com/ugorji/go
219.0
549    ugorji/go/codec/helper_unsafe.go
549    ugorji/go

cast-struct / serialization
549    NaN
549    NaN
549    NaN

urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

MATCH: not labeled by Costa

--------

/alloc.go
github.com/yuin/gopher-lua
65.0
585    yuin/gopher-lua/alloc.go
585    yuin/gopher-lua

pointer-arithmetic / serialization
585    NaN
585    NaN
585    NaN

e := *(*LValue)(unsafe.Pointer(al.nheader.Data + uintptr(al.top)*unsafe.Sizeof(_uv)))

MATCH: not labeled by Costa

--------

codec/helper_unsafe.go
github.com/ugorji/go
164.0
549    ugorji/go/codec/helper_unsafe.go
549    ugorji/go

cast-struct / serialization
549    NaN
549    NaN
549    NaN

urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

MATCH: not labeled by Costa

--------

codec/helper_unsafe.go
github.com/ugorji/go
189.0
549    ugorji/go/codec/helper_unsafe.go
549    ugorji/go

cast-basic / serialization
549    NaN
549    NaN
549    NaN

urv := (*unsafeReflectValue)(unsafe.Pointer(&rv))

MATCH: not labeled by Costa

--------

/inotify.go
github.com/fsnotify/fsnotify
244.0
180    fsnotify/fsnotify/inotify.go
180    fsnotify/fsnotify

cast-struct / serialization
180    System Call
180    Program
180    unix

raw := (*unix.InotifyEvent)(unsafe.Pointer(&buf[offset]))

MATCH: does not match, mine is more detailed

--------

codec/helper_unsafe.go
github.com/ugorji/go
50.0
549    ugorji/go/codec/helper_unsafe.go
549    ugorji/go

cast-header / efficiency
549    NaN
549    NaN
549    NaN

return *(*string)(unsafe.Pointer(&sx))

MATCH: not labeled by Costa

--------

/decode_object.go
github.com/francoispqt/gojay
275.0
179    francoispqt/gojay/decode_object.go
179    francoispqt/gojay

cast-basic / efficiency
179    NaN
179    NaN
179    NaN

return *(*string)(unsafe.Pointer(&d)), false, nil

MATCH: not labeled by Costa

--------

/segment.go
github.com/coocood/freecache
274.0
115    coocood/freecache/segment.go
115    coocood/freecache

cast-bytes / efficiency
115    Convert between Types
115    Program
115    Convert from byte

entryHdr := (*entryHdr)(unsafe.Pointer(&entryHdrBuf[0]))

MATCH: matches

--------

/mask.go
github.com/gorilla/websocket
24.0
262    gorilla/websocket/mask.go
262    gorilla/websocket

pointer-arithmetic / serialization
262    Performance Optimization
262    Program
262    NaN

if n := int(uintptr(unsafe.Pointer(&b[0]))) % wordSize; n != 0 {

MATCH: mine more detailed

--------

/alloc.go
github.com/yuin/gopher-lua
71.0
585    yuin/gopher-lua/alloc.go
585    yuin/gopher-lua

cast-pointer / efficiency
585    NaN
585    NaN
585    NaN

ep.word = unsafe.Pointer(fptr)

MATCH: not labeled by Costa

--------

proto/pointer_unsafe.go
github.com/golang/protobuf
291.0
231    golang/protobuf/proto/pointer_unsafe.go
231    golang/protobuf

delegate / atomic
231    NaN
231    NaN
231    NaN

return (*unmarshalInfo)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(p))))

MATCH: not labeled by Costa

--------

jlexer/bytestostr.go
github.com/mailru/easyjson
21.0
367    mailru/easyjson/jlexer/bytestostr.go
367    mailru/easyjson

cast-header / efficiency
367    NaN
367    NaN
367    NaN

h := (*reflect.SliceHeader)(unsafe.Pointer(&data))

MATCH: not labeled by Costa
```