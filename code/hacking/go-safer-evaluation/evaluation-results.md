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

I need to check the following files with go-safer, go vet, and gosec:

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

 1. type.go:475 in package runtime of module std: **NO**  
    go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 2. value.go:1815 in package reflect of module std: **NO**  
    go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 3. alg.go:317 in package runtime of module std: **NO**  
    go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 4. value.go:1800 in package reflect of module std: **NO**  
    go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)

**app**:

  1. label.go:118 in package golang.org/x/tools/internal/event/label of module golang.org/x/tools: **NO**  
     go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: 117 **POSITIVE** (false)
  2. iterator_native2.go:282 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 281 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
  3. iterator_native.go:1078 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 1073 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
  4. iterator_native2.go:324 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 319 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
  5. array.go:257 in package gorgonia.org/tensor of module gorgonia.org/tensor: **YES**  
     go-safer: 251 **POSTIIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
  6. iterator_native2.go:434 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 433 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
  7. alloc.go:62 in package github.com/yuin/gopher-lua of module github.com/yuin/gopher-lua: **NO**  
     go-safer: errors **NEGATIVE** (true), go vet: errors **NEGATIVE** (true), gosec: **POSITIVE** (false)
  8. iterator_native.go:654 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 653 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
  9. value_unsafe.go:92 in package google.golang.org/protobuf/reflect/protoreflect of module google.golang.org/protobuf: **NO**  
     go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 10. iterator_native.go:763 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 762 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 11. writer_unsafe.go:13 in package github.com/philhofer/fwd of module github.com/philhofer/fwd: **YES**  
     go-safer: **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 12. iterator_native2.go:134 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 129 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 13. iterator_native2.go:206 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 205 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 14. iterator_native.go:514 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 513 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 15. iterator_native.go:304 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 303 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 16. iterator_native.go:623 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 622 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 17. iterator_native.go:794 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 793 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 18. pcap.go:355 in package github.com/google/gopacket/pcap of module github.com/google/gopacket: **NO**  
     go-safer: errors **NEGATIVE** (true), go vet: errors **NEGATIVE** (true), gosec: **POSITIVE** (false)
 19. label.go:117 in package golang.org/x/tools/internal/event/label of module golang.org/x/tools: **NO**  
     go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 20. strings_unsafe.go:46 in package google.golang.org/protobuf/internal/strs of module google.golang.org/protobuf: **NO**  
     go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 21. marshalers.go:66 in package github.com/cilium/ebpf of module github.com/cilium/ebpf: **YES**  
     go-safer: 67 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: 73 **POSITIVE** (true)
 22. xxhash_unsafe.go:30 in package github.com/cespare/xxhash/v2 of module github.com/cespare/xxhash/v2: **NO**  
     go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 23. node.go:234 in package go.etcd.io/bbolt of module go.etcd.io/bbolt: **YES**  
     go-safer: **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 24. unsafe_slice.go:140 in package github.com/modern-go/reflect2 of module github.com/modern-go/reflect2: **NO**  
     go-safer: errors **NEGATIVE** (true), go vet: errors **NEGATIVE** (true), gosec: errors **NEGATIVE** (true)
 25. iterator_native.go:798 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 793 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 26. map.go:128 in package github.com/weaveworks/ps of module github.com/weaveworks/ps: **YES**  
     go-safer: 129 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 27. iterator_native.go:277 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 272 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 28. generic.go:61 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **NO**  
     go-safer: 62 **POSITIVE** (false), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 29. helper_unsafe.go:50 in package github.com/ugorji/go/codec of module github.com/ugorji/go: **YES**  
     go-safer: 49 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 30. helper_unsafe.go:52 in package github.com/hashicorp/go-msgpack/codec of module github.com/hashicorp/go-msgpack: **YES**  
     go-safer: errors **NEGATIVE** (false), go vet: errors **NEGATIVE** (false), gosec: errors **NEGATIVE** (false)
 31. iterator_native.go:907 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 902 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 32. iterator_native2.go:168 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 167 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 33. page.go:135 in package go.etcd.io/bbolt of module go.etcd.io/bbolt: **YES**  
     go-safer: **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 34. unsafe.go:31 in package github.com/elastic/go-structform/internal/unsafe of module github.com/elastic/go-structform: **YES**  
     go-safer: 32 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 35. iterator_native.go:1148 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 1143 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: 1143,1147 **POSITIVE** (true)
 36. iterator_native2.go:476 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 471 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 37. value_unsafe.go:75 in package google.golang.org/protobuf/reflect/protoreflect of module google.golang.org/protobuf: **NO**  
     go-safer: **NEGATIVE** (true), go vet: **NEGATIVE** (true), gosec: **POSITIVE** (false)
 38. iterator_native2.go:438 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 433 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)
 39. iterator_native.go:588 in package gorgonia.org/tensor/native of module gorgonia.org/tensor: **YES**  
     go-safer: 583 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: 586 **POSITIVE** (true)
 40. bytestostr.go:21 in package github.com/mailru/easyjson/jlexer of module github.com/mailru/easyjson: **YES**  
     go-safer: 22 **POSITIVE** (true), go vet: **NEGATIVE** (false), gosec: **POSITIVE** (true)

## Step 1 evaluation summary

In summary, we have:

| **Tool** | **TP** | **FP** | **TN** | **FN** | **Recall** | **Precision** | **Accuracy** |
|----------|--------|--------|--------|--------|------------|---------------|--------------|
| go-safer |   29   |    1   |   13   |    1   |   0.967    |     0.967     |    0.955     |
| go vet   |    0   |    0   |   14   |   30   |   0        |     -         |    0.318     |
| gosec   |   29   |   13   |    1   |    1   |   0.967    |     0.690     |    0.681     |

Recall = TP/(TP+FN)  
Precision = TP/(TP+FP)  
Accuracy = (TP+TN)/(TP+TN+FP+FN)

We see that Go Vet does not find anything at all, while Gosec just marks everything as potentially unsafe except for
when there is a compilation error. This shows that the existing tools are useless for this specific pattern.
Go-safer on the other hand achieved excellent scores with 96.7% recall and accuracy due to only 1 false negative / 
positive, yielding 95.5% accuracy for this task!


## Step 2: evaluation with manually analyzed projects
