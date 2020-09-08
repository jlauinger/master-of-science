# Data survey

This directory contains almost everything that is related to my empirical study. The subdirectories are:

 - `acquisition/` contains the Go based data acquisition tool that I used to extract data from open-source projects
 - `analysis/` contains Jupyter Notebooks that I used to analyze the data and create tables and plots
 - `classification` contains a Python Flask application to manually classify Go code snippets
 - `data` contains output data from my tools, as well as my labeled data set of Go unsafe usages and generated plots
 - `dumper` contains a Python tool to dump the labeled Go snippets into their own files
 - `zenodo-submission` contains a script to create a Git repository with all analyzed projects as submodules, as well
   as the Readme file that I put into the record published on Zenodo
