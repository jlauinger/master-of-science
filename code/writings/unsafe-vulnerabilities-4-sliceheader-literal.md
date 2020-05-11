# Exploitation Exercise with Go unsafe.Pointer: Unsafe Usages of Slice Headers (Part 4)

In this fourth part, we will explore a very common, but unsafe code pattern: creating `reflect.SliceHeader` and
`reflect.StringHeader` objects from scratch instead of deriving them by cast.


## Parts:

 1. [Information leakage](unsafe-vulnerabilities-1-information-leakage.md)
 2. [Code flow redirection](unsafe-vulnerabilities-2-code-flow-redirection.md)
 3. [ROP and spawning a shell](unsafe-vulnerabilities-3-rop-and-spawning-a-shell.md)
 4. SliceHeader literals (enjoy!)


## Executing code on the stack


## Introducing a static code analysis tool!


## Acknowledgements

This blog post was written as part of my work on my Master's thesis at the 
[Software Technology Group](https://www.stg.tu-darmstadt.de/stg/homepage.en.jsp) at TU Darmstadt.