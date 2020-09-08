# Labeled data set dumper script

This dumper script is used to generate a file for each labeled snippet in the labeled data set of unsafe usages in
open-source code, containing the labels as well as information about the snippet source and context lines.


## Usage

To run the dumper script, first install dependencies:

```shell
python3 -m venv venv
source venv/bin/activate
```

Then run the dumper like this:

```shell
GO_MOD_PATH=/root/go GO_LIB_PATH=/usr/local/go ./dumper.py \
    ../../../data/classification/sampled_usages_std.csv ../../../data/classified/std
GO_MOD_PATH=/root/go GO_LIB_PATH=/usr/local/go ./dumper.py \
    ../../../data/classification/sampled_usages_app.csv ../../../data/classified/app
```

The environment variables `GO_MOD_PATH` and `GO_LIB_PATH` need to be set to the respective directories containing
Go modules and the Go standard library packages. Usually, these are GOPATH and GOROOT.

Arguments to the script are the CSV file containing the labeled snippets and the directory to dump files to.


## Output format

The output of this script is one file for each labeled snippet, organized in directories representing the labels given
to the snippet. Each snippet has two labels: one indicating the purpose / goal of the operation, and one indicating what
is being done. Directories containing snippet files are named like this: `purpose__what`, with both labels contained in
the name separated by two underscores. An example would be `efficiency__cast-struct`.

The samples are provided as one file for each sample. The file name is a hash of line number, file, package etc. of the 
finding, providing a guaranteed unique name. The files contain 4 sections divided by dashes. The first section provides 
information about the module, version, package, file, and line of the snippet. It also states which project included
this snippet (but there can be more projects in the data set that share usage of the snippet), and the labels as 
already included in the directory name. The information is guaranteed to be in the same line number across files. 
The second section contains the snippet code line. The third and fourth section contain a +/- 5 lines and +/- 100 lines 
context, respectively.
