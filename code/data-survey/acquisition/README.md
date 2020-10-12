# Data Acquisition Tool

This tool is used to analyze the open-source Go projects and gather the evaluation data into CSV files.

This is iteration 3, a Go implementation.


## Available analysis operations

The data acquisition tool supports the following analysis modes:

 - `geiger` runs the logic of `go-geiger` to find all unsafe usages in the code (not in comments) as well as counting
   the number of unsafe usages in each package including its dependencies, and analyzing the dependency tree of a
   project including the import depth of each package. This analysis produces the `packages_0_499.csv` as well as the
   `geiger/geiger_findings_0_499.csv` files. The tool does not call the actual `go-geiger` program to save time.
 - `grep` runs a lexical search for unsafe usage sites using ripgrep. Therefore, it can also find occurrences in
   comments. This analysis was used as a first step before `go-geiger` existed and is not really used now anymore. It
   produces the `linters/grep_findings_0_499.csv` file.
 - `vet` runs `go vet` on all packages in each project and saves the warning message results. This is used to compare
   how many of the geiger findings are also found by go vet. This analysis produces the `linters/vet_findings_0_499.csv`
   file.
 - `gosec` runs `gosec` on all packages in each project and saves the warning message results. This is used to compare
   how many of the geiger findings are also found by gosec. This analysis produces the `linters/gosec_findings_0_499.csv`
   file. 
 - `gosafer` runs the `go-safer` unsafe-focused linter developed by me. This is used to find warnings in all of the
   projects, which is how I identified bugs that I submitted upstream patches for. This analysis produces the
   `linters/gosafer_findings_0_499.csv` file.
 - `ast` is an AST-based analysis similar to `go-geiger`, but in a less mature form. The benefit of this analysis is
   that it produces CSV files containing statements and functions that contain unsafe that could be used to check how
   often unsafe usages correlate within the same statement or function. It produces the `ast/ast_findings_0_499.csv`,
   `ast/statements_0_499.csv` and `ast/functions_0_499.csv` files.


## Tool architecture

This tool is used both for creating the projects data set as well as running analysis steps and generating CSV data
files.

The `cmd` package contains Cobra CLI commands and is responsible for the tool syntax and CLI parameters. For the
projects and analyze commands, which contain sub commands, the top level command contains most of the CLI parameters
because with Cobra commands inherit the parameters of their parents.

The `base` package contains logic that is common to several or all analysis operations. It contains all CSV and JSON
data structure definitions, CSV input/output functions that are responsible for marshalling, code that is responsible
for calling `wc` to count LOC for packages and some helper functions. It also contains the common analysis code, i.e.
identifying projects, taking the skip list into effect, enumerating project dependency packages, analyzing the
dependency tree / import depths etc. The concrete analysis step that needs to be done for each package or project is
passed in using a callback function of type `AnalysisCallback`. Since the callbacks are defined in the respective
analysis packages and only passed as values, the `base` package does not import any of the other packages which is
important to avoid cyclic imports.

The `projects` package contains the project data set scraper and the code to identify root modules and module support.
It is responsible for creating the `projects.csv` file.

The `geiger` packages contains the implementation of the `go-geiger` logic. Like the other analysis steps, it uses the
common logic provided by the base package.

The `ast` package contains the alternative, more basic AST analysis operation.

The `linters` package on the other hand contains the analysis steps that invoke external tools, i.e. the grep, govet,
gosec, and gosafer analysis operations.

Output CSV files are opened as needed and writing to them is all handled by the `base` package, which acts as a de-facto
singleton object for accepting write requests.


## Usage

To use the acquisition tool, first create the data directories for repositories and CSV data:

```shell script
go build
mkdir -p /path/to/data/{linters,geiger,ast}
mkdir -p /path/to/repositories
```

Identify, fork and download the top 500 (or another number) most-starred repositories:

```shell script
./acquisition projects --size=500 --download --fork --access-token=XXX --data-dir=/path/to/data --destination=/path/to/repositories
```

Then check module support and identify the root module of each project:

```shell script
./acquisition projects checkmodule --data-dir=/path/to/data
```

To run the analysis steps, use the appropriate command names:

```shell script
./acquisition analyze geiger --data-dir=/path/to/data --skip abc/def --skip xyz/xyz
./acquisition analyze grep --data-dir=/path/to/data
./acquisition analyze vet --data-dir=/path/to/data
./acquisition analyze gosec --data-dir=/path/to/data
./acquisition analyze gosafer --data-dir=/path/to/data
./acquisition analyze ast --data-dir=/path/to/data
```

The `geiger` analysis step is the most important. It uses `go-geiger` logic to identify unsafe usages by looking at AST
nodes. The `vet` and `gosec` analysis are conducted to compare the performance to existing tools. The others are useful
for data analysis, too.

To do better parallelization, you can split the analysis into buckets. `go vet` already automatically parallelizes as
best as possible, `gosec` also parallelizes pretty well. The geiger and grep analysis profit the most.
Even here, `ripgrep` does an excellent parallelization step, but its execution takes less time compared to the overall
program runtime, so chunking can give a little extra optimization. Not specifying offset and length assumes their
defaults 0 and number of projects (thus including all of them).

```shell script
./acquisition analyze geiger --offset 350 --length 50 --data-dir=/path/to/data
```

Then, concatenate the resulting CSV files, dropping the headers in all but the first.

You can skip projects with the skip argument. It can be applied multiple times. The skip list used for the evaluation
can be seen in the `data-survey/data/projects_with_errors.txt` file.

```shell script
./acquisition analyze grep --data-dir=/path/to/data --skip golang/go --skip avelino/awesome-go
```


## Development

To get the source code and compile the binary, run this:

```
$ git clone https://github.com/stg-tud/thesis-2020-lauinger-code
$ cd thesis-2020-lauinger-code/data-survey/acquisition
$ go build
```


## Notes

### Package or Module?

 - A Go module must be a VCS repository / a VCS repository should contain a single Go module.
 - A Go module should contain one or more packages
 - A package should contain one or more .go files in a single directory.

Read this:

https://medium.com/rungo/anatomy-of-modules-in-go-c8274d215c16


### Previous iterations

Formerly, the acquisition tool was built in Bash, then in Python, then the third version is this Go implementation. The
reason to implement it again was to gain speed. The Bash version would have taken months, the Python version would have
taken weeks, and this Go implementation can analyze the projects in a matter of a few hours.


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
