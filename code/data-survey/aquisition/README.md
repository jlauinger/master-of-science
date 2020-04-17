# Data Aquisition Tool

Iteration 3


## Usage

```shell script
go build
mkdir -p /path/to/data
mkdir -p /path/to/repositories
```

Download repositories:

```shell script
./data-aquisition projects --download --data-dir=/path/to/data --destination=/path/to/repositories
```

Run analysis steps. You can run one at a time:

```shell script
./data-aquisition analyze grep --offset 0 --length 500 --data-dir=/path/to/data
./data-aquisition analyze vet --offset 0 --length 500 --data-dir=/path/to/data
./data-aquisition analyze gosec --offset 0 --length 500 --data-dir=/path/to/data
```

To do better parallelization, you can split the analysis into buckets. `go vet` already automatically parallelizes as
best as possible, `gosec` also parallelizes pretty well. Therefore, the grep analysis is the one that profits the most.
Even here, `ripgrep` does an excellent parallelization step, but its execution takes less time compared to the overall
program runtime, so chunking can give a little extra optimization.

```shell script
./data-aquisition analyze grep --offset 350 --length 50--data-dir=/path/to/data
...
```

Then, concatenate the resulting CSV files, dropping the headers in all but the first.


## Package or Module?

 - A Go module must be a VCS repository or a VCS repository should contain a single Go module.
 - A Go module should contain one or more packages
 - A package should contain one or more .go files in a single directory.

Read this:

https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16