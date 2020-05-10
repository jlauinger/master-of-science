# Exploitation Exercise with Go unsafe.Pointer: Unsafe Usages of Slice Headers (Part 4)

In this fourth part, we will explore a very common, but unsafe code pattern: creating `reflect.SliceHeader` and
`reflect.StringHeader` objects from scratch instead of deriving them by cast.


## Parts:

 1. [Information leakage](unsafe-vulnerabilities-1-information-leakage.md)
 2. [Code flow redirection](unsafe-vulnerabilities-2-code-flow-redirection.md)
 3. [ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)
 4. SliceHeader literals (enjoy!)


## Executing code on the stack
