# Evaluation

This directory contains mostly notes on evaluation steps and results, as well as some general research.


## Assembly analysis

The `asm-analysis/` directory contains an evaluation into whether there is a difference between unsafely creating a
slice header by composite literal and then casting it to a slice in one or two statements. It contains Go files for
several possibilities, as well as output assembly. The assembly is divided into Go intermediate assembly (S.txt files) 
and actual amd64 assembly extracted from the linked ELF files (ELF.txt files)


## Study comparison with related work by Costa et al

The `costa-study-comparison/` directory contains a replication of the similar study about unsafe usage in Go code by
Costa et al, which is a relevant related work. It also contains comparison notes and results.


## go-safer evaluation

The `go-safer-evaluation/` directory contains the scientific evaluation of the `go-safer` linter tool. It is split into
two parts, an evaluation based on the labeled data set and an evaluation based on a thorough analysis of six complete
Go packages. It also contains notes about the evaluation strategy.


## Comparison with other external linters

The `other-linters-comparison/` directory contains comparison results about how many findings of `go-geiger` are also
found by the existing tools `go vet` and `gosec`.
