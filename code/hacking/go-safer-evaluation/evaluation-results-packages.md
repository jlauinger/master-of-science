# go-safer evaluation step 2: evaluation with manually analyzed packages

## Packages to analyze

Choose six packages: 2 with few, medium, and many unsafe usages each, and mixing small and large packages by LOC.

| **Package**                        | **Directory**                                                                       | **LOC** | **Number Go Files** | **Unsafe Usages** |
|------------------------------------|-------------------------------------------------------------------------------------|---------|---------------------|-------------------|
| k8s.io/kubernetes/pkg/apis/core/v1 | /root/download/kubernetes/kubernetes/pkg/apis/core/v1                               | 10,048  | 6                   | 675               |
| gopkg.in/olebedev/go-duktape.v3    | /root/go/pkg/mod/gopkg.in/olebedev/go-duktape.v3@v3.0.0-20200316214253-d7b0ff38cac9 | 2,687   | 7                   | 125               |
| github.com/cilium/cilium/pkg/bpf   | /root/download/cilium/cilium/pkg/bpf                                                | 2,851   | 13                  | 98                |
| github.com/tsg/gopacket/pcap       | /root/go/pkg/mod/github.com/tsg/gopacket@v0.0.0-20190320122513-dd3d0e41124a/pcap    | 1,009   | 3                   | 28                |
| github.com/elastic/go-perf         | /root/go/pkg/mod/github.com/elastic/go-perf@v0.0.0-20191212140718-9c656876f595      | 3,400   | 5                   | 27                |
| github.com/karalabe/usb            | /root/go/pkg/mod/github.com/karalabe/usb@v0.0.0-20190919080040-51dc0efba356         | 1,031   | 7                   | 20                |

These packages all contain unsafe usages, some contain slice header conversions as looked at in part 1 of the evaluation.
There are some packages around 30 usages, two with around 100 usages and one with more than 500. Furthermore, there is
a large package with more than 10k LOC, 3 around 3k LOC and two small around 1k LOC.

In this study, I only analyze the package on its own, no dependencies.


## Evaluation summary

| **Package**                        | **TP**                |            |           | **FP**                 |            |           | **TN**                |            |           | **FN**                 |            |           | **Recall**   |            |           | **Precision** |            |           | **Accuracy** |            |           |
|------------------------------------|-----------------------|------------|-----------|------------------------|------------|-----------|-----------------------|------------|-----------|------------------------|------------|-----------|--------------|------------|-----------|---------------|------------|-----------|--------------|------------|-----------|
|                                    | **go-safer**          | **vet**    | **gosec** | **go-safer**           | **vet**    | **gosec** | **go-safer**          | **vet**    | **gosec** | **go-safer**           | **vet**    | **gosec** | **go-safer** | **vet**    | **gosec** | **go-safer**  | **vet**    | **gosec** | **go-safer** | **vet**    | **gosec** |
| k8s.io/kubernetes/pkg/apis/core/v1 |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |            |           |
| gopkg.in/olebedev/go-duktape.v3    |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |            |           |
| github.com/cilium/cilium/pkg/bpf   |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |            |           |
| github.com/tsg/gopacket/pcap       |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |            |           |
| github.com/elastic/go-perf         |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |            |           |
| github.com/karalabe/usb            |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |            |           |
