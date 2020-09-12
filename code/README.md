# Identification and analysis of unsafe.Pointer usages in open-source Go code: Implementation

**It's dangerous to Go alone. Take \*this!**

This is the implementation repository for my Master's thesis "Identification and Analysis of unsafe.Pointer Usage 
Patterns in Open-Source Go Code".


## Repository structure

This section explains which directories contain which part of the implementation and evaluation. All directories also
contain individual README files explaining their structure and purpose.

 - `ansible/` contains configuration management code to setup the virtual machine that I used to gather my evaluation
   data
 - `blog-posts/` contains four blog posts showcasing possible unsafe vulnerabilities that I published on the DEV
   community
 - `data-survey/` contains a tool to analyze my data set of 500 open-source Go projects and write them into CSV data
   files, various Jupyter notebooks that I used to analyze the data, a web application to manually label unsafe code
   usages to obtain a labeled data set of 1,400 usage sites, almost all data files (although excluded from Git), as well
   as meta data for my data set submission to Zenodo.
 - `evaluation/` contains various evaluation results, including the evaluation of go-safer, a replication of and
   comparison to a related, concurrent work by Costa et al., and assembler-level analysis of a specific unsafe misuse
 - `exploits/` contains proof-of-concept code to show different possible vulnerabilities that are introduced through
   a misuse of unsafe code
 - `go-geiger/` contains the source code for go-geiger, a tool to identify and count unsafe usage sites in Go packages
   and their dependencies
 - `go-safer` contains the source code for go-safer, an unsafe-focused linter for Go code that detects two misuses of
   unsafe.Pointer that can lead to security vulnerabilities and were previously undetected with existing tools


## Data Access

Jupyter Notebook Server reachable from within the TU Darmstadt network:

[http://vm6.rbg.informatik.tu-darmstadt.de/notebooks](http://vm6.rbg.informatik.tu-darmstadt.de/notebooks)

HTTP Basic Auth:

 - Benutzer: `johannes`
 - Passwort: `Syp9393`
 
Jupyter server password, if needed: `admin`

All data is stored on the server in the  `/root/data` directory.

Python code to load the CSV data:

```python
projects_df = pd.read_csv('/root/data/projects.csv',
                         parse_dates=['project_created_at', 'project_last_pushed_at', 'project_updated_at'])

grep_df = pd.read_csv('/root/data/lexical/grep_findings_0_499.csv')
package_df = pd.read_csv('/root/data/packages_0_499.csv')

vet_df = pd.read_csv('/root/data/lexical/vet_findings_0_499.csv')
gosec_df = pd.read_csv('/root/data/lexical/gosec_findings_0_499.csv')
ast_df = pd.read_csv('/root/data/ast/ast_findings_0_499.csv')
function_df = pd.read_csv('/root/data/ast/functions_0_499.csv')
statement_df = pd.read_csv('/root/data/ast/statements_0_499.csv')
```


## License

Copyright (c) 2020 Johannes Lauinger

Licensed under the terms of the <a rel="license" href="https://www.gnu.org/licenses/gpl-3.0.en.html">GNU GENERAL PUBLIC LICENSE, Version 3</a>.
