# Exploitation Exercise with Goland unsafe.Pointer: Code Flow Redirection (Part 2)

Go is a generally safe language. It has memory safety measures that should render common buffer overflow vulnerabilities
which often exist in C programs impossible.

The `unsafe` std package defeats this memory safety. With `unsafe.Pointer`, we can create a pointer of arbitrary type.
The compiler can't and won't enforce safety measures on this pointer.

In this first of a three-part series on practically exploiting unsafe.Pointer usage, we will cover the possibilities
that come with `unsafe.Pointer` and show a first problem: a possible information leakage.

## Parts:

 1. [Information leakage](unsafe-vulnerabilities-1-information-leakage.md)
 2. Code flow redirection (enjoy!)
 3. [ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)


## Buffer overflow, part 2: controlling the return address


## Crafting a binary exploit


## Code flow redirection POC

First, we compile a `win` function into the binary. This is so that we have a target to redirect the code flow to. This
is a good first step in code flow exploitation. The function does not do very much, it simply prints "win!" so that we
know we did good:

```go
func win() {
    fmt.Println("win!")
}
```

Next, we create an address function that uses `fmt.Sprintf` to convert a reference to an actual integer. This is used to
automatically get the address of the `win` function:

```go
func address(i interface{}) int {
    addr, _ := strconv.ParseUint(fmt.Sprintf("%p", i), 0, 0)
    return int(addr)
}
```

We assume there is a utility function to copy arrays 

Running this program creates the following output:

```shell script
$ go run main.go 
win!
unexpected fault address 0x0
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x80 addr=0x0 pc=0x4940ee]

goroutine 1 [running]:
runtime.throw(0x4c2d75, 0x5)
	/usr/lib/go/src/runtime/panic.go:1112 +0x72 fp=0xc0000a2f58 sp=0xc0000a2f28 pc=0x42ebd2
runtime.sigpanic()
	/usr/lib/go/src/runtime/signal_unix.go:694 +0x3cc fp=0xc0000a2f88 sp=0xc0000a2f58 pc=0x4429dc
runtime: unexpected return pc for main.win called from 0x4848484847474747
stack: frame={sp:0xc0000a2f88, fp:0xc0000a2f90} stack=[0xc0000a2000,0xc0000a3000)
000000c0000a2e88:  0000000000000001  0000000000000000 
000000c0000a2e98:  000000c0000a2ed8  0000000000430404 <runtime.gwrite+164> 
000000c0000a2ea8:  0000000000000002  00000000004c2ad6 
000000c0000a2eb8:  0000000000000001  0000000000000001 
000000c0000a2ec8:  000000c0000a2f45  0000000000000003 
000000c0000a2ed8:  000000c0000a2f28  0000000000430c28 <runtime.printstring+120> 
000000c0000a2ee8:  000000000042ed97 <runtime.fatalthrow+87>  000000c0000a2ef8 
000000c0000a2ef8:  0000000000458580 <runtime.fatalthrow.func1+0>  000000c000000180 
000000c0000a2f08:  000000000042ebd2 <runtime.throw+114>  000000c0000a2f28 
000000c0000a2f18:  000000c0000a2f48  000000000042ebd2 <runtime.throw+114> 
000000c0000a2f28:  000000c0000a2f30  0000000000458500 <runtime.throw.func1+0> 
000000c0000a2f38:  00000000004c2d75  0000000000000005 
000000c0000a2f48:  000000c0000a2f78  00000000004429dc <runtime.sigpanic+972> 
000000c0000a2f58:  00000000004c2d75  0000000000000005 
000000c0000a2f68:  0000000000000000  0000000000000000 
000000c0000a2f78:  4444444443434343  00000000004940ee <main.win+126> 
000000c0000a2f88: <4848484847474747 >4a4a4a4a49494949 
000000c0000a2f98:  4c4c4c4c4b4b4b4b  4e4e4e4e4d4d4d4d 
000000c0000a2fa8:  505050504f4f4f4f  0000000000000000 
000000c0000a2fb8:  000000c000000180  000000c0000a2fae 
000000c0000a2fc8:  00000000004ca800  0000000000000000 
000000c0000a2fd8:  000000000045b911 <runtime.goexit+1>  0000000000000000 
000000c0000a2fe8:  0000000000000000  0000000000000000 
000000c0000a2ff8:  0000000000000000 
main.win()
	/home/johannes/studium/s14/masterarbeit/code/exploits/code-flow-redirection/main.go:25 +0x7e fp=0xc0000a2f90 sp=0xc0000a2f88 pc=0x4940ee
exit status 2
```

Quite obvious from the big stack trace, we see that the program crashed. But more importantly, we see the `win!` output
right at the top, which means that the `win` function was indeed executed. We don't actually care about the program
crash, the objective was to decide which code should be executed and this was successful!


Continue reading with [Part 3: ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)