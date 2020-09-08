# Data Acquisition Tool

Iteration 3 (Go implementation)


## Package or Module?

 - A Go module must be a VCS repository or a VCS repository should contain a single Go module.
 - A Go module should contain one or more packages
 - A package should contain one or more .go files in a single directory.

Read this:

https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16


## Usage

```shell script
go build
mkdir -p /path/to/data/{analysis,ast,classification,lexical}
mkdir -p /path/to/repositories
mkdir -p /path/to/copied/files
```

Download repositories:

```shell script
./acquisition projects --download --data-dir=/path/to/data --destination=/path/to/repositories
```

Run analysis steps. You can run one at a time:

```shell script
./acquisition analyze grep --data-dir=/path/to/data
./acquisition analyze vet --data-dir=/path/to/data
./acquisition analyze gosec --data-dir=/path/to/data
./acquisition analyze ast --data-dir=/path/to/data
./acquisition analyze linter --data-dir=/path/to/data
```

To do better parallelization, you can split the analysis into buckets. `go vet` already automatically parallelizes as
best as possible, `gosec` also parallelizes pretty well. Therefore, the grep analysis is the one that profits the most.
Even here, `ripgrep` does an excellent parallelization step, but its execution takes less time compared to the overall
program runtime, so chunking can give a little extra optimization. Not specifying offset and length assumes their
defaults 0 and 500.

```shell script
./acquisition analyze grep --offset 350 --length 50--data-dir=/path/to/data
```

You can skip projects with the skip argument. It can be applied multiple times.

```shell script
./acquisition analyze grep --data-dir=/path/to/data --skip golang/go --skip avelino/awesome-go
```

It is recommended to copy the vulnerable files into a specific directory. The resulting path will be written into the
CSV findings files, and later analysis can use those files to do context expansion on the finding context.

```shell script
./acquisition analyze grep --data-dir=/path/to/data --copy --copy-destination=/path/to/copied/files
```

Then, concatenate the resulting CSV files, dropping the headers in all but the first.


## Development

To get the source code and compile the binary, run this:

```
$ git clone https://github.com/stg-tud/thesis-2020-lauinger-code
$ cd thesis-2020-lauinger-code/data-survey/acquisition
$ go build
```
