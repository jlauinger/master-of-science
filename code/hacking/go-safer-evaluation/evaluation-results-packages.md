# go-safer evaluation step 2: evaluation with manually analyzed packages

## Packages to analyze

Choose six packages: 2 with few, medium, and many unsafe usages each, and mixing small and large packages by LOC.

| **Package**                        | **Directory**                                                                       | **LOC** | **Number Go Files** | **Unsafe Usages** |
|------------------------------------|-------------------------------------------------------------------------------------|---------|---------------------|-------------------|
| gopkg.in/olebedev/go-duktape.v3    | /root/go/pkg/mod/gopkg.in/olebedev/go-duktape.v3@v3.0.0-20200316214253-d7b0ff38cac9 | 313     | 3                   | 906               |
| k8s.io/kubernetes/pkg/apis/core/v1 | /root/download/kubernetes/kubernetes/pkg/apis/core/v1                               | 10,048  | 6                   | 675               |
| github.com/karalabe/usb            | /root/go/pkg/mod/github.com/karalabe/usb@v0.0.0-20190919080040-51dc0efba356         | 166     | 2                   | 106               |
| github.com/cilium/cilium/pkg/bpf   | /root/download/cilium/cilium/pkg/bpf                                                | 2,851   | 13                  | 98                |
| github.com/tsg/gopacket/pcap       | /root/go/pkg/mod/github.com/tsg/gopacket@v0.0.0-20190320122513-dd3d0e41124a/pcap    | 100     | 1                   | 28                |
| github.com/elastic/go-perf         | /root/go/pkg/mod/github.com/elastic/go-perf@v0.0.0-20191212140718-9c656876f595      | 3,400   | 5                   | 27                |

These packages all contain unsafe usages, some contain slice header conversions as looked at in part 1 of the evaluation.
There are 2 packacges with more than 500, 2 around 100, and 2 projects around 30 unsafe usages, and each group except for
the first one also has one large and one small project.

In this study, I only analyze the package on its own, no dependencies.


## Evaluation summary

| **Package**                        | **True Positives TP** |            |           | **False Positives FP** |            |           | **True Negatives TN** |            |           | **False Negatives FN** |            |           | **Recall**   |            |           | **Precision** |            |           | **Accuracy** |
|------------------------------------|-----------------------|------------|-----------|------------------------|------------|-----------|-----------------------|------------|-----------|------------------------|------------|-----------|--------------|------------|-----------|---------------|------------|-----------|--------------|
|                                    | **go-safer**          | **go vet** | **gosec** | **go-safer**           | **go vet** | **gosec** | **go-safer**          | **go vet** | **gosec** | **go-safer**           | **go vet** | **gosec** | **go-safer** | **go vet** | **gosec** | **go-safer**  | **go vet** | **gosec** | **go-safer** |
| gopkg.in/olebedev/go-duktape.v3    |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |
| k8s.io/kubernetes/pkg/apis/core/v1 |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |
| github.com/cilium/cilium/pkg/bpf   |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |
| github.com/karalabe/usb            |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |
| github.com/tsg/gopacket/pcap       |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |
| github.com/elastic/go-perf         |                       |            |           |                        |            |           |                       |            |           |                        |            |           |              |            |           |               |            |           |              |