# Comparison with study of Costa et al.

Costa et al. did a similar study on the prevalence of unsafe code blocks in Go and published in the IEEE Transactions
on Software Engineering: https://arxiv.org/abs/2006.09973

I try to replicate the study to the best of my possibilities and compare the results.


## Replication package

https://zenodo.org/record/3871931#.XzrconVfjQr

It does contain:
 
 - the rawdata (unsafe usages) and data (authors' interpretations and labelings) as CSV files
 - Python scripts to build plots and count numbers based on the data CSV files
 
So far, similar to my own replication package for our paper. However, it does not contain:

 - the raw project data
 - the scripts to find unsafe usages in the code (sadly). The authors say that this is not included because it would
   not make much sense to have it without the project code.
   
There are further problems with the replication package:

 - the authors state that their data set was collected at 01.10.2019, but the data does not include the specific
   revisions of the projects under examination. There are files with labeled unsafe usages, but they do not contain
   the code but a link to Github, which points to the `master` branch (outdated by now).
 - the labeling is rather incomplete and authors disagree a lot
 - the documentation is very bad, I don't understand most of the CSV files
 

## Comparison

From an overall look, I find the following similarities:

 - the label classes that the authors use are pretty much the same as I use. They split CGo and atomic/sync into two
   classes, which I combined into FFI, and they combined Cast while I split into different classes. Other than that,
   they only have an additional "find architectural details" class which I put into memory layout.
   
And these differences:

 - Costa data set contains a set of unsafe related issues and pull requests (about 280), which I did not collect
   
   
## Replication approaches

 1. Run the Costa et al. script on my data set (downloaded projects) and then use my analysis scripts to check my
    numbers. **Impossible**, because the scripts to extract unsafe usages are not provided.
 2. Run my scripts on the Costa data set. **Impossible**, because it does not provide exact project revisions.
 3. Recount fraction of projects using unsafe with a subset of Costa data, the projects that I analyzed too. **Results**:
    305 of 343 projects of my data set were analyzed by Costa, too. The others are not included in their data set.
    Of those, there were 7 where I found a few usages and Costa did not, and 25 where Costa found something and I did
    not. However, if I remove the filter to include only direct project code, these 25 go away. Therefore, I strongly
    think that the differences are not in the data collection, but in the filtering of 3rd-party / 1st-party code. I
    can not validate this because the data set of Costa's is already filtered to only include their root projects. For
    86 projects, my and Costa's count of found unsafe code blocks differ. The difference in number of found unsafe code 
    blocks is 50 on average, with about 8 projects with notable differences of about 100 while the total block count is
    at about 200-300. However, in general the numbers align very well and the differences are very likely to be caused
    by different filtering of main project packages. 
 4. Compare fraction of projects transitively using unsafe. **Impossible**, because Costa filtered everything so that
    only the root projects are contained in the data
 5. Compare the manual labeling. **Results**: comparison is only possible for 20 instances. Many even lack labels on
    Costa side. For the four actual matches, the labeling is perfectly the same. See the other document for details.
 6. Costa data set contains about 7500 Vet errors "possible misuse of unsafe pointer", did I find those?