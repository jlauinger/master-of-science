# Identification and analysis of unsafe.Pointer usages in open-source Go code

![CI](https://github.com/stg-tud/thesis-2020-Lauinger/workflows/CI/badge.svg)

Modern programming languages like Rust and Go have mechanisms to protect potential
unsafe usages, e.g., derefences of raw pointers or modifying static variables. Thus, it is
recommended to avoid unsafe usages. However, if a developer wants to avoid unsafe
usages and potential security vulnerabilities caused by these usages, they do not only need
to check their code but also their dependencies.

To provide a developer the power to understand the unsafe usages within their code base,
tools like cargo-geiger [0] exists. Unfortunately, such a tool does not exist by today for Go.
Thus, an in-depth analysis of how many Go projects include direct and indirect unsafe
usages and if these projects are vulnerable does not exist, e.g., to buffer overflows [2,3,4].
Within this work, the aim is to develop a tool, probably based on go vet [1], that can identify
unsafe usages in Go projects - similar to cargo-geiger. Based on this tool, the thesis should
evaluate how common unsafe usages in Go are and try to analyze if vulnerabilities caused
by unsafe usages exist.

[0] ​ https://github.com/anderejd/cargo-geiger  
[1] ​ https://golang.org/cmd/vet/  
[2] Larochelle, D., & Evans, D. (2001). Statically detecting likely buffer overflow
    vulnerabilities. In ​ 10th USENIX Security Symposium.  
[3] Alnaeli, S. M., Sarnowski, M., Aman, M. S., Abdelgawad, A., & Yelamarthi, K. (2017).
    Source Code Vulnerabilities in IoT Software Systems. ​ Advances in Science, Technology and
    Engineering Systems Journal ​ , ​ 2 ​ (3), 1502-1507.  
[4] Wang, C., Zhang, M., Jiang, Y., Zhang, H., Xing, Z., & Gu, M. Escape from Escape
    Analysis of Golang. ICSE 2020.  


## Important Links

The TUDa Corporate Design and fonts: https://www.intern.tu-darmstadt.de/arbeitsmittel/corporate_design_vorlagen/index.de.jsp

The github page of the new TUDa Latex classes: https://github.com/tudace/tuda_latex_templates

The new TUDa Latex classes: https://www.ce.tu-darmstadt.de/ce/latex_tuda/index.de.jsp

(The old TUDa Latex classes: http://exp1.fkp.physik.tu-darmstadt.de/tuddesign/)

Information about theses can be found on the [webpage of the FB 20](https://www.informatik.tu-darmstadt.de/studium_fb20/im_studium/studienbuero/abschlussarbeiten_fb20/index.de.jsp).

Information on how to submit your final thesis in an electronic form can be found on the [FAQ of the university](https://www.tu-darmstadt.de/studieren/tucan_studienorganisation/tucan_faq/details_96256.de.jsp).

Information about plagiarism and scientific ethic are available on the [webpage of the FB 20](https://www.informatik.tu-darmstadt.de/studium_fb20/im_studium/studienbuero/plagiarismus/index.de.jsp). 

## Other Material to Write a Good Thesis

*  “How to write a successful Bachelor’s/Master’s thesis” by Elmar Jürgens from TUM <https://thesisguide.org/>
*  ["Writing academic papers"](https://sarahnadi.org/writing-papers/) by Sarah Nadi from [University of Alberta](https://sarahnadi.org/smr/), former post-doc in our group. 


## License

Copyright (c) 2020 Johannes Lauinger

<a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/4.0/"><img alt="Creative Commons Lizenzvertrag" style="border-width:0" src="https://i.creativecommons.org/l/by-nc-nd/4.0/88x31.png" /></a><br />This work is licensed under a <a rel="license" href="http://creativecommons.org/licenses/by-nc-nd/4.0/">Creative Commons Attribution-NonCommercial-NoDerivs  4.0 International License</a>.


