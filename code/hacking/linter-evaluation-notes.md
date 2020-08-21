# How do linters get evaluated?

Can I evaluate `go-safer` to an even better, scientific degree? Let's dive into the case studies!


## Cryptolint

Paper: egele2013 https://sites.cs.ucsb.edu/~chris/research/doc/ccs13_cryptolint.pdf

Cryptolint has 6 rules which are statically analyzed on Dalvik byte code in Android applications.

Evaluation: the authors take a data set of about 11,500 Android applications, run cryptolint and report how many
projects broke each rule.

In a few case studies, the authors manually validate some of the positives. There are no manual evaluations on
false negatives.

There is no hard evaluation on recall / precision / accuracy.

Cryptolint is not actually available open-source, Krüger et al. re-implemented their rules in CrySL.


## CrySL / CogniCryptSAST

Paper: kruger2018 https://bodden.de/pubs/ksa+18crysl.pdf

Rules are developed based on the JCA documentation and then refined through selective discussions with developers.

Evaluation: RQ1 is about precision / recall of CogniCryptSAST. Setup: 50 random apps get evaluated manually and check
against the output of CogniCryptSAST. This gives precision / recall values.

The authors find that in 228 usages in the 50 apps, CogniCryptSAST finds 156 misuses. There are 27 type-state misuses 
with 2 false positives and 4 false negatives, yielding precision and recall of 92% and 86%.
There are 129 constraint warnings with 19 false positives, yielding precision and recall of 85% and 100%.

With RQ2, the authors look into how many *projects* contained a misuse.


## Wickert19 MSR

Paper: wickert2019

A labeled data set of Java cryptography API misuses is presented and precision and recall of FindBugs is evaluated
upon this data set.

Difference to my set: there is no other linter to compare against, go-safer *adds* new linting rules that can only
be manually evaluated.


## Dlint

Paper: gong2015 https://dl.acm.org/doi/pdf/10.1145/2771783.2771809

Dlint is a dynamic Javascript linter. It reimplements some of the rules that ESLint already has, and adds some more.

Evaluation: comparison to ESLint which shows (unsurprisingly) that Dlint finds more problems. This is however a good
baseline recall / precision evaluation. I can not do that because go-safer does not reimplement anything.

The evaluation however is more of a counting of how many problems are identified. That I can do too. The authors then
manually check all warnings for validity. That is the same as I have done.


## Go-Safer plan

For me, I'd use my labeled data set of unsafe usages and check if go-safer finds all and no additional instances of the
slice cast pattern. Then, I see how many projects are affected by go-safer findings. I need to compare to Go Vet, which
finds none of the patterns found by go-safer.

Plan:

 - [ ] Identify data set instances with `cast-header` label
 - [ ] Run go-safer on all of those files
 - [ ] See if there are false negatives or false positives
 - [ ] Probably conclude that there is 100% recall and precision
 