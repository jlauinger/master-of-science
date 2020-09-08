# go-safer evaluation step 2: evaluation with manually analyzed packages

## Packages to analyze

Choose six packages: 2 with few, medium, and many unsafe usages each, and mixing small and large packages by LOC.

| **Package**                             | **Directory**                                                                               | **LOC** | **Number Go Files** | **Unsafe Usages** |
|-----------------------------------------|---------------------------------------------------------------------------------------------|---------|---------------------|-------------------|
| k8s.io/kubernetes/pkg/apis/core/v1      | /root/download/kubernetes/kubernetes/pkg/apis/core/v1                                       | 10,048  | 6                   | 675               |
| gorgonia.org/tensor/native              | /root/go/pkg/mod/gorgonia.org/tensor@v0.9.6/native                                          | 1,867   | 4                   | 151               |
| github.com/anacrolix/mmsg/socket        | /root/go/pkg/mod/github.com/anacrolix/mmsg@v1.0.0/socket                                    | 3,782   | 86                  | 114               |
| github.com/cilium/ebpf                  | /root/go/pkg/mod/github.com/cilium/ebpf@v0.0.0-20191113100448-d9fb101ca1fb                  | 2,851   | 14                  | 58                |
| golang.org/x/tools/internal/event/label | /root/go/pkg/mod/golang.org/x/tools@v0.0.0-20200502202811-ed308ab3e770/internal/event/label | 213     | 1                   | 6                 |
| github.com/mailru/easyjson/jlexer       | /root/go/pkg/mod/github.com/mailru/easyjson@v0.7.0/jlexer                                   | 1,234   | 4                   | 5                 |

These packages all contain unsafe usages, some contain slice header conversions as looked at in part 1 of the evaluation.
There are two packages with less than 10 usages, three around 50 to 200 usages and one with more than 500. Furthermore, there is
a large package with more than 10k LOC, 4 around 1k to 4k LOC and one small around 200 LOC.

In this study, I only analyze the package on its own, no dependencies. I also look at all the files except test files,
other than go-safer which will only look at the current architecture files.


## Evaluation

Plan:

 - [x] use a script that executes `grep` to find usages of `unsafe.`, `reflect.` and `uintptr` and writes those to a CSV
       file
 - [x] when manually going through the code, annotate that CSV file manually
 - [x] run `go-safer`, `go vet` and `gosec`, capture output and write all findings in a CSV file
 - [x] then use Pandas to calculate the resulting PRF scores
 - [x] when `go-safer` etc. find something that was not even listed in the unsafe usages, add that finding as negative to
       be correctly counted as a false positive
 - [x] when looking through the code, when there is a finding that is outside of an unsafe usage line, add that finding as
       positive
       
Evaluated packages:

 - [x] jlexer
 - [x] label
 - [x] native
 - [x] v1
 - [x] ebpf
 - [x] socket


## Evaluation summary

| **Package**                              | **TP**                |            |           | **FP**                 |            |           | **TN**                |            |           | **FN**                 |            |           | **Precision** |            |           | **Recall**    |            |           | **Accuracy** |            |           |
|------------------------------------------|-----------------------|------------|-----------|------------------------|------------|-----------|-----------------------|------------|-----------|------------------------|------------|-----------|---------------|------------|-----------|---------------|------------|-----------|--------------|------------|-----------|
|                                          | **go-safer**          | **vet**    | **gosec** | **go-safer**           | **vet**    | **gosec** | **go-safer**          | **vet**    | **gosec** | **go-safer**           | **vet**    | **gosec** | **go-safer**  | **vet**    | **gosec** | **go-safer**  | **vet**    | **gosec** | **go-safer** | **vet**    | **gosec** |
| k8s.io/kubernetes/pkg/apis/core/v1       | 0                     | 0          | 0         | 0                      | 0          | 676       | 676                   | 676        | 1         | 0                      | 0          | 0         | -             | -          | 0         | -             | -          | -         | 1            | 1          | 0.001     |
| gorgonia.org/tensor/native               | 48                    | 0          | 0         | 9                      | 0          | 98        | 101                   | 109        | 11        | 0                      | 48         | 48        | 0.842         | -          | 0         | 1             | 0          | 0         | 0.943        | 0.694      | 0.070     |
| github.com/anacrolix/mmsg/socket         | 0                     | 0          | 0         | 0                      | 0          | 17        | 115                   | 115        | 99        | 0                      | 0          | 0         | -             | -          | 0         | -             | -          | -         | 1            | 1          | 0.853     |
| github.com/cilium/ebpf                   | 0                     | 0          | 0         | 1                      | 0          | 52        | 57                    | 58         | 27        | 0                      | 0          | 0         | 0             | -          | 0         | -             | -          | -         | 0.983        | 1          | 0.342     |
| golang.org/x/tools/internal/event/label  | 0                     | 0          | 0         | 0                      | 0          | 7         | 5                     | 5          | 1         | 0                      | 0          | 0         | -             | -          | 0         | -             | -          | -         | 1            | 1          | 0.125     |
| github.com/mailru/easyjson/jlexer        | 1                     | 0          | 0         | 0                      | 0          | 2         | 4                     | 4          | 2         | 0                      | 1          | 1         | 1             | -          | 0         | 1             | 0          | 0         | 1            | 0.8        | 0.4       |
| **total**                                | **49**                | **0**      | **0**     | **10**                 | **0**      | **852**   | **958**               | **967**    | **141**   | **0**                  | **49**     | **49**    | **0.831**     | **-**      | **0**     | **1**         | **0**      | **0**     | **0.990**    | **0.952**  | **0.135** |
