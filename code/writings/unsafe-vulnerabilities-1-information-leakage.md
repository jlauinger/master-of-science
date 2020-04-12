# Exploitation Exercise with Go unsafe.Pointer

April 12, 2020. Johannes Lauinger

Go is a generally safe language. It has memory safety measures that should render common buffer overflow vulnerabilities
which often exist in C programs impossible.

The `unsafe` std package defeats this memory safety. With `unsafe.Pointer`, we can create a pointer of arbitrary type.
The compiler can't and won't enforce safety measures on this pointer.

In this first of a three-part series on practically exploiting unsafe.Pointer usage, we will cover the possibilities
that come with `unsafe.Pointer` and show a first problem: a possible information leakage.


## Parts:

 1. Information leakage (enjoy!)
 2. [Code flow redirection](unsafe-vulnerabilities-2-code-flow-redirection.md)
 3. [ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)


## Buffer overflows, part 1: the stack layout 

A CPU uses a stack to keep track of data that is meaningful in the current context. Most importantly, it is used for
calling functions. A stack is a data structure that grows like a tower of things. New items can be pushed onto the
stack, and items on the stack can be removed or popped.

The stack used in the `x86_64` architecture is an area in the RAM which is identified by the stack pointer register
`$rsp`. When the current program calls a function, the return address as well as some function parameters (more on this
later) are pushed on the stack, and the processor jumps to the first instruction of the function. This jump is done
by setting the instruction pointer register `$rip`. Then, when the function returns (by executing the `ret` instruction),
the return address is popped from the stack and put into the `$rip` register.

Pushing something onto the stack is done be decrementing the stack pointer by some amount, e.g. a processor word (8 byte
on 64-bit architecture). Then the data is written to the address where the stack pointer now points to. Decrementing the
stack pointer marks the memory region as belonging to the stack. When popping values from the stack, the stack pointer
is incremented again, marking the memory region as free again. Because the stack pointer decrements with new data, we
can say that the stack grows to the bottom, starting from high addresses in memory and growing to low addresses.

The function can store local variables on the stack (inside its so-called stack frame). These are pushed onto the stack
after the return address, meaning the variables are at lower memory addresses than the return address. Furthermore,
variables on the stack are located directly next to each other. This is why bounds checking is very important for
buffers. Reading or writing outside the bounds of a variable means we are reading or writing other variables. We call
this buffer overflow.


## Go memory safety

Go employs some safety techniques that prevent buffer overflows, amongst other vulnerabilities. The type system strictly
encodes the buffer length of variables, e.g. we have `[8]byte` and `[16]byte` as completely different types with no
casting from the short buffer to the long buffer. This prevents the misuse of memory regions which will eventually lead
to a potentially exploitable buffer overflow.

Dangerous operations common to C programs such as pointer casting and the infamous, no-bounds-checking `gets()` function
are therefore impossible with safe Go programs.

However, there exists the `unsafe` package and with it the `unsafe.Pointer` type. This pointer type is special in that
it can participate in type operations that would otherwise be forbidden:

 1. we can cast any pointer type into `unsafe.Pointer`
 2. we can cast `unsafe.Pointer` into any pointer type
 3. we can cast `unsafe.Pointer` into `uintptr`, which is essentially the address as an integer
 4. we can cast `uintptr` into `unsafe.Pointer`
 
Points 1 and 2 allow type-casting between arbitrary types, and points 3 and 4 allow pointer arithmetic. With these
powers however comes great responsibility: using them removes the safety net of Go, meaning we're back at the security
madness of plain C code. The `unsafe` package must therefore be used only with extreme caution.

In the following proof of concepts, we will demonstrate some of the potential vulnerabilities that can be introduced
surprisingly fast when using `unsafe.Pointer`.


## Information leakage POC

In this short proof of concept, let's assume there is a buffer of harmless, public data. It is called `harmlessData`
and it might store e.g. the version of the program, or the name of a logged-in user.

Behind it, there is a declaration of a secret data buffer. For the sake of the argument, imagine that it might be some
private information about a logged-in user. Similar to the famous Heartbleed bug, it might also be private key data of
e.g. a TLS certificate.

```go
func main() {
    // this could be some public information, e.g. version information
    harmlessData := [8]byte{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'}
    // however, this could be critical private information such as a TLS private key
    secret := [17]byte{'l', '3', '3', 't', '-', 'h', '4', 'x', 'x', '0', 'r', '-', 'w', '1', 'n', 's', '!'}
    
    // ...
}
```

Next, the buffer is cast. Using the `unsafe.Pointer` type, we can do any type casting we want, defeating the Go memory
safety measures. Here, we cast the buffer into a new byte buffer, but with a bigger size. After this, we print the new
(dangerous) buffer.

```go
    // ...
    // (accidentally) cast harmless buffer into a new buffer type of wrong size
    var dangerousData = (*[8+17]byte)(unsafe.Pointer(&harmlessData[0]))
    
    // print (misused) buffer
    fmt.Println(string((*dangerousData)[:]))
}
```

Running this script will read the newly created, dangerous buffer. The length information will be inaccurate, and thus
the program will read memory after the end of the harmless data, revealing the secret data:

```shell script
$ go run main.go 
AAAAAAAAl33t-h4xx0r-w1ns!
```

This is a clear information leak, actually quite similar to what happened in the Heartbleed accident.

The threat model is a simple miscast into the wrong buffer size, which might occur in a big project e.g. if there is
a miscommunication or API documentation failure which leads a developer to assume false information about the buffer
size.


Continue reading with [Part 2: Code Flow Redirection](unsafe-vulnerabilities-2-code-flow-redirection.md)