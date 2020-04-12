# Exploitation Exercise with Goland unsafe.Pointer: ROP and Spawning a Shell (Part 3)

In this third part, we will see how to take the idea of code flow redirection one step further. We will do an arbitrary
remote code execution using the Return Oriented Programming (ROP) technique. In the end, we will reach the classic goal
for binary exploitation and run a shell in the program context.


## Parts:

 1. [Information leakage](unsafe-vulnerabilities-1-information-leakage.md)
 2. [Code flow redirection](unsafe-vulnerabilities-2-code-flow-redirection.md)
 3. ROP and spawning a shell (enjoy!)


## DEP, ASLR, canaries and more: mitigations against buffer overflows


## Executing code on the stack


## Return2libc


## Return oriented programming


## POC: Spawning a shell

We analyze the following short Go program which contains a critical flaw:

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

Putting the ROP techniques from above into play, the plan looks like this:

 1. Set the executable and writable flags for a memory page belonging to the program
 2. Write some code into the page that spawns a shell
 3. Jump to that code
 
**Step 1: Get a memory page with RWX permissions**

**Step 2: Write shell code into the page**

**Step 3: Jump to the code**


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
