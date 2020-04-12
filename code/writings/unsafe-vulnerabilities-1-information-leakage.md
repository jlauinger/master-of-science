# Exploitation Exercise with Goland unsafe.Pointer

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

The stack used in the x86_64 architecture is an area in the RAM which is identified by the stack pointer register
`$rsp`.


## Go memory safety


## unsafe.Pointer


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
 