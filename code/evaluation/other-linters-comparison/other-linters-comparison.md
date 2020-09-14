# Comparison with go vet and gosec

This file contains evaluation results for a comparison of the `go-geiger` findings with the existing tools `go vet` and
`gosec`.


## Approach

Using the data acquisition tool, I run `go vet` and `gosec` on all the projects in the `go-geiger` evaluation, that is
343 of the top 500 most-starred Go projects, and saved all the warning messages into CSV files.

Then I merged the resulting data sets and identified which lines in which files were flagged by both `go-geiger` and
the existing tools. I just count how many of those lines there are.
 
There is no need to compare to a function unsafe analysis done by my other tool `go-safer`, because there is a separate
specialized `go-safer` evaluation which is done on my data set of labeled unsafe usages and on a subset of 6 selected
packages that I manually analyzed completely.


## Results

**go vet**:

[tp] lines that were flagged by geiger and vet (unsafeptr): 213
[fn] lines that were not flagged by vet: 76744
[fp] lines that were flagged by vet (unsafeptr) but not geigered: 0


**gosec**:

tbd


## Discussion

There is no point in calculating precision or recall because these tools are designed completely differently. Since
go-geiger simply counts all unsafe findings, it *should* find many more that go vet. This is simply to illustrate just
how much more findings there are than what gets flagged by Vet.

Gosec on the other hand is similar to go-geiger in that its only unsafe rule is "presence of unsafe call sites".
