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

Any vet message:

[tp] lines that were flagged by geiger and vet (any message): 219
[fn] lines that were not flagged by vet: 76738
[fp] lines that were flagged by vet (any message) but not geigered: 31224

Only the vet message that is related (unsafeptr):

[tp] lines that were flagged by geiger and vet (unsafeptr): 213
[fn] lines that were not flagged by vet: 76744
[fp] lines that were flagged by vet (unsafeptr) but not geigered: 0


**gosec**:

Any gosec message:

[tp] lines that were flagged by geiger and gosec (any message): 36279
[fn] lines that were not flagged by gosec: 40678
[fp] lines that were flagged by gosec (any message) but not geigered: 114306

Only the gosec message that is related (unsafe should be audited):

[tp] lines that were flagged by geiger and gosec (only unsafe-related): 36267
[fn] lines that were not flagged by gosec: 40690
[fp] lines that were flagged by gosec (only unsafe-related) but not geigered: 0

Only the gosec message that is related (unsafe should be audited) and geiger unsafe package which is the only usage type
found by gosec:

[tp] lines that were flagged by geiger (only unsafe pkg matches) and gosec (only unsafe-related): 36267
[fn] lines that were not flagged by gosec: 18019
[fp] lines that were flagged by gosec (only unsafe-related) but not geigered: 0


## Discussion

There is no point in calculating precision or recall because these tools are designed completely differently. Since
go-geiger simply counts all unsafe findings, it *should* find many more that go vet. This is simply to illustrate just
how much more findings there are than what gets flagged by Vet.

Gosec on the other hand is similar to go-geiger in that its only unsafe rule is "presence of unsafe call sites".
We can see that there are no false positives. We see a lot of overlap, although there is still a significant amount of
false negatives. Some of that is due to gosec not analyzing usages of uintptr and reflect.SliceHeader / StringHeader,
but after filtering it out there are still more than 18000 false negatives. This is because go-geiger uses a much more
thorough approach to finding unsafe sites, while gosec misses some forms. On the other hand, gosec is not designed
to find all usages so I guess this is okay.
