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
   
Differences in project selection: within 6 months of time between data collection, the stars counts can change a lot!   
I have no possibility to compare the star situation at Costa collect time.

   
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
    Costa side. Also Costa labels files, not lines. For the four actual matches, the labeling is perfectly the same. 
    See the other document for details.
 6. Costa data set contains about 7500 Vet errors "possible misuse of unsafe pointer", did I find those? **Results**:
    yes, I found about 60k of those Vet errors. It is only about 200 after deduplication, but Costa would also need
    deduplication and also the data sets do not overlap that much anyways.
    
    
## Unsafe count differences evaluation

There are 86 projects with a different unsafe count. 5 of those differ around 100 usages.

Top 10 projects to compare, with an unsafe difference ranging between 49 and 234.

Example compare: `grep -- '^+' v1.16.1...fb9e1946b0.diff | grep -v '^+++' | grep -Pno 'unsafe\.' | wc -l`

**jetstack/cert-manager**: 

difference is +234.

Project description: Automatically provision and manage TLS certificates in Kubernetes

The diff between v0.11.0-alpha.0 (27.09.2019) and 78ee463a98 (28.05.2020) shows 248 additions and 8 
deletions of `unsafe`, this explains a difference of +240. Dependending on when exactly the code is
downloaded in Costa, this matches the recorded difference.

The changes add and remove (change) usages of unsafe exclusively for efficient in-place conversion between struct types.

When running my tool at v.0.11.0-alpha.0, I get a difference of -4, which fits the suggestion by the diff reasonably
well. I would have expected -6 here.

Therefore, I can safely assume that the counts match for this project.

 
**kubernetes/kubernetes**:

difference is -173.

Project description: Production-Grade Container Scheduling and Management

The diff between v1.16.1 (02.10.2019) and fb9e1946b0 (28.05.2020) shows 951 additions and 331 deletions of `unsafe`,
meaning there would be a difference of +620.

All of those unsafe additions are within generated code, so the difference might be caused by Costa et al ignoring
this code? No, this does not explain it, as it is actually exactly wrong. If Costa had ignored the generated code, then
I should have seen an even bigger positive difference.

When running my tool at v1.16.1, I get a difference of -370. This means that I am probably not counting some modules
that are part of the repository but not the root module, in fact there are 26 `go.mod` files in the repository
excluding the vendor directory.
Then the difference gets smaller over time since the diff suggests there are unsafe usages being added. 

Using grep, I can manually verify that I get 2065 matches (reasonably close to the 2058 reported by Costa) with plain
search in the repository: `rg 'unsafe\.' -g '!vendor' -g '*.go' | wc -l`.
The difference to my implementation here is that it double-counts code that is available for multiple architectures,
while I don't.


**golang/mobile**: 

difference is -118.

Project description: Go on Mobile

The diff between 6d0d39b (02.10.2019) and 4c31acba00 (28.05.2020) shows 2 additions and 1 deletion of `unsafe`, so this
really doesn't explain the difference at all.

When running my tool at 6d0d39b, I get a difference of -119 exactly as the diff would suggest.
Additionally, there is only one `go.mod` file in the repository.

Using grep, I can manually verify that I get 207 matches (reasonably close to the 211 reported by Costa) with plain
search in the repository: `rg 'unsafe\.' -g '!vendor' -g '*.go' | wc -l`.
The difference to my implementation here is that it double-counts code that is available for multiple architectures,
while I don't.


**TykTechnologies/tyk**: 

difference is +108.

Project description: Tyk Open Source API Gateway written in Go

The diff between v2.8.5 (01.10.2019) and ce0ee257b6 (28.05.2020) shows 4002 additions and 1757 deletions of unsafe. 
In version v2.8.5, there isn't a `go.mod` file yet, instead there is a vendor directory inside the directory.

I cannot run my tool at v2.8.5 because of the missing `go.mod` file. 

Running grep shows an unsafe count of 36, which is still different from the 79 reported by Costa.
Therefore, I assume there are differences in vendoring here.

TODO: refine?


**elastic/beats**: 

difference is -94.

Project description: Beats - Lightweight shippers for Elasticsearch & Logstash
                     
The diff between v7.4.0 (01.10.2019) and df6f2169c5 (28.05.2020) shows 3367 additions and 775 deletions of unsafe.
In version v7.4.0, there isn't a `go.mod` file yet, instead there is a vendor directory inside the directory. 

Using grep, excluding vendor and test files as well as comments, but counting all occurences even on the same line with
`rg '^[^/]*unsafe\.' -g '!vendor' -g '*.go' -g '!*_test.go' | rg 'unsafe\.' -o | wc -l`, I can manually verify that
the repository then contains 160 usages, which is reasonably close to the 164 reported by Costa. The difference to my
number six months later is then due to different vendoring.



**golang/tools**: 

difference is -69.

Project description: Go Tools golang.org/x/tools

The diff between c337991 (30.09.2019) and 6be401e3f7 (28.05.2020) shows 13 additions and 1 deletion of unsafe, therefore
there is a difference of +12 unsafe usages.

When running my tool at c337991, I see a difference of -77, which fits the suggestion by the diff reasonably well, I
would have expected -81.

Using grep, I can manually verify that at c337991, counting unsafe matches excluding those in comments etc. using 
`rg '^[^/]*unsafe\.' -g '!vendor' -g '*.go' | wc -l`, there are 111 usages which is reasonably close to the 113 as
reported by Costa. This is because this repository contains lots of architecture-duplicated code which I ignore.


**peterq/pan-light**: 

difference is -64.

Project description: 百度网盘不限速客户端, golang + qt5, 跨平台图形界面

The diff between 482eb093f (31.08.2019) and 867eee7a92 (28.05.2020) shows no additions or deletions of unsafe.
Additionally, there are no dependencies whatsoever.

When running my tool at 482eb093f, I still see a difference of -64 as the diff also suggests.

This repository contains 7 `go.mod` files excluding the demo directory, which means my count of 0 is explained by
submodules in the repository which Costa counted.


**cilium/cilium**: 

difference is +55.

Project description: eBPF-based Networking, Security, and Observability

The diff between v1.6.2 (25.09.2019) and 9b0ae85b5f (28.05.2020) shows 2689 additions and 2485 deletions, therefore I
should see a difference of +204.
In version v1.6.2, there isn't a `go.mod` file yet, instead there is a vendor directory inside the directory.

I cannot run my tool at v1.6.2 because of the missing `go.mod` file.

Using grep, excluding vendor and test files as well as comments, but counting all occurences even on the same line with
`rg '^[^/]*unsafe\.' -g '!vendor' -g '*.go' -g '!*_test.go' | rg 'unsafe\.' -o | wc -l`, I can manually verify that
the repository then contains 286 usages, which is reasonably close to the 290 reported by Costa. The difference to my
number six months later is then due to different vendoring.


**go-delve/delve**: 

difference is -52.

Project description: Delve is a debugger for the Go programming language

The diff between v1.3.0 (28.08.2019) and 4a9b3419d1 (28.05.2020) shows 56 additions and 35 deletions, therefore I
should see a difference of +21.

When running my tool at v1.3.0, I see a difference of -32. Following the diff, I should have seen -73 actually.

Using grep, I can manually verify that the repository contains 70 usages with `rg '^[^/]*unsafe\.' -g '!vendor' -g '*.go' | wc -l`,
which reasonably matches the 72 as reported by Costa. This is because the repository contains architecture-duplicated
code.


**ethereum/go-ethereum**: 

difference is +49.

Project description: Official Go implementation of the Ethereum protocol

The diff between v1.9.5 (20.09.2019) and 389da6aa48 (28.05.2020) shows 1 addition and 6858 deletions, therefore I
should see a difference of -6857. This is clearly completely different from the actual evidence.
In version v1.9.5, there isn't a `go.mod` file yet, instead there is a vendor directory inside the directory.

I cannot run my tool at v1.9.5 because of the missing `go.mod` file.

Using grep `rg '^[^/]*unsafe\.' -g '!vendor' -g '*.go' | wc -l`, I get 42 usages which exactly matches the number
reported by Costa. The difference to my value six months later then is due to vendoring differences.
