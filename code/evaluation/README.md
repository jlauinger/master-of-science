# Evaluation

This directory contains mostly notes on evaluation steps and results, as well as some general research.


## Assembly analysis

The `asm-analysis/` directoy contains an evaluation into whether there is a difference between unsafely creating a
slice header by composite literal and then casting it to a slice in one or two statements. It contains Go files for
several possibilities, as well as output assembly. The assembly is divided into Go intermediate assembly (S.txt files) 
and actual amd64 assembly extracted from the linked ELF files (ELF.txt files)


## 