# Data Acquisition Tool

Iteration 3 (Go implementation)

This tool is used to analyze the open-source Go projects and gather the evaluation data into CSV files.


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


## Notes

### Package or Module?

 - A Go module must be a VCS repository or a VCS repository should contain a single Go module.
 - A Go module should contain one or more packages
 - A package should contain one or more .go files in a single directory.

Read this:

https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16


### Previous iterations

Formerly, the acquisition tool was built in Bash, then in Python, then the third version is this Go implementation.


### Use grep to find out how many unsafe.Pointers exist in module dependencies

```
for dir in $(go mod vendor -v 2>&1 | grep -v "#" | sort | uniq); do 
lines=$(rg unsafe.Pointer vendor/$dir | wc -l); 
echo "$lines $dir"; 
done | sort -n | grep -ve "^0 "
```

There are some ways to show the dependency modules:

 - `go mod vendor -v`: stores the required dependency modules in the `vendor` directory, therefore I can analyze the
   exact version that is used by the module. The `-v` switch prints the modules (and thus their filesystem path) to
   stderr.
 - `go mod vendor && find . -type d vendor`: simply checks all directories in the vendor directory. This is very
   inaccurate because not every directory is a package. The first approach prints the logical modules rather than the
   low-level directories and is thus better.
 - `go list -m all`: using the `-m` flag, the list command prints module information. The `all` keyword instructs it
   to recursively show all dependency modules. This command output includes the versions of the modules, which is very
   nice. Some preprocessing using `cut` is necessary to strip the version and including module from the directory path.
   The problem is that this command lists modules that do not exist in the vendor directory, making the grepping throw
   errors. TODO: find out why this is the case.
 - `go mod graph`: This is very similar to the command above, and shares the problem that not all directories exist. 
