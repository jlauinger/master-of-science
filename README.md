# Identification and Analysis of unsafe.Pointer Usage Patterns in Open-Source Go Code

It's dangerous to Go alone. Take \*this!  
Master Thesis by Johannes Lauinger

Submitted: October 26th, 2020

Advisor: Prof. Dr.-Ing. Mira Mezini  
Supervisors: Anna-Katharina Wickert, M.Sc., Dr. rer. nat. Lars Baumgärtner

Software Technology Group  
Department of Computer Science  
Technische Universität Darmstadt


## Citation

**Master Thesis:**

Lauinger, Johannes Tobias. "Identification and Analysis of unsafe.Pointer Usage Patterns in Open-Source Go Code." M.Sc. Thesis. Technische Universität Darmstadt, 2020.

```latex
@mastersthesis{lauinger2020,
    type    = {M.Sc. Thesis},
    author  = {Lauinger, Johannes Tobias},
    title   = {Identification and Analysis of unsafe.Pointer Usage Patterns in Open-Source Go Code},
    school  = {Technische Universität Darmstadt},
    year    = {2020}
}
```

**Paper:**

Johannes Lauinger, Lars Baumgärtner, Anna-Katharina Wickert, and Mira Mezini. "Uncovering the Hidden Dangers: Finding Unsafe Go Code in the Wild." In *19th IEEE International Conference on Trust, Security and Privacy in Computing and Communications, TrustCom 2020, Gouangzhou, China, December 29, 2020 - January 1, 2021*. IEEE, 2021.

```latex
@inproceedings{lauinger2020,
    author={Lauinger, Johannes and Baumgärtner, Lars and Wickert, Anna-Katharina and Mezini, Mira},
    title={{Uncovering the Hidden Dangers}: {Finding Unsafe Go Code in the Wild}},
    booktitle={19th {IEEE} International Conference on Trust, Security and Privacy
               in Computing and Communications, TrustCom 2020, Guangzhou,
               China, December 29, 2020 -- January 1, 2021},
    publisher={{IEEE}},
    year={2021}
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

Ein Jahrzehnt nach der ersten veröffentlichten Version ist die Programmiersprache Go heute eine
beliebte und weit verbreitete, moderne Sprache. Sie strebt Speicher- und Threadsicherheit durch
Maßnahmen wie ein striktes Typsystem und automatische Speicherverwaltung, die ungültige
Speicherzugriffe verhindert, an. Es gibt allerdings ebenfalls das unsafe Package, eine API, die
es Entwickler\*innen erlaubt, diese Maßnahmen zu umgehen. In manchen Fällen kann dies
gerechtfertigt sein, beispielsweise bei der Konvertierung von Daten in einen anderen Typ, ohne
diese im Speicher zu kopieren, um so die Effizienz des Programms zu steigern, oder um externen
C Code über das Foreign Function Interface zu nutzen.

Eine falsche Benutzung der unsafe API kann jedoch zu Sicherheitsproblemen wie Buffer Overflows
und Use-After-Frees führen. Diese Arbeit analysiert Verwendungsmuster von unsafe Code im
Hinblick auf Sicherheitsrisiken. Dabei werden mögliche Code Injection und Information Leak
Verwundbarkeiten sowohl in Proof-of-Concepts als auch in realem Anwendungscode zu Tage
gebracht.

Um die Risiken durch unsafe Code in Anwendungen abzuschätzen, stellt diese Arbeit go-geiger
vor. Es handelt sich dabei um ein neues Werkzeug, das Entwickler\*innen dabei hilft, unsafe
Nutzungen in Projekten und deren Abhängigkeiten zu finden. Mit go-geiger wird eine Studie zur
Nutzung von unsafe in den 500 beliebtesten Open-Source Go Projekten auf GitHub durchgeführt,
inklusive einer manuellen Analyse von 1,400 individuellen Codestücken in Bezug darauf wie
und zu welchem Zweck unsafe benutzt wird. Die Studie zeigt, dass 5.5% der Packages, die von
Projekten importiert werden, welche das Go Modules System unterstützen, unsafe verwenden.
Darüber hinaus nutzen 38.19% der Projekte unsafe direkt, und 90.96% enthalten unsafe Code
durch ihre Abhängigkeiten. Eine Replikation sowie ein Vergleich mit einer zeitgleichen Studie
von Costa et al. [10] bestätigt diese Ergebnisse.

Weiterhin präsentiert diese Arbeit go-safer, ein neues statisches Analysewerkzeug, das Entwickler\*innen
hilft, zwei gefährliche und häufig vorkommende inkorrekte Verwendungen der unsafe API, die mit
bisher existierenden Tools nicht gefunden werden, zu identifizieren. Mittels go-safer
konnten 64 Fehler in realem Code gefunden und entsprechende Patches eingereicht werden, die
von den Maintainern bestätigt wurden. Eine Evaluation des Tool ergibt eine Accuracy von 95.5%
auf dem Datensatz von unsafe Codezeilen, und 99% Genauigkeit auf händisch analysierten
Open-Source Go Packages.


## License

Copyright (c) 2020 Johannes Lauinger  

### Thesis Document and Source Code

<a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/4.0/"><img alt="Creative Commons Lizenzvertrag" style="border-width:0" src="https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/4.0/">Creative Commons Attribution-NonCommercial-NoDerivs  4.0 International License</a>.

### Implementation

Licensed under the terms of the <a rel="license" href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GENERAL PUBLIC LICENSE, Version 3</a>.

