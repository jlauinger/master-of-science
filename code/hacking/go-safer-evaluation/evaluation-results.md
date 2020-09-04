# go-safer evaluation

## Step 1: evaluation with labeled data set

Instances with label `cast-header` in the labeled data set: 40 in app, 4 in std.
 
Are those a security problem? Let's look at all of them.

**std**:

 1. type.go:475 in package runtime of module std: **NO**
    no uintptr but unsafe.Pointer, no back-conversion from header to slice
 2. value.go:1815 in package reflect of module std: **NO**
    conversion from real slice
 3. alg.go:317 in package runtime of module std: **NO**
    conversion from real slice
 4. value.go:1800 in package reflect of module std: **NO**
    conversion from real slice, internal reflect code
   
**app**:

  1. label.go:118 in package golang.org/x/tools/internal/event/label of module golang.org/x/tools: **NO**
     conversion from real string
  2. iterator_native2.go:282 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
  3. iterator_native.go:1078 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
  4. iterator_native2.go:324 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
  5. array.go:257 in package gorgonia.org/tensor of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
  6. iterator_native2.go:434 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
  7. alloc.go:62 in package github.com/yuin/gopher-lua of module github.com/yuin/gopher-lua: **NO**
     direct array access, header derived by cast from real slice
  8. iterator_native.go:654 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
  9. value_unsafe.go:92 in package google.golang.org/protobuf/reflect/protoreflect of module google.golang.org/protobuf: **NO**
     conversion from real slice
 10. iterator_native.go:763 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 11. writer_unsafe.go:13 in package github.com/philhofer/fwd of module github.com/philhofer/fwd: **YES**
     slice header literal is converted to real slice later on, possible GC race
 12. iterator_native2.go:134 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 13. iterator_native2.go:206 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 14. iterator_native.go:514 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 15. iterator_native.go:304 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 16. iterator_native.go:623 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 17. iterator_native.go:794 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 18. pcap.go:355 in package github.com/google/gopacket/pcap of module github.com/google/gopacket: **NO**
     conversion from real slice
 19. label.go:117 in package golang.org/x/tools/internal/event/label of module golang.org/x/tools: **NO**
     conversion from real slice
 20. strings_unsafe.go:46 in package google.golang.org/protobuf/internal/strs of module google.golang.org/protobuf: **NO**
     conversion from real slice
 21. marshalers.go:66 in package github.com/cilium/ebpf of module github.com/cilium/ebpf: **YES**
     slice header literal is converted to real slice later on, possible GC race
 22. xxhash_unsafe.go:30 in package github.com/cespare/xxhash/v2 of module github.com/cespare/xxhash/v2: **NO**
     conversion from real slice
 23. node.go:234 in package go.etcd.io/bbolt of module go.etcd.io/bbolt: **YES**
     slice header literal is converted to real slice later on, possible GC race
 24. unsafe_slice.go:140 in package github.com/modern-go/reflect2 of module github.com/modern-go/reflect2: **NO**
     conversion from unsafe pointer, but only to extract Cap, no back-conversion
 25. iterator_native.go:798 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 26. map.go:128 in package github.com/weaveworks/ps of module github.com/weaveworks/ps: **YES**
     slice header literal is converted to real slice later on, possible GC race
 27. iterator_native.go:277 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 28. generic.go:61 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **NO**
     conversion from real slice
 29. helper_unsafe.go:50 in package github.com/ugorji/go/codec of module github.com/ugorji/go: **YES**
     slice header literal is converted to real slice later on, possible GC race
 30. helper_unsafe.go:52 in package github.com/hashicorp/go-msgpack/codec of module github.com/hashicorp/go-msgpack: **YES**
     slice header literal is converted to real slice later on, possible GC race
 31. iterator_native.go:907 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 32. iterator_native2.go:168 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 33. page.go:135 in package go.etcd.io/bbolt of module go.etcd.io/bbolt: **YES**
     slice header literal is converted to real slice later on, possible GC race
 34. unsafe.go:31 in package github.com/elastic/go-structform/internal/unsafe of module github.com/elastic/go-structform: **YES**
     slice header literal is converted to real slice later on, possible GC race
 35. iterator_native.go:1148 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 36. iterator_native2.go:476 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 37. value_unsafe.go:75 in package google.golang.org/protobuf/reflect/protoreflect of module google.golang.org/protobuf: **NO**
     access to internals provided, but no back-conversion from uintptr slice representation
 38. iterator_native2.go:438 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 39. iterator_native.go:588 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**
     slice header literal is converted to real slice later on, possible GC race
 40. bytestostr.go:21 in package github.com/mailru/easyjson/jlexer of module github.com/mailru/easyjson: **YES**
     slice header literal is converted to real slice later on, possible GC race

Therefore in summary, I have 44 labeled verification snippets of which 14 are negative and 30 are positive.

I need to check the following files with go-safer, go vet, and golint:

 - github.com/cespare/xxhash/v2@v2.1.0/xxhash_unsafe.go
 - github.com/cilium/ebpf@v0.0.0-20191113100448-d9fb101ca1fb/marshalers.go
 - github.com/elastic/go-structform@v0.0.6/internal/unsafe/unsafe.go
 - github.com/google/gopacket@v1.1.17/pcap/pcap.go
 - github.com/hashicorp/go-msgpack@v1.1.5/codec/helper_unsafe.go
 - github.com/mailru/easyjson@v0.7.0/jlexer/bytestostr.go
 - github.com/modern-go/reflect2@v1.0.1/unsafe_slice.go
 - github.com/philhofer/fwd@v1.0.0/writer_unsafe.go
 - github.com/ugorji/go@v0.0.0-20170918222552-54210f4e076c/codec/helper_unsafe.go
 - github.com/weaveworks/ps@v0.0.0-20160725183535-70d17b2d6f76/map.go
 - github.com/yuin/gopher-lua@v0.0.0-20170403160031-b402f3114ec7/alloc.go
 - go.etcd.io/bbolt@v1.3.4/node.go
 - go.etcd.io/bbolt@v1.3.4/page.go
 - golang.org/x/tools@v0.0.0-20200428021058-7ae4988eb4d9/internal/event/label/label.go
 - google.golang.org/protobuf@v1.23.0/internal/strs/strings_unsafe.go
 - google.golang.org/protobuf@v1.23.0/reflect/protoreflect/value_unsafe.go
 - gorgonia.org/tensor@v0.9.6/array.go
 - gorgonia.org/tensor@v0.9.6/native/generic.go
 - gorgonia.org/tensor@v0.9.6/native/iterator_native.go
 - gorgonia.org/tensor@v0.9.6/native/iterator_native2.go
 - reflect/value.go
 - runtime/alg.go
 - runtime/type.go
 
The following are the finding results of the tools:

**std**:

 1. Label: **NO**
    go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 2. Label: **NO**
    go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 3. Label: **NO**
    go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 4. Label: **NO**
    go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**

**app**:

  1. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  2. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  3. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  4. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  5. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  6. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  7. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  8. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
  9. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 10. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 11. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 12. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 13. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 14. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 15. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 16. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 17. Label: **YES**
      go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 18. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 19. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 20. Label: **NO**
      go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 21. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 22. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 23. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 24. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 25. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 26. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 27. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 28. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 29. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 30. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 31. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 32. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 33. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 34. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 35. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 36. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 37. Label: **NO**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 38. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 39. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**
 40. Label: **YES**
     go-safer: **UNKNOWN**, go vet: **UNKNOWN**, golint: **UNKNOWN**

In summary, we have:

| **Tool** | **TP** | **FP** | **TN** | **FN** | **Recall** | **Precision** | **Accuracy** |
|----------|--------|--------|--------|--------|------------|---------------|--------------|
| go-safer |        |        |        |        |            |               |              |
| go vet   |        |        |        |        |            |               |              |
| golint   |        |        |        |        |            |               |              |


## Step 2: evaluation with manually analyzed projects
