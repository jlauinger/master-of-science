# Identification and Analysis of unsafe.Pointer Usage Patterns in Open-Source Go Code

Master Thesis by Johannes Lauinger

Submitted: October 26th, 2020

Advisor: Prof. Dr.-Ing. Mira Mezini  
Supervisors: Anna-Katharina Wickert, M.Sc., Dr. rer. nat. Lars Baumgärtner

Software Technology Group  
Department of Computer Science  
Technische Universität Darmstadt


## Citation

Cite this work as follows:

 - Lauinger, Johannes Tobias. "Identification and Analysis of unsafe.Pointer Usage Patterns in Open-Source Go Code." M.Sc. Thesis. Technische Universität Darmstadt, 2020.

BibTex:

```latex
@mastersthesis{lauinger2020,
    type    = {M.Sc. Thesis},
    author  = {Lauinger, Johannes Tobias},
    title   = {Identification and Analysis of unsafe.Pointer Usage Patterns in Open-Source Go Code},
    school  = {Technische Universität Darmstadt},
    year    = {2020}
}
```


## Abstract

One decade after its first published version, the Go programming language has become a popular
and widely-used modern programming language. It aims to achieve thorough memory and
thread safety by using measures such as a strict type system and automated memory management
with garbage collection, which prevents invalid memory access. However, there is also the unsafe
package, which allows developers to deliberately circumvent this safety net. There are a number
of legitimate use cases for doing this, for example, an in-place type conversion saving reallocation
costs to improve efficiency, or interacting with C code through the foreign function interface.

Misusing the unsafe API can however lead to security vulnerabilities such as buffer overflow
and use-after-free bugs. This work contributes an analysis of unsafe usage patterns with respect
to a security context. It reveals possible code injection and information leak vulnerabilities in
proof-of-concept exploits as well as common usages from real-world code.

To assess the risk of unsafe code in their applications, this work presents go-geiger, a novel
tool to help developers quantify unsafe usages not only in their project itself, but including its
dependencies. Using go-geiger, a study on unsafe usage in the top 500 most popular open-source
Go projects on GitHub was conducted, including a manual study of 1,400 individual code samples
on how unsafe is used and for what purpose. The study shows that 5.5% of packages imported
by the projects using the Go module system use unsafe. Furthermore, 38.19% of the projects
use unsafe directly, and 90.96% include unsafe usages through any of their dependencies. A
replication and comparison of a concurrent study by Costa et al. [10] matches these results.

This work further presents go-safer, a novel static code analysis tool that helps developers to
identify two dangerous and common misuses of the unsafe API, which were previously undetected
with existing tools. Using go-safer, 64 bugs in real-world code were identified and patches have
been submitted to and accepted by the maintainers. An evaluation of the tool shows 95.5%
accuracy on the data set of labeled unsafe usages, and 99% accuracy on a set of manually
inspected open-source Go packages.


## Zusammenfassung


## License

Copyright (c) 2020 Johannes Lauinger  

### Thesis Document and Source Code

<a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/4.0/"><img alt="Creative Commons Lizenzvertrag" style="border-width:0" src="https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/4.0/">Creative Commons Attribution-NonCommercial-NoDerivs  4.0 International License</a>.

### Implementation

Licensed under the terms of the <a rel="license" href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GENERAL PUBLIC LICENSE, Version 3</a>.

